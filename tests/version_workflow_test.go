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
	if steps[app].If != "${{ !inputs.dry-run && inputs.release-app-client-id != '' }}" ||
		steps[app].With["owner"] != "${{ github.repository_owner }}" ||
		steps[app].With["repositories"] != "${{ github.event.repository.name }}" ||
		steps[app].With["permission-contents"] != "write" {
		t.Error("an App token is optional, single-repository Contents: write, and only for tagging")
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
	// Inheriting the caller's token lets it choose: read with an App token
	// (separate workflows) or write with GITHUB_TOKEN (same workflow).
	jobPermissions := w.Jobs["version"].Permissions
	if !w.Permissions.IsZero() || !jobPermissions.IsZero() {
		t.Error("the wrapper must inherit the caller job's token permissions")
	}
	if steps[version].ID == "" {
		t.Fatal("the GitVersion step needs an id to expose its outputs")
	}
	for output, source := range map[string]string{
		"tagged":  "tagged",
		"tag":     "pushed-tag",
		"version": "semVer",
	} {
		want := "${{ steps." + steps[version].ID + ".outputs." + source + " }}"
		if w.Jobs["version"].Outputs[output] != want ||
			w.On.WorkflowCall.Outputs[output].Value != "${{ jobs.version.outputs."+output+" }}" {
			t.Errorf("same-workflow callers need the %s output from GitVersion's %s", output, source)
		}
	}
}
