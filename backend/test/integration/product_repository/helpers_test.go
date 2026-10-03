//go:build integration

package integration

import (
	"strings"
	"testing"

	"github.com/rs/xid"
)

func uniqueID(t *testing.T, prefix string) string {
	t.Helper()
	return prefix + "_" + xid.New().String()
}

func uniqueTitle(t *testing.T, prefix string) string {
	t.Helper()
	name := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	return prefix + "_" + name + "_" + xid.New().String()
}
