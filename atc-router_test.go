package goatcrouter

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type field struct {
	f string
	t FieldType
}

func newSchema(fields ...field) *Schema {
	schema := NewSchema()
	for _, fld := range fields {
		schema.AddField(fld.f, fld.t)
	}
	return schema
}

func verify(atc string) error {
	schema := newSchema(field{"http.path", String}, field{"tcp.port", Int})
	defer schema.Free()

	router := NewRouter(schema)
	defer router.Free()

	return router.AddMatcher(1, uuid.New(), atc)
}

func get_fields(atc string) ([]string, error) {
	schema := newSchema(field{"http.path", String}, field{"tcp.port", Int})
	defer schema.Free()

	router := NewRouter(schema)
	defer router.Free()

	if err := router.AddMatcher(1, uuid.New(), atc); err != nil {
		return nil, err
	}

	flds, err := router.GetFields()
	return flds, err
}

func Test_Verify(t *testing.T) {
	require.NoError(t, verify("tcp.port == 1"))
	require.Error(t, verify("bad.var == 9"))
}

func Test_GetFields(t *testing.T) {
	fields, err := get_fields("tcp.port == 1")
	require.NoError(t, err)
	require.ElementsMatch(t, []any{"tcp.port"}, fields)

	fields, err = get_fields(`http.path==""`)
	require.NoError(t, err)
	require.ElementsMatch(t, []any{"http.path"}, fields)

	fields, err = get_fields(`tcp.port == 1 && http.path==""`)
	require.NoError(t, err)
	require.ElementsMatch(t, []any{"http.path", "tcp.port"}, fields)
}

func Test_splitByNulls(t *testing.T) {
	for _, test := range []struct {
		name string
		buf  []byte
		n    int
		out  []string
	}{
		{name: "empty input", buf: nil, out: []string{}},
		{name: "empty input, n:1", buf: nil, n: 1, out: []string{""}},
		{name: "single string", buf: []byte("one\x00"), n: 1, out: []string{"one"}},
		{name: "single string, n:0", buf: []byte("one\x00"), n: 0, out: []string{}},
		{name: "single string, n:2", buf: []byte("one\x00"), n: 2, out: []string{"", "one"}},
		{name: "unterminated", buf: []byte("one"), n: 1, out: []string{""}},
		{name: "two, unterminated", buf: []byte("one\x00two"), n: 2, out: []string{"", "one"}},
		{name: "consecutive separator 3", buf: []byte("one\x00\x00two\x00"), n: 3, out: []string{"", "one", "two"}},
		{name: "consecutive separator 2", buf: []byte("one\x00\x00two\x00"), n: 2, out: []string{"", "one"}},
	} {
		t.Run(fmt.Sprintf("splitByNulls %q:", test.name), func(t *testing.T) {
			out := splitByNulls(test.buf, test.n)
			require.Equal(t, test.n, len(out))
			require.EqualValues(t, test.out, out)
		})
	}
}

func Test_ValidateExpression(t *testing.T) {
	schema := newSchema(field{"http.path", String}, field{"tcp.port", Int})
	defer schema.Free()

	res := ValidateExpression(*schema, "tcp.port == 1")
	require.Equal(t, &ValidationResult{
		fields:    []string{"tcp.port"},
		operators: BinaryOperatorFlags_EQUALS,
	}, res)

	require.Equal(t, &ValidationResult{
		fields:    []string{"http.path", "tcp.port"},
		operators: BinaryOperatorFlags_EQUALS + BinaryOperatorFlags_NOT_EQUALS,
	}, ValidateExpression(*schema, `tcp.port == 1 && http.path!=""`))
}
