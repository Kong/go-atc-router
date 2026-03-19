//go:build linux && arm64 && musl && !atcrouter_dynamic

package goatcrouter

// #cgo LDFLAGS: ${SRCDIR}/libatc_router_vendor/libatc_router_musl_linux_arm64.a -lm -ldl -lpthread -lrt
import "C"

const atcRouterLinkInfo = "static musl_linux_arm64 from libatc_router_vendor"
