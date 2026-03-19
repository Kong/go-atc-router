.PHONY: all clean build-deps test build-host-lib import-lib generate-header

all: test

clean:
	rm -f libatc_router.*
	@if [ -n "$${ATC_ROUTER_SRC:-}" ] && [ -f "$${ATC_ROUTER_SRC}/Makefile" ]; then \
		$(MAKE) -C "$${ATC_ROUTER_SRC}" clean; \
	elif [ -f ./atc-router/Makefile ]; then \
		$(MAKE) -C ./atc-router clean; \
	elif [ -f ../atc-router/Makefile ]; then \
		$(MAKE) -C ../atc-router clean; \
	fi

build-deps:
	@ATC_ROUTER_SRC="$$(./scripts/find-atc-router-src.sh)"; \
		$(MAKE) -C "$${ATC_ROUTER_SRC}" clean build; \
		cp "$${ATC_ROUTER_SRC}"/target/release/libatc_router.* .

test:
	go test -failfast -count 1 -p 1 -timeout 15m ./...

build-host-lib:
	./scripts/build-host-library.sh

import-lib:
	@test -n "$(VARIANT)" || (echo "VARIANT is required" >&2; exit 1)
	@test -n "$(SRC)" || (echo "SRC is required" >&2; exit 1)
	./scripts/import-library.sh "$(VARIANT)" "$(SRC)"

generate-header:
	@ATC_ROUTER_SRC="$$(./scripts/find-atc-router-src.sh)"; \
		cd "$${ATC_ROUTER_SRC}" && cbindgen -l c > "$(CURDIR)/atc-router.h"
