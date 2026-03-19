//go:build atcrouter_dynamic

package goatcrouter

// #cgo LDFLAGS: -latc_router
import "C"

const atcRouterLinkInfo = "dynamic external libatc_router"
