package format

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorsRetainCauseAndOperation(t *testing.T) {
	cause := errors.New("failure")
	err := Error("operation", cause)
	require.ErrorIs(t, err, cause)
	require.Contains(t, err.Error(), "operation")
}

func TestStructFormattingAndNilError(t *testing.T) {
	require.NoError(t, Error("op", nil))
	require.JSONEq(t, `{"title":"Desk"}`, Struct(struct {
		Title string `json:"title"`
	}{Title: "Desk"}))
	require.NotEmpty(t, Struct(make(chan int)))
}
