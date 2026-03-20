//go:build darwin && arm64 && !atcrouter_dynamic

package goatcrouter

// #cgo LDFLAGS: ${SRCDIR}/libatc_router_vendor/libatc_router_darwin_arm64.a
import "C"

const atcRouterLinkInfo = "static darwin_arm64 from libatc_router_vendor"
