package tests

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"
)

type callable struct {
	On struct {
		WorkflowCall struct {
			Inputs  map[string]yaml.Node `yaml:"inputs"`
			Secrets map[string]yaml.Node `yaml:"secrets"`
			Outputs map[string]yaml.Node `yaml:"outputs"`
		} `yaml:"workflow_call"`
	} `yaml:"on"`
}

type caller struct {
	Jobs map[string]struct {
		Uses    string            `yaml:"uses"`
		With    map[string]string `yaml:"with"`
		Secrets map[string]string `yaml:"secrets"`
	} `yaml:"jobs"`
}

var (
	sharedCall = regexp.MustCompile(`^vhco-pro/\.github/\.github/workflows/([\w.-]+\.yml)@`)
	outputRef  = regexp.MustCompile(`needs\.(\w+)\.outputs\.([\w-]+)`)
)

// The examples are what callers copy, so they must match the shared
// workflows' current inputs, secrets and outputs.
func TestExamplesMatchSharedWorkflows(t *testing.T) {
	files, err := filepath.Glob("../examples/*/*/*.yml")
	if err != nil || len(files) == 0 {
		t.Fatalf("no examples found: %v", err)
	}
	load := func(path string, into any) {
		t.Helper()
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := yaml.Unmarshal(raw, into); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
	}
	for _, file := range files {
		raw, _ := os.ReadFile(file)
		var c caller
		load(file, &c)
		called := map[string]callable{}
		for name, job := range c.Jobs {
			m := sharedCall.FindStringSubmatch(job.Uses)
			if m == nil {
				continue
			}
			var w callable
			load(filepath.Join("..", ".github", "workflows", m[1]), &w)
			called[name] = w
			for input := range job.With {
				if _, ok := w.On.WorkflowCall.Inputs[input]; !ok {
					t.Errorf("%s job %s passes unknown input %q to %s", file, name, input, m[1])
				}
			}
			for secret := range job.Secrets {
				if _, ok := w.On.WorkflowCall.Secrets[secret]; !ok {
					t.Errorf("%s job %s passes unknown secret %q to %s", file, name, secret, m[1])
				}
			}
		}
		for _, ref := range outputRef.FindAllStringSubmatch(string(raw), -1) {
			w, ok := called[ref[1]]
			if !ok {
				continue
			}
			if _, ok := w.On.WorkflowCall.Outputs[ref[2]]; !ok {
				t.Errorf("%s reads unknown output %s.%s", file, ref[1], ref[2])
			}
		}
		if !strings.Contains(string(raw), "@FULL_COMMIT_SHA") && len(called) > 0 {
			t.Errorf("%s must use the FULL_COMMIT_SHA placeholder for this repository's workflows", file)
		}
	}
}
