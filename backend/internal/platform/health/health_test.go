package health

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLivenessAndReadiness(t *testing.T) {
	require.Equal(t, "ok", Live().Status)
	s, code := Ready(nil)
	require.Equal(t, "ready", s.Status)
	require.Equal(t, 200, code)
	s, code = Ready(errors.New("database down"))
	require.Equal(t, "not ready", s.Status)
	require.Equal(t, 503, code)
}
