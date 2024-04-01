package goatcrouter

// #cgo LDFLAGS: -L/tmp/lib -latc_router
// #include <stdlib.h>
// #include "atc-router.h"
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/google/uuid"
)

// FieldType serves to indicate the desired data type for a Schema field.
type FieldType int

const (
	String FieldType = C.Type_String
	IpCidr FieldType = C.Type_IpCidr
	IpAddr FieldType = C.Type_IpAddr
	Int    FieldType = C.Type_Int
	Regex  FieldType = C.Type_Regex
)

// The Schema type holds the names and types of fields available to the router.
type Schema struct {
	s *C.Schema
}

// NewSchema creates a new empty Schema object
func NewSchema() *Schema {
	s := &Schema{s: C.schema_new()}
	runtime.SetFinalizer(s, (*Schema).Free)
	return s
}

// The Free method deallocates a Schema object
// can be called manually or automatically by the GC.
func (s *Schema) Free() {
	runtime.SetFinalizer(s, nil)
	C.schema_free(s.s)
}

// AddField is used to define fields and their associated type.
func (s *Schema) AddField(field string, typ FieldType) {
	fieldC := unsafe.Pointer(C.CString(field))
	defer C.free(fieldC)

	C.schema_add_field(s.s, (*C.schar)(fieldC), uint32(typ))
}

// The Router type holds the Matcher rules.
type Router struct {
	r *C.Router
}

// NewRouter creates a new empty Router object associated with
// the given Schema.
func NewRouter(s *Schema) *Router {
	if s == nil {
		return nil
	}

	r := &Router{r: C.router_new(s.s)}
	runtime.SetFinalizer(r, (*Router).Free)
	return r
}

// The Free method deallocates a Router object
// can be called manually or automatically by the GC.
func (r *Router) Free() {
	runtime.SetFinalizer(r, nil)
	C.router_free(r.r)
}

// AddMatcher parses a new ATC rule and adds to the Router
// under the given priority and ID.
func (r *Router) AddMatcher(priority int, id uuid.UUID, atc string) error {
	idC := unsafe.Pointer(C.CString(id.String()))
	defer C.free(idC)

	errLen := C.ulong(1024)
	errBuf := [1024]C.uchar{}
	atcC := unsafe.Pointer(C.CString(atc))
	defer C.free(atcC)

	ok := C.router_add_matcher(r.r, C.ulong(priority), (*C.schar)(idC), (*C.schar)(atcC), &errBuf[0], &errLen)
	if !ok {
		return fmt.Errorf(string(errBuf[:errLen]))
	}
	return nil
}

func (r *Router) GetFields() ([]string, error) {
	num_flds := C.router_get_fields(r.r, nil, nil)
	if num_flds == 0 {
		return nil, nil
	}

	c_flds := make([](*C.uchar), num_flds)
	c_lens := make([]C.ulong, num_flds)
	c_lens[0] = num_flds

	C.router_get_fields(r.r, &c_flds[0], &c_lens[0])

	flds := make([]string, num_flds)

	for i := range flds {
		flds[i] = C.GoStringN((*C.char)(unsafe.Pointer(c_flds[i])), (C.int)(c_lens[i]))
	}

	return flds, nil
}
