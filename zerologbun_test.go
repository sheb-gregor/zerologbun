package zerologbun

import (
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

func TestLogging(t *testing.T) {
	w := zerolog.NewConsoleWriter()
	w.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%-6s", i)
	}
	var log = zerolog.New(w).
		Level(zerolog.DebugLevel)

	db := bun.DB{}
	db.AddQueryHook(NewQueryHook(QueryHookOptions{Logger: log}))
	// @TODO against empty db (stub?)
}
