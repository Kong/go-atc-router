.PHONY: clean build-deps test generate-header

all: clean build-deps test

clean:
	rm -rf libatc_router.*

build-deps:
	cd atc-router && \
		make clean build && \
		cp target/release/libatc_router.* ../

test:
	go test -failfast -count 1 -p 1 -timeout 15m ./...

generate-header:
	cd atc-router && cbindgen -l c > ../atc-router.h
