// +skip_license_check
/*
Copyright 2022 The cert-manager Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

// Note for developers:
// If you want to edit how tests are generated, change: ./pkg/
// If you want to edit which tests are generated on each branch / k8s version, change: ./prowspecs/

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"

	"prowgen/prowspecs"
)

const (
	generateProwCommand         = "prowgen"
	generateProwDescription     = "Generate YAML specifying prow tests for cert-manager"
	generateProwLongDescription = `prowgen creates prow test specifications for a given cert-manager release branch. These specifications
define the Prow tests available to be run against a given branch.

Generated tests include both presubmit tests (tests that can be run against PRs) and periodic
tests (tests which are run on a schedule, independently of PRs).

By generating this config we avoid the need for humans to edit YAML manually
which is error-prone.

If --output-dir is set, the generated YAML will be written to the specified
directory with the correct directory format which prow expects. Otherwise,
generated output will be written to stdout.
`
)

var (
	generateProwExample = fmt.Sprintf(`
To generate tests for the a branch called foo:

	%s --branch=foo
`, generateProwCommand)
)

type generateProwOptions struct {
	// Branch specifies the name of the branch whose tests should be generated
	Branch string

	// OutputDir specifies the dir to output the yaml files to. If empty, output
	// will be written to stdout.
	OutputDir string
}

func (o *generateProwOptions) AddFlags(fs *flag.FlagSet, markRequired func(string)) {
	fs.StringVar(&o.Branch, "branch", "", fmt.Sprintf("Type of tests to generate; one of ('*' generates all branches) %v", append(prowspecs.KnownBranches(), "*")))
	fs.StringVarP(&o.OutputDir, "output-dir", "o", "", "OutputDir specifies the dir to output the yaml files to. If empty, output will be written to stdout.")

	markRequired("branch")
}

func generateProwCmd() *cobra.Command {
	o := &generateProwOptions{}

	cmd := &cobra.Command{
		Use:          generateProwCommand,
		Short:        generateProwDescription,
		Long:         generateProwLongDescription,
		Example:      generateProwExample,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.Branch == "*" {
				for _, branch := range prowspecs.KnownBranches() {
					if err := o.runGenerateProw(branch); err != nil {
						return err
					}
				}
				return nil
			}
			return o.runGenerateProw(o.Branch)
		},
	}

	o.AddFlags(cmd.Flags(), func(s string) {
		if err := cmd.MarkFlagRequired(s); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	})

	return cmd
}

// sanitizedArgs strips the path from the command which was used to invoke the script,
// so we don't include things like "/home/workspace/release/bin/prowgen"
func sanitizedArgs() []string {
	args := os.Args[:]

	for i := range args {
		if !strings.Contains(args[i], "/") {
			continue
		}

		args[i] = filepath.Base(args[i])
	}

	return args
}

func (o *generateProwOptions) runGenerateProw(branch string) error {
	spec, err := prowspecs.SpecForBranch(branch)
	if err != nil {
		return err
	}

	jobFile := spec.GenerateJobFile()

	out, err := yaml.Marshal(jobFile)
	if err != nil {
		return err
	}

	prelude := fmt.Sprintf(
		`# THIS FILE HAS BEEN AUTOMATICALLY GENERATED
# Don't manually edit it; instead edit the "prowgen" tool which generated it
# Generated with: %s

`,
		strings.Join(sanitizedArgs(), " "),
	)

	data := prelude + string(out)

	if o.OutputDir == "" {
		fmt.Println(data)
		return nil
	}

	branchPath := path.Join(o.OutputDir, branch)

	if err := os.MkdirAll(branchPath, 0755); err != nil {
		return err
	}

	path := filepath.Join(branchPath, fmt.Sprintf("cert-manager-%s.yaml", branch))
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, strings.NewReader(data)); err != nil {
		return err
	}

	return nil
}

func main() {
	cmd := generateProwCmd()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
