#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
atc_router_src=$("${repo_root}/scripts/find-atc-router-src.sh")

variant=${ATC_ROUTER_VARIANT:-}
if [[ -z "${variant}" ]]; then
	goos=$(uname -s)
	goarch=$(uname -m)

	case "${goos}/${goarch}" in
		Darwin/arm64|Darwin/aarch64) variant="darwin_arm64" ;;
		Linux/x86_64)
			if ldd --version 2>&1 | grep -qi musl; then
				variant="musl_linux_amd64"
			else
				variant="glibc_linux_amd64"
			fi
			;;
		Linux/arm64|Linux/aarch64)
			if ldd --version 2>&1 | grep -qi musl; then
				variant="musl_linux_arm64"
			else
				variant="glibc_linux_arm64"
			fi
			;;
		*)
			echo "unsupported host platform: ${goos}/${goarch}" >&2
			exit 1
			;;
		esac
fi

cargo_args=(build --release)
release_dir="target/release"
if [[ -n "${CARGO_BUILD_TARGET:-}" ]]; then
	cargo_args+=(--target "${CARGO_BUILD_TARGET}")
	release_dir="target/${CARGO_BUILD_TARGET}/release"
fi

(
	cd "${atc_router_src}"
	cargo "${cargo_args[@]}"
)

archive="${atc_router_src}/${release_dir}/libatc_router.a"
"${repo_root}/scripts/import-library.sh" "${variant}" "${archive}"
