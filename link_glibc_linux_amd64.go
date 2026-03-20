//go:build linux && amd64 && !musl && !atcrouter_dynamic

package goatcrouter

// #cgo LDFLAGS: ${SRCDIR}/libatc_router_vendor/libatc_router_glibc_linux_amd64.a -lm -ldl -lpthread -lrt
import "C"

const atcRouterLinkInfo = "static glibc_linux_amd64 from libatc_router_vendor"
