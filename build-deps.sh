#!/bin/bash -e

# ATC_ROUTER_REPO=https://github.com/Kong/atc-router
ATC_ROUTER_REPO=https://github.com/javierguerragiraldez/atc-router
ATC_ROUTER_VERSION=feat/golang-binding

while [ -n "$*" ]; do
  case "$1" in
    --build ) DO_BUILD=y ;;
    --header ) DO_HEADER=y ;;
    --install=*) INSTALL_DEST="${1#--install=}" ;;
    --install ) INSTALL_DEST="/usr/local/lib" ;;
    --cache ) DO_CACHE=y ;;
    --rm ) DO_REMOVE=y ;;
  esac
  shift
done

DESTDIR="${DESTDIR:-$PWD}"
BUILDDIR="$(mktemp -d)"
LIBNAME="target/release/libatc_router.so"

mkdir -p "${DESTDIR}"
pushd "${BUILDDIR}"

  git clone "${ATC_ROUTER_REPO}" atc-router
  pushd atc-router
    git checkout "${ATC_ROUTER_VERSION}"

    if [ -n "$DO_BUILD" ]; then
      make build
    fi

    if [ -n "$DO_CACHE" -a -e "$LIBNAME" ]; then
      cp "$LIBNAME" /tmp/
    fi

    if [ -n "$INSTALL_DEST" -a -e "$LIBNAME" ]; then
      sudo install "$LIBNAME" "$INSTALL_DEST"
    fi

    if [ -n "$DO_HEADER" ]; then
      cbindgen -l c > "${DESTDIR}/atc-router.h"
    fi

  popd
popd

if [ -n "$DO_REMOVE" ]; then
  rm -rf "${BUILDDIR}"
fi
