package tests

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"
)

func namedWorkflow(t *testing.T, name string) workflow {
	t.Helper()
	raw, err := os.ReadFile("../.github/workflows/" + name)
	if err != nil {
		t.Fatal(err)
	}
	var w workflow
	if err := yaml.Unmarshal(raw, &w); err != nil {
		t.Fatal(err)
	}
	return w
}

func policyScript(t *testing.T, id string) string {
	t.Helper()
	w := namedWorkflow(t, "reusable-version.yml")
	for _, job := range w.Jobs {
		for _, s := range job.Steps {
			if s.ID == id {
				if s.With["script"] == "" {
					t.Fatalf("%s has no executable policy", id)
				}
				return s.With["script"]
			}
		}
	}
	t.Fatalf("missing policy step %s", id)
	return ""
}

type policyResult struct {
	Outputs  map[string]string `json:"outputs"`
	Dispatch []map[string]any  `json:"dispatch"`
	Error    string            `json:"error"`
	Notices  []string          `json:"notices"`
}

func executePolicy(t *testing.T, id string, changes map[string]any) policyResult {
	t.Helper()
	fixture := map[string]any{
		"script": policyScript(t, id),
		"env": map[string]string{
			"DRY_RUN": "false", "CHECKED_RUN_ID": "42",
			"CHECKS_WORKFLOW":  ".github/workflows/ci.yml",
			"VERSION_WORKFLOW": "version.yml", "PR_NUMBER": "7", "LABEL_EVENT_ID": "12",
			"REQUIRED_JOBS": "verify / verify,checker,checker-macos",
		},
		"context": map[string]any{
			"eventName": "workflow_dispatch", "sha": strings.Repeat("a", 40),
			"ref":  "refs/heads/feat/example",
			"repo": map[string]string{"owner": "example", "repo": "app"},
			"payload": map[string]any{
				"repository": map[string]any{"default_branch": "main"},
			},
		},
		"run": map[string]any{
			"id": 42, "status": "completed", "conclusion": "success", "event": "pull_request",
			"path": ".github/workflows/ci.yml", "head_sha": strings.Repeat("a", 40),
			"head_branch":     "feat/example",
			"head_repository": map[string]string{"full_name": "example/app"},
		},
		"pr": map[string]any{
			"number": 7, "state": "open", "labels": []map[string]string{{"name": "release:preview"}},
			"head": map[string]any{"sha": strings.Repeat("a", 40), "ref": "feat/example",
				"repo": map[string]string{"full_name": "example/app"}},
		},
		"events": []map[string]any{{
			"id": 12, "event": "labeled", "label": map[string]string{"name": "release:preview"},
			"actor": map[string]string{"login": "writer"},
		}},
		"permission": "write",
		"jobs": []map[string]string{
			{"name": "verify / verify", "conclusion": "success"},
			{"name": "checker", "conclusion": "success"},
			{"name": "checker-macos", "conclusion": "success"},
		},
	}
	for key, value := range changes {
		fixture[key] = value
	}
	input, err := json.Marshal(fixture)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("node", "--input-type=module", "-e", `
import {readFileSync} from "node:fs";
const f = JSON.parse(readFileSync(0, "utf8"));
for (const key of ["run", "pr", "context", "env"]) {
  if (f[key + "Patch"]) Object.assign(f[key], f[key + "Patch"]);
}
if (f.request) {
  f.context.eventName = "workflow_run";
  f.context.payload.workflow_run = f.run;
}
const result = {outputs: {}, dispatch: [], notices: []};
const api = value => async () => {
  if (f.apiError) throw new Error("GitHub unavailable");
  return {data: value};
};
const github = {
  rest: {
    actions: {
      getWorkflowRun: api(f.run),
      listJobsForWorkflowRun: api({jobs: f.jobs}),
      createWorkflowDispatch: async request => { result.dispatch.push(request); },
    },
    pulls: {list: api(f.pr ? [f.pr] : []), get: api(f.pr)},
    issues: {listEvents: api(f.events)},
    repos: {getCollaboratorPermissionLevel: api({permission: f.permission}),
      getBranch: api({commit: {sha: f.branchSha || f.context.sha}})},
  },
  paginate: async (method, parameters) => {
    const data = (await method(parameters)).data;
    return data.jobs || data;
  },
};
const core = {
  setOutput: (key, value) => { result.outputs[key] = String(value); },
  notice: message => result.notices.push(message),
};
try {
  const AsyncFunction = Object.getPrototypeOf(async function(){}).constructor;
  await new AsyncFunction("github", "context", "core", "process", f.script)(
    github, f.context, core, {env: f.env});
} catch (error) {
  result.error = error.message;
}
console.log(JSON.stringify(result));
`)
	cmd.Stdin = bytes.NewReader(input)
	raw, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("executing actual workflow script: %v\n%s", err, raw)
	}
	var result policyResult
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("invalid policy result: %v\n%s", err, raw)
	}
	return result
}

func TestPreviewEligibility(t *testing.T) {
	cases := []struct {
		name    string
		changes map[string]any
		publish bool
		failed  bool
	}{
		{"writer enabled label", nil, true, false},
		{"maintainer enabled label", map[string]any{"permission": "maintain"}, true, false},
		{"admin enabled label", map[string]any{"permission": "admin"}, true, false},
		{"triager cannot enable publication", map[string]any{"permission": "triage"}, false, true},
		{"reader cannot enable publication", map[string]any{"permission": "read"}, false, true},
		{"missing PR cannot publish", map[string]any{"pr": nil}, false, false},
		{"dry run never authorizes publication", map[string]any{"env": map[string]string{"DRY_RUN": "true"}}, false, false},
		{"API failure is explicit", map[string]any{"apiError": true}, false, true},
		{"failed CI cannot authorize", map[string]any{"runPatch": map[string]any{"conclusion": "failure"}}, false, true},
		{"unfinished CI cannot authorize", map[string]any{"runPatch": map[string]any{"status": "in_progress"}}, false, true},
		{"different source cannot authorize", map[string]any{"runPatch": map[string]any{"head_sha": strings.Repeat("b", 40)}}, false, true},
		{"different workflow cannot authorize", map[string]any{"runPatch": map[string]any{"path": ".github/workflows/other.yml"}}, false, true},
		{"different branch cannot authorize", map[string]any{"runPatch": map[string]any{"head_branch": "feat/other"}}, false, true},
		{"feature push is not PR permission", map[string]any{"runPatch": map[string]any{"event": "push"}}, false, true},
		{"fork CI cannot authorize", map[string]any{"runPatch": map[string]any{"head_repository": map[string]string{"full_name": "fork/app"}}}, false, true},
		{"unlabeled PR stays read-only", map[string]any{"prPatch": map[string]any{"labels": []any{}}}, false, false},
		{"stale head stays read-only", map[string]any{"prPatch": map[string]any{"head": map[string]any{"sha": strings.Repeat("b", 40), "repo": map[string]string{"full_name": "example/app"}}}}, false, false},
		{"missing label history is explicit", map[string]any{"events": []any{}}, false, true},
		{"unlabel event cannot enable", map[string]any{"events": []map[string]any{{"id": 13, "event": "unlabeled", "label": map[string]string{"name": "release:preview"}}}}, false, true},
		{"missing CI identity is explicit", map[string]any{"envPatch": map[string]string{"CHECKED_RUN_ID": ""}}, false, true},
		{"skipped validation is not a pass", map[string]any{"jobs": []map[string]string{{"name": "verify / verify", "conclusion": "skipped"}}}, false, true},
		{"missing platform checks block publication", map[string]any{"jobs": []map[string]string{{"name": "verify / verify", "conclusion": "success"}}}, false, true},
		{"configured stable release branch is not a preview", map[string]any{
			"contextPatch": map[string]any{"ref": "refs/heads/release/1.2.3"},
			"runPatch":     map[string]any{"head_branch": "release/1.2.3"},
		}, false, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := executePolicy(t, "eligibility", tc.changes)
			if (got.Outputs["publish"] == "true") != tc.publish || (got.Error != "") != tc.failed {
				t.Fatalf("result = %+v", got)
			}
		})
	}
}

func TestSuccessfulChecksRequestBranchDispatch(t *testing.T) {
	for _, tc := range []struct {
		name    string
		changes map[string]any
		want    bool
	}{
		{"labeled PR update", nil, true},
		{"unlabeled PR", map[string]any{"prPatch": map[string]any{"labels": []any{}}}, false},
		{"feature push without PR event", map[string]any{"runPatch": map[string]any{"event": "push"}}, false},
		{"failed checks", map[string]any{"runPatch": map[string]any{"conclusion": "failure"}}, false},
		{"all jobs skipped", map[string]any{"jobs": []any{}}, false},
		{"wrong workflow", map[string]any{"runPatch": map[string]any{"path": ".github/workflows/other.yml"}}, false},
		{"fork", map[string]any{"runPatch": map[string]any{"head_repository": map[string]string{"full_name": "fork/app"}}}, false},
		{"PR closed", map[string]any{"pr": nil}, false},
		{"main push", map[string]any{"runPatch": map[string]any{"head_branch": "main", "event": "push"}}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			changes := map[string]any{"request": true}
			for key, value := range tc.changes {
				changes[key] = value
			}
			got := executePolicy(t, "request", changes)
			if got.Error != "" || (len(got.Dispatch) == 1) != tc.want {
				t.Fatalf("result = %+v", got)
			}
			if tc.want {
				branch := "feat/example"
				if tc.name == "main push" {
					branch = "main"
				}
				if got.Dispatch[0]["ref"] != branch {
					t.Fatalf("versioning was not dispatched on its source branch: %+v", got.Dispatch)
				}
			}
		})
	}
}

func TestEligibilityIsRechecked(t *testing.T) {
	for _, tc := range []struct {
		name    string
		changes map[string]any
		failed  bool
	}{
		{"same label remains enabled", nil, false},
		{"label removed", map[string]any{"prPatch": map[string]any{"labels": []any{}}}, true},
		{"PR closed", map[string]any{"prPatch": map[string]any{"state": "closed"}}, true},
		{"new source arrived", map[string]any{"contextPatch": map[string]any{"sha": strings.Repeat("b", 40)}}, true},
		{"label removed and added again", map[string]any{"envPatch": map[string]string{"LABEL_EVENT_ID": "11"}}, true},
		{"main unchanged", map[string]any{"envPatch": map[string]string{"PR_NUMBER": ""}}, false},
		{"main advanced", map[string]any{"envPatch": map[string]string{"PR_NUMBER": ""}, "branchSha": strings.Repeat("b", 40)}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := executePolicy(t, "recheck", tc.changes)
			if (got.Error != "") != tc.failed || (got.Outputs["publish"] == "true") == tc.failed {
				t.Fatalf("result = %+v", got)
			}
		})
	}
}

func TestVersioningRemainsBranchBased(t *testing.T) {
	w := namedWorkflow(t, "reusable-version.yml")
	if !strings.Contains(w.Jobs["request"].If, "!inputs.dry-run") {
		t.Error("read-only mode must also prohibit publication-request dispatches")
	}
	version := w.Jobs["version"]
	checkout := version.Steps[role(t, version.Steps, "actions/checkout")]
	if checkout.With["ref"] != "${{ github.ref }}" {
		t.Error("GitVersion must receive the selected branch, not a synthetic PR ref")
	}
	tagger := version.Steps[role(t, version.Steps, "michielvha/gitversion-tag-action")]
	if tagger.With["prefix"] != "v" || !strings.Contains(tagger.With["dry-run"], "publish") {
		t.Error("only the public GitVersion action may calculate and push the authorized exact tag")
	}
	for _, job := range w.Jobs {
		for _, s := range job.Steps {
			if strings.Contains(s.Run, "git tag") || strings.Contains(s.Run, "git push") ||
				strings.Contains(s.Run, "scripts/tag.sh") {
				t.Error("shared orchestration must not implement tagging")
			}
		}
	}
}

func TestSnapshotChecksDoNotPublish(t *testing.T) {
	w := namedWorkflow(t, "reusable-go-ci.yml")
	job := w.Jobs["verify"]
	snapshot := job.Steps[role(t, job.Steps, "goreleaser/goreleaser-action")]
	if snapshot.If != "${{ inputs.test-release-packaging }}" {
		t.Error("the caller must be able to select Go CI without snapshot packaging")
	}
	for _, flag := range []string{"--snapshot", "--skip=publish,sign", "--clean"} {
		if !strings.Contains(snapshot.With["args"], flag) {
			t.Errorf("snapshot lacks %s", flag)
		}
	}
	raw, err := os.ReadFile("../.github/workflows/reusable-go-ci.yml")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "pgp-private-key") || strings.Contains(string(raw), "release-app-private-key") {
		t.Error("validation and packaging must not receive production credentials")
	}
}
