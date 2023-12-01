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

func Test_Verify(t *testing.T) {
	require.NoError(t, verify("tcp.port == 1"))
	require.Error(t, verify("bad.var == 9"))
}

func Test_AddMatcher_Neg_Priority(t *testing.T) {
	schema := NewSchema()
	defer schema.Free()

	schema.AddField("http.path", String)
	schema.AddField("tcp.port", Int)

	router := NewRouter(schema)
	defer router.Free()
	require.Error(t, router.AddMatcher(-123, uuid.New(), "tcp.port == 1"))
}
