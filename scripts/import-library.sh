#!/usr/bin/env bash

set -euo pipefail

variant=${1:-}
src=${2:-}

if [[ -z "${variant}" || -z "${src}" ]]; then
	echo "usage: $0 <variant> <path-to-libatc_router.a>" >&2
	exit 1
fi

case "${variant}" in
	darwin_arm64|glibc_linux_amd64|glibc_linux_arm64|musl_linux_amd64|musl_linux_arm64)
		;;
	*)
		echo "unsupported variant: ${variant}" >&2
		exit 1
		;;
esac

if [[ ! -f "${src}" ]]; then
	echo "archive not found: ${src}" >&2
	exit 1
fi

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
dest_dir="${repo_root}/libatc_router_vendor"
dest="${dest_dir}/libatc_router_${variant}.a"

mkdir -p "${dest_dir}"
cp "${src}" "${dest}"

echo "imported ${src} -> ${dest}"
