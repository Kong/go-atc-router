package goatcrouter

// #cgo CFLAGS: -DDEFINE_ATC_ROUTER_FFI=1
// #cgo CFLAGS: -DDEFINE_ATC_ROUTER_EXPR_VALIDATION=1
// #cgo LDFLAGS: -L/tmp/lib -latc_router
// #include "atc-router.h"
import "C"

import (
	"fmt"
	"runtime"
	"slices"
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

const (
	defaultNumFields    = 100
	defaultMaxFieldSize = 100
	errBufsize          = 500
)

var (
	fieldsBuf = []C.uchar{}
	fieldsLen = C.ulong(defaultMaxFieldSize * defaultNumFields)
	fieldsNum = C.ulong(defaultNumFields)
	errorBuf  = [errBufsize]C.uchar{}
)

type ValidationResult struct {
	fields    []string
	operators uint64
	errorMsg  string
}

func ValidateExpression(s Schema, atc string) *ValidationResult {
	atcC := unsafe.Pointer(C.CString(atc))
	defer C.free(atcC)

	operators := C.ulong(0)
	errorBufLen := C.ulong(errBufsize)

loop:
	for {
		if len(fieldsBuf) < int(fieldsLen) {
			fieldsBuf = make([]C.uchar, fieldsLen)
		}

		switch C.expression_validate(
			(*C.uchar)(atcC), s.s,
			(*C.uchar)(&fieldsBuf[0]), &fieldsLen, &fieldsNum,
			&operators,
			&errorBuf[0], &errorBufLen) {
		case 0:
			break loop
		case 1:
			return &ValidationResult{
				errorMsg: C.GoStringN(
					(*C.char)(unsafe.Pointer(&errorBuf[0])),
					C.int(errorBufLen),
				),
			}
		case 2:
			continue
		}
	}
	return &ValidationResult{
		fields:    splitByNulls(fieldsBuf, int(fieldsNum)),
		operators: uint64(operators),
	}
}

const (
	BinaryOperatorFlags_EQUALS = 1 << iota
	BinaryOperatorFlags_NOT_EQUALS
	BinaryOperatorFlags_REGEX
	BinaryOperatorFlags_PREFIX
	BinaryOperatorFlags_POSTFIX
	BinaryOperatorFlags_GREATER
	BinaryOperatorFlags_GREATER_OR_EQUAL
	BinaryOperatorFlags_LESS
	BinaryOperatorFlags_LESS_OR_EQUAL
	BinaryOperatorFlags_IN
	BinaryOperatorFlags_NOT_IN
	BinaryOperatorFlags_CONTAINS
)

func splitByNulls(b []C.uchar, maxN int) []string {
	out := make([]string, maxN)
	pos := 0
	for i := range out {
		fieldLen := slices.Index(b[pos:], C.uchar(0))
		if fieldLen <= 0 {
			break
		}

		str := C.GoStringN((*C.char)(unsafe.Pointer(&b[pos])), C.int(fieldLen))

		out[i] = str
		pos += fieldLen + 1
		if pos >= len(b) {
			break
		}
	}
	slices.Sort(out)
	return out
}
