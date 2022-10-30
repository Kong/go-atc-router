#!/bin/bash -e

DESTDIR="$1"

: ${ATC_ROUTER_REPO:=https://github.com/Kong/atc-router}
: ${ATC_ROUTER_VERSION:=main}

mkdir -p "$DESTDIR"
pushd "$DESTDIR"
  if [ ! -e "$DESTDIR/atc-router/.git" ]; then
    echo "downloading ${ATC_ROUTER_REPO}..."
    git clone "${ATC_ROUTER_REPO}" atc-router
  fi

  pushd atc-router
    echo "checking version ${ATC_ROUTER_VERSION}..."
    git checkout "${ATC_ROUTER_VERSION}"

    echo "building library..."
    make build
  popd

  echo "copying library file(s)"
  cp atc-router/target/release/libatc_router.* .
popd
