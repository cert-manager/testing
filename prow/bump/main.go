package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// bump is a tool for automatically bumping the Prow version needed in the
// various places required in order to roll out a new version.
// It will create a pull request against the jetstack/testing repository with
// the updated changes.

var (
	versionFile string
	repoRoot    string
	// a file containing a line `image: foo:tag` that is used to detect the
	// 'source' image tag that is being bumped *from*.
	// This is used to string replace files contained in the directoryList
	// directories.
	existingImageFile = "prow/cluster/tide_deployment.yaml"
	directoryList     = []string{
		"prow/cluster",
		"config/jobs/testing",
	}
)

func init() {
	flag.StringVar(&versionFile, "version-file", "prow/version", "path to a file containing the image tag that should be set")
	flag.StringVar(&repoRoot, "repo-root", "", "base path used as a prefix for all other file paths")
}

func main() {
	flag.Parse()

	existingVersion, err := detectExistingVersion()
	if err != nil {
		log.Printf("error detecting existing version: %v", err)
		os.Exit(1)
	}

	newVersion, err := getNewVersion()
	if err != nil {
		log.Printf("error detecting new version: %v", err)
		os.Exit(1)
	}

	files, err := findFiles(directoryList...)
	if err != nil {
		log.Printf("error enumerating files to patch: %v", err)
		os.Exit(1)
	}

	log.Printf("detected files to patch: %v", files)

	patchedFiles, err := patchFiles(existingVersion, newVersion, files...)
	if err != nil {
		log.Printf("error patching files: %v", err)
		os.Exit(1)
	}

	log.Printf("patched %d files", len(patchedFiles))
}

func patchFiles(old, new string, paths ...string) ([]string, error) {
	var updated []string
	for _, p := range paths {
		d, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, err
		}
		mode := os.FileMode(0644)
		fi, err := os.Stat(p)
		if err == nil {
			mode = fi.Mode()
		}
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}

		in := string(d)
		out := strings.ReplaceAll(in, old, new)

		if in == out {
			log.Printf("no change to file %q detected, skipping", p)
			continue
		}

		if err := ioutil.WriteFile(p, []byte(out), mode); err != nil {
			return nil, err
		}

		log.Printf("updated file %q", p)
		updated = append(updated, p)
	}
	return updated, nil
}

var existingVersionRE = regexp.MustCompile(`image: gcr\.io/k8s-prow/tide:(.+)`)

func detectExistingVersion() (string, error) {
	d, err := ioutil.ReadFile(path.Join(repoRoot, existingImageFile))
	if err != nil {
		return "", err
	}

	matches := existingVersionRE.FindStringSubmatch(string(d))
	if len(matches) != 2 {
		return "", fmt.Errorf("error extracting image tag from file %q (matches: %v)", existingImageFile, matches)
	}

	tag := matches[1]
	log.Printf("detected old image tag %q", tag)

	return tag, nil
}

func getNewVersion() (string, error) {
	d, err := ioutil.ReadFile(path.Join(repoRoot, versionFile))
	if err != nil {
		return "", err
	}

	v := strings.TrimSpace(string(d))
	log.Printf("detected new image tag %q", v)
	return v, nil
}

func findFiles(paths ...string) ([]string, error) {
	var files []string
	for _, p := range paths {
		if err := filepath.Walk(path.Join(repoRoot, p), func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			files = append(files, path)
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return files, nil
}
