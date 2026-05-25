<!--
Sync Impact Report
Version change: N/A → 1.0.0
Modified principles: N/A
Added principles:
- I. Consumer Financial Clarity
- II. Deterministic Calculation Core
- III. Evidence-Backed Assumptions
- IV. CLI/TUI Contract Stability
- V. Test-First Financial Changes
Added sections:
- Domain Constraints
- Development Workflow & Quality Gates
Removed sections: None
Templates requiring updates:
- ⚠ pending: .specify/templates/plan-template.md
- ⚠ pending: .specify/templates/spec-template.md
- ⚠ pending: .specify/templates/tasks-template.md
Follow-up TODOs: None
-->

# BCPS Constitution

## Core Principles

### I. Consumer Financial Clarity

BCPS exists to help users understand the approximate repayment burden of an Income Sharing
Agreement. Every feature MUST improve transparency for the user, not merely add presentation
or implementation novelty.

The tool MUST make the financial consequence of an ISA understandable through concrete outputs:
annual salary assumptions, annual repayment amounts, CPI-adjusted threshold values, total repaid,
multiples of borrowed sum repaid, and equivalent education-loan interest rate where applicable.
Any feature that changes the repayment model MUST preserve or improve the user’s ability to
compare ISA repayment against a conventional loan-style mental model.

BCPS MUST clearly communicate that outputs are projections based on user-provided assumptions
and embedded model assumptions. It MUST NOT present results as legal, tax, investment, or
contractual advice.

Rationale: The project’s purpose is transparency for students evaluating ISA repayment exposure.
A technically correct calculation is insufficient if users cannot understand the consequence.

### II. Deterministic Calculation Core

Financial calculations MUST be deterministic, reproducible, and separated from terminal UI concerns.

Repayment, salary progression, CPI adjustment, threshold calculation, total repayment capping,
and equivalent loan-interest calculation MUST be expressible as pure domain logic with explicit
inputs and outputs. New calculation behavior MUST NOT depend on terminal state, interactive form
state, package-level mutable globals, wall-clock time, environment variables, or hidden network
calls.

The CLI/TUI MAY collect input interactively, but it MUST pass validated values into the calculation
core explicitly. Rendering and styling MUST remain outside the calculation core.

Rationale: Users need to trust the schedule. Deterministic calculations make the model reviewable,
testable, and reusable outside the interactive UI.

### III. Evidence-Backed Assumptions

All financial assumptions that are not directly entered by the user MUST be visible, justified, and
versioned.

Historical CPI values, forecast CPI values, default model assumptions, repayment caps, and any
contract-specific rules MUST cite their source in code comments or documentation. When assumptions
change, the change MUST include a rationale and tests showing the effect on representative schedules.

Forecast values MUST be clearly distinguished from historical values. Future projected values MUST
be configurable or explicitly documented when hardcoded.

Rationale: A payment schedule predictor can silently become misleading if embedded assumptions
age, drift, or are treated as facts after they were only forecasts.

### IV. CLI/TUI Contract Stability

The command-line entry point and interactive flow are product contracts.

The primary command MUST remain runnable as `bcps` or `go run main.go` unless a migration path
is documented. Prompt wording, input semantics, validation behavior, and output labels SHOULD be
stable across releases. Breaking changes to input interpretation or output meaning MUST be treated
as user-facing breaking changes, even if the Go API does not change.

The TUI MUST validate numeric inputs before calculation. Invalid input MUST fail early with a
human-readable error and MUST NOT produce a partial or silently coerced schedule. Quit behavior
for common terminal exits such as `esc`, `ctrl+c`, and `q` MUST remain predictable.

Rationale: BCPS is primarily an interactive CLI. Users experience the product through prompts,
validation, and rendered schedules, so those are part of the compatibility surface.

### V. Test-First Financial Changes

Changes to financial behavior are NON-NEGOTIABLE test-first changes.

Any change to repayment calculation, threshold calculation, CPI handling, salary progression,
equivalent interest-rate calculation, input validation, or schedule termination behavior MUST begin
with tests that describe the intended behavior. These tests MUST fail against the old behavior before
implementation begins, unless the change is a pure refactor with no behavior change.

At minimum, calculation changes MUST include:
- golden or table-driven tests for representative repayment schedules;
- edge-case tests for zero/low repayment, high salary growth, cap-triggered repayment termination,
  and repayment start year validation;
- regression tests for any previously discovered calculation bug.

Rationale: Financial projections are easy to break with small changes. Test-first development
protects the user-facing promise of clarity and correctness.

## Domain Constraints

BCPS is a financial projection tool for Income Sharing Agreement repayment schedules.

The following constraints are mandatory:

1. Projection, not advice: outputs MUST be framed as approximate projections based on supplied
   assumptions.
2. Explicit units: monetary values, percentages, years, CPI rates, and salary-growth rates MUST be
   labeled unambiguously in prompts, docs, and output.
3. No hidden data fetches: calculations MUST NOT depend on runtime network calls unless a feature
   specification explicitly introduces data refresh behavior, caching, provenance display, and offline
   behavior.
4. Source transparency: embedded financial data MUST include source and date context.
5. Reproducibility: a user SHOULD be able to reproduce a schedule from the displayed or exported
   inputs and documented assumptions.
6. Privacy by default: user-entered salary, repayment, and borrowing values MUST remain local unless
   a future feature explicitly introduces export or sharing, with clear user consent.
7. Go-first implementation: the project remains a Go CLI/TUI unless a specification explicitly justifies
   a broader architecture. New dependencies MUST be justified by user value, maintenance cost, and
   impact on terminal UX.

## Development Workflow & Quality Gates

All feature work MUST follow spec-driven development:

1. Start with a feature specification that names the user goal, input assumptions, expected output, and
   acceptance scenarios.
2. Define calculation examples before implementation for any feature that changes financial behavior.
3. Add tests before changing domain logic.
4. Keep domain logic independent from Bubble Tea, Huh, Lipgloss, Cobra, and other UI/command
   packages.
5. Run formatting, tests, and static checks before merge.
6. Update README or user-facing documentation when prompts, outputs, assumptions, or run commands
   change.
7. Include a compatibility note for any change that alters existing prompt semantics or output meaning.

Pull requests MUST identify whether they affect:
- calculation behavior;
- assumption data;
- CLI/TUI interaction;
- output format;
- dependency/runtime setup;
- documentation only.

A change affecting calculation behavior or assumption data MUST include before/after examples in
the PR description or feature spec.

## Governance

This constitution supersedes ad hoc development preferences for BCPS. When a specification,
implementation plan, or task list conflicts with this constitution, the constitution takes precedence.

Amendments require:
1. a written rationale for the change;
2. an explicit version bump;
3. updates to affected Spec Kit templates and project documentation;
4. migration guidance when principles or user-facing contracts change.

Versioning follows semantic versioning:
- MAJOR: removes or redefines a core principle, weakens user-protection guarantees, or changes
  governance compatibility.
- MINOR: adds a new principle, adds a required quality gate, or materially expands an existing rule.
- PATCH: clarifies wording, fixes ambiguity, or updates non-semantic guidance.

Every feature plan MUST include a Constitution Check that verifies:
- calculation transparency;
- deterministic calculation core boundaries;
- assumption provenance;
- CLI/TUI contract impact;
- test-first coverage for financial behavior.

Compliance is reviewed during planning and PR review. A PR that changes financial behavior without
tests or assumption provenance MUST NOT be merged.

**Version**: 1.0.0 | **Ratified**: 2026-05-25 | **Last Amended**: 2026-05-25
