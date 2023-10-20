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
- [x] `irisslog.New()`
- [ ] `irisslog.NewWithRecover()`
- [ ] example
- [ ] skip logging
- [ ] custom fields

## Example

See the [example](_examples/example_1/main.go).

```go
package main

import (
	"fmt"
	"github.com/fnsne/irisslog"
	"github.com/kataras/iris/v12"
	"golang.org/x/exp/slog"
	"os"
	"time"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}))

	slog.SetDefault(logger)

	app := iris.New()
	app.Use(irisslog.New(logger))

	// Example ping request.
	app.Get("/ping", func(ctx iris.Context) {
		//can get logger from context that has trace id in it.
		handlerLogger := irisslog.GetLogger(ctx)
		handlerLogger.Info("info ping")
		ctx.Text("pong " + fmt.Sprint(time.Now().Unix()))
	})

	// Example when panic happen.
	app.Get("/panic", func(ctx iris.Context) {
		panic("An unexpected error happen!")
	})

	// Listen and Server in 0.0.0.0:8080
	app.Listen(":8081")

}
```