# Roadmap

## Phase 0 - Public Boundary

- Create public-boundary docs. Done.
- Add synthetic fixture schema. Done.
- Add gitleaks gate. Done.
- Publish only after local CI, full-history gitleaks, and manual boundary
  review pass. Done.

## Phase 1 - Synthetic CLI

- Add local JSON fixtures. Done.
- Add read-only opportunity list and resume-tailoring demo commands. Done.
- Add tests that reject non-synthetic fixture markers. Done.

## Phase 2 - MCP-Style Surface

- Add a small manifest of MCP-style tools. Done.
- Add in-process demo calls for opportunity review, resume tailoring, and
  interview prep. Done.

## Phase 3 - Public Release Review

- Run full secret scan and history review. Done.
- Add CI and community files. Done.
- Add one-command local release gate. Done with `make ci`.
- Switch visibility only after the public boundary is reviewed. Done.

## Phase 4 - Next Proof Improvements

- Add a small README demo transcript generated from `make smoke`.
- Add JSON-schema snippets for the MCP-style tool inputs and outputs.
- Add a tagged public release after CI has passed on the public repository.
