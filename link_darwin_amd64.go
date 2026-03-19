//go:build darwin && amd64 && !atcrouter_dynamic

package goatcrouter

// #error "darwin/amd64 does not have a vendored libatc_router archive yet; import one with scripts/import-library.sh or build with -tags atcrouter_dynamic"
import "C"
