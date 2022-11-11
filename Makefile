
# ATC_ROUTER_REPO=https://github.com/Kong/atc-router
ATC_ROUTER_REPO=https://github.com/javierguerragiraldez/atc-router
ATC_ROUTER_VERSION=feat/golang-binding

LIBRARY=lib/libatc_router.a
HEADER=lib/atc-router.h

.PHONY: clean build-deps

all: build-deps

build-deps:
	./build-deps.sh --build --header --rm

clean:
	rm -rf lib/
