package goatcrouter

import (
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
