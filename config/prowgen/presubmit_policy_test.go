package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"
)

// policyConfig models just enough of the Prow config and job-config schema to
// enforce the presubmit credential policy below.
type policyConfig struct {
	Presets []struct {
		Labels map[string]string `yaml:"labels"`
		Env    []struct {
			ValueFrom struct {
				SecretKeyRef map[string]string `yaml:"secretKeyRef"`
			} `yaml:"valueFrom"`
		} `yaml:"env"`
		Volumes []struct {
			Secret map[string]interface{} `yaml:"secret"`
		} `yaml:"volumes"`
	} `yaml:"presets"`
	Presubmits map[string][]struct {
		Name   string            `yaml:"name"`
		Labels map[string]string `yaml:"labels"`
	} `yaml:"presubmits"`
}

// TestNoPresubmitUsesSecretBearingPreset asserts that no presubmit job in the
// whole repository carries a preset label which injects a Secret, whether via
// env secretKeyRef or a Secret volume. Presubmits run unreviewed code from
// pull requests, so a presubmit with access to live credentials lets anyone
// who can open a PR exfiltrate them.
//
// The forbidden label set is derived from the preset definitions themselves,
// so a newly added credential preset is forbidden automatically, and the
// check covers hand-written job files as well as prowgen-generated ones
// (which are committed and kept in sync by verify-prowgen). This is the
// repo-wide enforcement backing the generation-time strip of
// presubmitForbiddenLabels in pkg/context.go.
func TestNoPresubmitUsesSecretBearingPreset(t *testing.T) {
	files := []string{"../config.yaml"}
	err := filepath.WalkDir("../jobs", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".yaml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	configs := make(map[string]policyConfig, len(files))
	for _, file := range files {
		raw, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		var config policyConfig
		if err := yaml.Unmarshal(raw, &config); err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		configs[file] = config
	}

	// Pass 1: labels of presets which inject a Secret, from any config file.
	secretPresetLabels := map[string]string{}
	for file, config := range configs {
		for _, preset := range config.Presets {
			secretBearing := false
			for _, env := range preset.Env {
				if len(env.ValueFrom.SecretKeyRef) > 0 {
					secretBearing = true
				}
			}
			for _, volume := range preset.Volumes {
				if len(volume.Secret) > 0 {
					secretBearing = true
				}
			}
			if secretBearing {
				for label := range preset.Labels {
					secretPresetLabels[label] = file
				}
			}
		}
	}
	if len(secretPresetLabels) == 0 {
		t.Fatal("found no secret-bearing presets at all; the preset parsing in this test is probably broken")
	}

	// Pass 2: no presubmit anywhere may carry one of those labels.
	for file, config := range configs {
		for repo, presubmits := range config.Presubmits {
			for _, presubmit := range presubmits {
				for label := range presubmit.Labels {
					if definedIn, forbidden := secretPresetLabels[label]; forbidden {
						t.Errorf("%s: presubmit %s (%s) uses preset %s (defined in %s), which injects a Secret; presubmits run unreviewed PR code and must not have access to live credentials",
							file, presubmit.Name, repo, label, definedIn)
					}
				}
			}
		}
	}
}
