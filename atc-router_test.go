package goatcrouter

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func verify(atc string) error {
	schema := NewSchema()
	defer schema.Free()

	schema.AddField("http.path", String)
	schema.AddField("tcp.port", Int)

	router := NewRouter(schema)
	defer router.Free()

	return router.AddMatcher(1, uuid.New(), atc)
}

func get_fields(atc string) ([]string, error) {
	schema := NewSchema()
	defer schema.Free()

	schema.AddField("http.path", String)
	schema.AddField("tcp.port", Int)

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
