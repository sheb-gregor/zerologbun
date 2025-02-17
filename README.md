# zerologbun

A simple hook for bun that enables logging with zerolog


    go get github.com/sheb-gregor/zerologbun


## Usage

```golang
db := bun.NewDB(...)
log := zerolog.New(zerolog.NewConsoleWriter())
db.AddQueryHook(NewQueryHook(QueryHookOptions{Logger: log}))

```

### QueryHookOptions

* _LogSlow_ time.Duration value of queries considered 'slow'
* _Logger_ logger following zerolog.FieldLogger interface
* _QueryLevel_ zerolog.Level for logging queries, eg: QueryLevel: zerolog.DebugLevel
* _SlowLevel_ zerolog.Level for logging slow queries
* _ErrorLevel_ zerolog.Level for logging errors
* _MessageTemplate_ alternative message string template, avialable variables listed below
* _ErrorTemplate_ alternative error string template, available variables listed below

### Message template variables

* {{.Timestamp}} Event timestmap
* {{.Duration}} Duration of query
* {{.Query}} Query string
* {{.Operation}} Operation name (eg: SELECT, UPDATE...)
* {{.Error}} Error message if available

### Kitchen sink example
```golang
db.AddQueryHook(NewQueryHook(QueryHookOptions{
    LogSlow:    time.Second,
    Logger:     log,
    QueryLevel: zerolog.DebugLevel,
    ErrorLevel: zerolog.ErrorLevel,
    SlowLevel:  zerolog.WarnLevel,
    MessageTemplate: "{{.Operation}}[{{.Duration}}]: {{.Query}}",
    ErrorTemplate: "{{.Operation}}[{{.Duration}}]: {{.Query}}: {{.Error}",
}))

```
