`libatc_router_vendor` stores the native static archives that `go-atc-router`
links against at build time.

Naming convention:

- `libatc_router_darwin_arm64.a`
- `libatc_router_glibc_linux_amd64.a`
- `libatc_router_glibc_linux_arm64.a`
- `libatc_router_musl_linux_amd64.a`
- `libatc_router_musl_linux_arm64.a`

Useful commands:

- `./scripts/build-host-library.sh`
- `./scripts/import-library.sh <variant> /path/to/libatc_router.a`

The release workflow publishes variant-specific tarballs so the Linux archives
can be refreshed and imported without relying on ambiguous `linux+arch` asset
names.
