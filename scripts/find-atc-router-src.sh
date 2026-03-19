#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)

candidates=()
if [[ -n "${ATC_ROUTER_SRC:-}" ]]; then
	if [[ "${ATC_ROUTER_SRC}" = /* ]]; then
		candidates+=("${ATC_ROUTER_SRC}")
	else
		candidates+=("${repo_root}/${ATC_ROUTER_SRC}")
	fi
fi

candidates+=(
	"${repo_root}/atc-router"
	"${repo_root}/../atc-router"
)

for candidate in "${candidates[@]}"; do
	if [[ -f "${candidate}/Cargo.toml" ]]; then
		cd "${candidate}"
		pwd
		exit 0
	fi
done

echo "unable to locate an atc-router checkout; set ATC_ROUTER_SRC to a directory containing Cargo.toml" >&2
exit 1
