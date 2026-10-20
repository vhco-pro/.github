package tests

import (
	"os"
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestVersionWorkflowPlumbing(t *testing.T) {
	raw, err := os.ReadFile("../.github/workflows/reusable-version.yml")
	if err != nil {
		t.Fatal(err)
	}
	var w workflow
	if err := yaml.Unmarshal(raw, &w); err != nil {
		t.Fatal(err)
	}
	steps := w.Jobs["version"].Steps
	app := role(t, steps, "actions/create-github-app-token")
	checkout := role(t, steps, "actions/checkout")
	version := role(t, steps, "michielvha/gitversion-tag-action")
	if app >= checkout || checkout >= version {
		t.Fatal("GitVersion must follow authentication and full-history branch checkout")
	}
	if steps[app].If != "${{ !inputs.dry-run }}" ||
		steps[app].With["owner"] != "${{ github.repository_owner }}" ||
		steps[app].With["repositories"] != "${{ github.event.repository.name }}" ||
		steps[app].With["permission-contents"] != "write" {
		t.Error("only publishing may mint a single-repository Contents: write token")
	}
	if steps[checkout].With["ref"] != "${{ github.ref }}" ||
		steps[checkout].With["fetch-depth"] != "0" ||
		steps[checkout].With["persist-credentials"] != "${{ !inputs.dry-run }}" {
		t.Error("GitVersion needs branch context and full history, without write credentials in dry-run")
	}
	if steps[version].With["dry-run"] != "${{ inputs.dry-run }}" ||
		steps[version].With["configFilePath"] != "${{ inputs.config-file-path }}" {
		t.Error("the public GitVersion action must own calculation and optional tagging")
	}
	for _, forbidden := range []string{"github-script@", "workflow_run", "checked-run-id", "required-jobs", "git tag ", "git push "} {
		if strings.Contains(string(raw), forbidden) {
			t.Errorf("thin versioning must not contain %q", forbidden)
		}
	}
}
