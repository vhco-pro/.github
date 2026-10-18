import { expect, test } from "bun:test";

const workflow = Bun.YAML.parse(await Bun.file(
  new URL("../.github/workflows/reusable-go-release.yml", import.meta.url),
).text());
const steps = workflow.jobs.release.steps;

test("the shared Go workflow only composes checkout, CI, and release actions", () => {
  expect(steps.map((step) => step.uses.split("@")[0])).toEqual([
    "actions/checkout", "michielvha/go-ci-action", "michielvha/goreleaser-action",
  ]);
  expect(steps.every((step) => !step.run && !step["continue-on-error"])).toBe(true);
  expect(steps[2].with["go-version"]).toBe("${{ steps.ci.outputs.go-version }}");
});

test("publication uses the event commit and excludes branch or deleted-tag events", () => {
  expect(steps[0].with.ref).toBe("${{ github.sha }}");
  expect(steps[0].with["persist-credentials"]).toBe(false);
  expect(workflow.jobs.release.if).toContain("github.event_name == 'push'");
  expect(workflow.jobs.release.if).toContain("refs/tags/v");
  expect(workflow.jobs.release.if).toContain("!github.event.deleted");
});

test("the workflow does not provision taggers or release helper scripts", () => {
  expect(Object.keys(workflow.on.workflow_call.inputs)).toEqual([
    "go-version-file", "golangci-lint-config",
  ]);
  expect(steps.slice(0, 2).every((step) => !step.with["pgp-private-key"])).toBe(true);
  expect(steps[2].with["pgp-private-key"]).toBe("${{ secrets.pgp-private-key }}");
});

test("all action implementations are pinned and publication is serialized", () => {
  expect(workflow.concurrency["cancel-in-progress"]).toBe(false);
  for (const step of steps) {
    expect(step.uses).toMatch(/@[0-9a-f]{40}$/);
  }
});
