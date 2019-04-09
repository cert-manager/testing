/*
Copyright 2019 The Jetstack contributors.

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

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// builder builds docker images given a build.yaml file provided in the
// images build directory.
// It makes it easy to build a matrix of images, setting different build args
// for each one.
// It also handles templating image names and pushing docker images.
// It is used by the testing repository to push testing docker images used in
// ProwJobs.

var (
	confirm  bool
	registry string
	buildDir string
)

func init() {
	flag.BoolVar(&confirm, "confirm", false, "set to true to confirm pushing images")
	flag.StringVar(&registry, "registry", "eu.gcr.io/jetstack-build-infra-images", "docker image registry to push images to")
	flag.StringVar(&buildDir, "build-dir", "", "path to a directory containing a build.yaml file")
}

func main() {
	flag.Parse()

	// validate flags
	if errs := validateFlags(); len(errs) > 0 {
		for _, err := range errs {
			log.Println(err.Error())
		}
		os.Exit(1)
	}

	if !confirm {
		log.Printf("--confirm is set to false, not pushing images")
	}

	cfg, err := parseConfig(buildDir + "/build.yaml")
	if err != nil {
		log.Printf("error reading build.yaml: %v", err)
		os.Exit(1)
	}

	ctxs, err := buildContexts(*cfg)
	if err != nil {
		log.Printf("error constructing build contexts: %v", err)
		os.Exit(1)
	}

	for name, ctx := range ctxs {
		log.Printf("building variant %q", name)
		if err := ctx.Build(); err != nil {
			log.Printf("error building variant %q: %v", name, err)
			os.Exit(1)
		}
		log.Printf("built variant %q", name)
	}

	log.Printf("build all variants")
	if !confirm {
		log.Printf("skipping pushing images")
		os.Exit(0)
	}

	for name, ctx := range ctxs {
		imageNames, err := allImageNames(cfg, ctx, name, cfg.Images...)
		if err != nil {
			log.Printf("error determining image names: %v", err)
			os.Exit(1)
		}

		for _, img := range imageNames {
			log.Printf("pushing image %q", img)
			if err := ctx.Push(img); err != nil {
				log.Printf("error pushing image %q: %v", img, err)
				os.Exit(1)
			}
			log.Printf("pushed image %q", img)
		}
	}

	log.Printf("SUCCESS")
	os.Stdout.Write([]byte(path.Join(registry, cfg.Name)))
}

func allImageNames(cfg *buildConfig, ctx *buildContext, variant string, templates ...string) ([]string, error) {
	switch variant {
	case "":
		templates = append(templates,
			"${_REGISTRY}/${_NAME}:${_DATE_STAMP}-${_GIT_REF}",
			"${_REGISTRY}/${_NAME}:latest",
		)
	default:
		templates = append(templates,
			"${_REGISTRY}/${_NAME}:${_DATE_STAMP}-${_GIT_REF}-${_VARIANT}",
			"${_REGISTRY}/${_NAME}:latest-${_VARIANT}",
		)
	}

	imageNames := make(strSet)
	for _, t := range templates {
		img, err := formatImageName(cfg, ctx, variant, t)
		if err != nil {
			log.Printf("error generating image name: %v", err)
			return nil, err
		}

		imageNames.Add(img)
	}

	return imageNames.Slice(), nil
}

type strSet map[string]struct{}

func (s strSet) Slice() []string {
	out := make([]string, len(s))
	i := 0
	for k := range s {
		out[i] = k
		i++
	}
	return out
}

func (s strSet) Add(strs ...string) {
	for _, str := range strs {
		s[str] = struct{}{}
	}
}

func formatImageName(cfg *buildConfig, ctx *buildContext, variant string, tmpl string) (string, error) {
	tmplMap := make(map[string]string)
	for k, v := range ctx.BuildArgs {
		tmplMap[k] = v
	}
	gitRef, err := getGitRef()
	if err != nil {
		return "", err
	}
	tmplMap["_NAME"] = cfg.Name
	tmplMap["_REGISTRY"] = registry
	tmplMap["_DATE_STAMP"] = time.Now().Format("20060102")
	tmplMap["_GIT_REF"] = gitRef
	tmplMap["_VARIANT"] = variant

	img := tmpl
	for k, v := range tmplMap {
		img = strings.ReplaceAll(img, fmt.Sprintf("${%s}", k), v)
	}

	return img, nil
}

func getGitRef() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--always", "--dirty")
	cmd.Dir = buildDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func validateFlags() []error {
	var errs []error
	if buildDir == "" {
		errs = append(errs, fmt.Errorf("build-dir must be specified"))
	}
	return errs
}

type buildConfig struct {
	Name       string             `json:"name"`
	Dockerfile string             `json:"dockerfile"`
	Arguments  map[string]string  `json:"arguments"`
	Variants   map[string]variant `json:"variants"`
	Images     []string           `json:"images"`
}

type variant struct {
	Arguments map[string]string `json:"arguments"`
}

func parseConfig(path string) (*buildConfig, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg buildConfig
	if err := yaml.Unmarshal(d, &cfg); err != nil {
		return nil, err
	}

	applyDefaults(&cfg)

	if errs := validateConfig(cfg); len(errs) > 0 {
		return nil, fmt.Errorf("config file error: %v", errs)
	}

	return &cfg, nil
}

func applyDefaults(cfg *buildConfig) {
	if cfg.Dockerfile == "" {
		cfg.Dockerfile = "Dockerfile"
	}
}

func validateConfig(cfg buildConfig) []error {
	var errs []error
	if cfg.Name == "" {
		errs = append(errs, fmt.Errorf("image field must be set"))
	}
	if cfg.Dockerfile == "" {
		errs = append(errs, fmt.Errorf("dockerfile field must be set"))
	}
	return errs
}

// buildContexts constructs a slice of buildContexts for the given config
// variations will be expanded in this function.
func buildContexts(cfg buildConfig) (map[string]*buildContext, error) {
	if len(cfg.Variants) == 0 {
		ctx := constructContext(cfg, nil)
		return map[string]*buildContext{"": ctx}, nil
	}

	ctxs := make(map[string]*buildContext)
	for name, v := range cfg.Variants {
		ctx := constructContext(cfg, v.Arguments)
		ctxs[name] = ctx
	}

	return ctxs, nil
}

func constructContext(cfg buildConfig, extraArgs map[string]string) *buildContext {
	ctx := buildContext{
		Dockerfile: cfg.Dockerfile,
		Directory:  buildDir,
	}
	buildArgs := make(map[string]string)
	for k, v := range cfg.Arguments {
		buildArgs[k] = v
	}
	for k, v := range extraArgs {
		buildArgs[k] = v
	}
	ctx.BuildArgs = buildArgs
	return &ctx
}

// buildContext provides an abstraction to build docker images using different
// docker build systems.
// Initially only docker is supported.
type buildContext struct {
	Dockerfile string
	Directory  string
	BuildArgs  map[string]string

	name  string
	built bool

	nameLock  sync.Mutex
	buildLock sync.Mutex
}

// Build will build the docker image given the context config
func (b *buildContext) Build() error {
	b.buildLock.Lock()
	defer b.buildLock.Unlock()
	if b.built {
		return nil
	}

	log.Printf("building docker image dockerfile=%s, directory=%s, buildArgs=%v", b.Dockerfile, b.Directory, b.BuildArgs)
	args := b.buildCmd()
	if err := b.runDocker(args...); err != nil {
		return err
	}
	log.Printf("built docker image")
	b.built = true
	return nil
}

func (b *buildContext) buildCmd() []string {
	args := []string{"build", "-t", b.temporaryImageName(), "-f", path.Join(b.Directory, b.Dockerfile)}
	for k, v := range b.BuildArgs {
		args = append(args, "--build-arg", k+"="+v)
	}
	args = append(args, b.Directory)
	return args
}

func (b *buildContext) temporaryImageName() string {
	b.nameLock.Lock()
	defer b.nameLock.Unlock()

	if b.name == "" {
		b.name = randString(16)
	}

	return "builder:" + b.name
}

// Push will push the docker image that has been built with the image name
// provided.
// If Build has not been called, the image will be built.
// It is safe to call this function multiple times in parallel.
func (b *buildContext) Push(name string) error {
	if err := b.Build(); err != nil {
		return err
	}

	if err := b.runDocker("tag", b.temporaryImageName(), name); err != nil {
		return err
	}

	if err := b.runDocker("push", name); err != nil {
		return err
	}

	return nil
}

func (b *buildContext) runDocker(args ...string) error {
	log.Printf("running with args %v", args)
	cmd := exec.Command("docker", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
