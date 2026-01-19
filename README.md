# ATC Router Go wrapper

## Building instruction

### MAKE Targets

- `clean`  
  Removes build artifacts in the repository root. It deletes files matching `libatc_router.*`.


- `build-deps`  
  Builds the Rust component located in the `atc-router` subdirectory and copies the produced shared/static library into the repository root. Effectively runs:
    - change directory to `atc-router`
    - `make clean build`
    - copy `target/release/libatc_router.*` to the root


- `test`  
  Runs Go tests for the repository with specific flags:
    - `-failfast` — stop on first failing test
    - `-count 1` — disable test caching
    - `-p 1` — use a single test package at a time
    - `-timeout 15m` — set test timeout to 15 minutes  
      Command executed: `go test -failfast -count 1 -p 1 -timeout 15m ./...`


- `all`  
  Default pipeline executed by `make` with no arguments. It runs, in order: `clean`, `build-deps`, `test`.


- `generate-header`  
  Uses `cbindgen` inside the `atc-router` crate to produce the C header file in the repository root:
    - change directory to `atc-router`
    - `cbindgen -l c > ../atc-router.h`

## Notes
- Ensure you have Rust and Cargo installed to build the Rust component.
- `atc-router` is included as git submodule. Please check out the submodule after cloning the repository:
  ```bash
  git submodule update --init --recursive
  ```
