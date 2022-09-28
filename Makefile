
# ATC_ROUTER_REPO=https://github.com/Kong/atc-router
ATC_ROUTER_REPO=https://github.com/javierguerragiraldez/atc-router
ATC_ROUTER_VERSION=feat/golang-binding

LIBRARY=lib/libatc_router.a
HEADER=lib/atc-router.h

.PHONY: clean

all: $(LIBRARY) $(HEADER)

lib/atc-router/Makefile:
	mkdir -p lib/
	cd lib/ && git clone $(ATC_ROUTER_REPO)
	cd lib/atc-router &&git checkout $(ATC_ROUTER_VERSION)

$(LIBRARY): lib/atc-router/Makefile
	cd lib/atc-router && make build
	cp lib/atc-router/target/release/libatc_router.a $(LIBRARY)

$(HEADER): lib/atc-router/cbindgen.toml
	cd lib/atc-router && cbindgen -l c > ../atc-router.h

clean:
	rm -rf lib/
