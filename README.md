# irisslog

Alternative logging through [slog](https://pkg.go.dev/golang.org/x/exp/slog). A clone of https://github.com/iris-contrib/middleware 's zap.

This middleware will log all http requests with slog.

Current support go 1.20. So slog is `"golang.org/x/exp/slog"`


## Usage

### Start using it

Download and install it:

[//]: # (todo after pull request, change the script below to the correct path)
```sh
go get github.com/fnsne/irisslog@master
```

Import it in your code:

```go
import "github.com/fnsne/irisslog"
```
## todos
- [ ] `irisslog.New()`
- [ ] example
- [ ] skip logging
- [ ] custom fields

## Example

See the [example](_examples/example_1/main.go).

todo: add the right example

## Skip logging

When you want to skip logging for specific path,
please use `NewWithConfig`


## Custom Zap fields
Example for custom log request body, response request ID or log [Open Telemetry](https://opentelemetry.io/) TraceID.

