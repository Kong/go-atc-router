#!/bin/bash -e

# ATC_ROUTER_REPO=https://github.com/Kong/atc-router
ATC_ROUTER_REPO=https://github.com/javierguerragiraldez/atc-router
ATC_ROUTER_VERSION=feat/golang-binding

while [ -n "$*" ]; do
  case "$1" in
    --build ) DO_BUILD=y ;;
    --header ) DO_HEADER=y ;;
    --rm ) DO_REMOVE=y ;;
  esac
  shift
done


DESTDIR="$(realpath "$(dirname "$0")")/lib"
BUILDDIR="$(mktemp -d)"

mkdir -p "${DESTDIR}"
pushd "${BUILDDIR}"

  git clone "${ATC_ROUTER_REPO}" atc-router
  pushd atc-router
    git checkout "${ATC_ROUTER_VERSION}"

    if [ -n "$DO_BUILD" ]; then
      make build
      cp target/release/libatc_router.a "${DESTDIR}"
    fi
    if [ -n "$DO_HEADER" ]; then
      cbindgen -l c > "${DESTDIR}/atc-router.h"
    fi

  popd
popd

if [ -n "$DO_REMOVE" ]; then
  rm -rf "${BUILDDIR}"
fi
