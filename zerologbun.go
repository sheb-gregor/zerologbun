package zerologbun

import (
	"bytes"
	"context"
	"database/sql"
	"strings"
	"text/template"
	"time"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

// QueryHookOptions logging options
type QueryHookOptions struct {
	LogSlow         time.Duration
	Logger          zerolog.Logger
	QueryLevel      zerolog.Level
	SlowLevel       zerolog.Level
	ErrorLevel      zerolog.Level
	MessageTemplate string
	ErrorTemplate   string
}

// QueryHook wraps query hook
type QueryHook struct {
	opts            QueryHookOptions
	errorTemplate   *template.Template
	messageTemplate *template.Template
}

// LogEntryVars variables made available to template
type LogEntryVars struct {
	Timestamp time.Time
	Query     string
	Operation string
	Duration  time.Duration
	Error     error
}

// NewQueryHook returns new instance
func NewQueryHook(opts QueryHookOptions) *QueryHook {
	h := new(QueryHook)

	if opts.ErrorTemplate == "" {
		opts.ErrorTemplate = "{{.Operation}}: {{.Query}}: {{.Error}}"
	}
	if opts.MessageTemplate == "" {
		opts.MessageTemplate = "{{.Operation}}: {{.Query}}"
	}

	h.opts = opts
	errorTemplate, err := template.New("ErrorTemplate").Parse(h.opts.ErrorTemplate)
	if err != nil {
		panic(err)
	}

	messageTemplate, err := template.New("MessageTemplate").Parse(h.opts.MessageTemplate)
	if err != nil {
		panic(err)
	}

	h.errorTemplate = errorTemplate
	h.messageTemplate = messageTemplate
	return h
}

// BeforeQuery does nothing tbh
func (h *QueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

// AfterQuery convert a bun QueryEvent into a zerolog message
func (h *QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	var isError bool
	var msg bytes.Buffer

	now := time.Now()
	dur := now.Sub(event.StartTime)

	level := zerolog.Disabled

	switch event.Err {
	case nil, sql.ErrNoRows:
		isError = false
		if h.opts.LogSlow > 0 && dur >= h.opts.LogSlow {
			level = h.opts.SlowLevel
		} else {
			level = h.opts.QueryLevel
		}
	default:
		isError = true
		level = h.opts.ErrorLevel
	}

	if level == zerolog.Disabled {
		return
	}

	args := &LogEntryVars{
		Timestamp: now,
		Query:     event.Query,
		Operation: eventOperation(event),
		Duration:  dur,
		Error:     event.Err,
	}

	if isError {
		if err := h.errorTemplate.Execute(&msg, args); err != nil {
			panic(err)
		}
	} else {
		if err := h.messageTemplate.Execute(&msg, args); err != nil {
			panic(err)
		}
	}

	h.opts.Logger.WithLevel(level).
		Stringer("duration", args.Duration).
		Msg(msg.String())
}

// taken from bun
func eventOperation(event *bun.QueryEvent) string {
	switch event.QueryAppender.(type) {
	case *bun.SelectQuery:
		return "SELECT"
	case *bun.InsertQuery:
		return "INSERT"
	case *bun.UpdateQuery:
		return "UPDATE"
	case *bun.DeleteQuery:
		return "DELETE"
	case *bun.CreateTableQuery:
		return "CREATE TABLE"
	case *bun.DropTableQuery:
		return "DROP TABLE"
	}
	return queryOperation(event.Query)
}

// taken from bun
func queryOperation(name string) string {
	if idx := strings.Index(name, " "); idx > 0 {
		name = name[:idx]
	}
	if len(name) > 16 {
		name = name[:16]
	}
	return name
}
