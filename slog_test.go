package irisslog_test

import (
	"context"
	"github.com/fnsne/irisslog"
	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareSuite))
}

type MiddlewareSuite struct {
	suite.Suite
	app         *iris.Application
	mockHandler *observeHandler
}

var checkNotEmpty = "skip check value"

func (suite *MiddlewareSuite) SetupTest() {
	suite.mockHandler = newObserveHandler(suite.T())
	logger := slog.New(suite.mockHandler)

	suite.app = iris.New()
	suite.app.Use(irisslog.New(logger))

	suite.app.Get("/test", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusNoContent)
		_ = ctx.JSON(nil)
	})

	err := suite.app.Build()
	require.NoError(suite.T(), err)

	//init mock handler for slog init
	suite.mockHandler.On("Enabled", mock.Anything, mock.Anything).Return(true)
	suite.mockHandler.On("Handle", mock.Anything, mock.Anything).Return(nil)
}

func (suite *MiddlewareSuite) Test_middleware_with_success_request() {
	suite.doRequest(http.MethodGet, "/test")
	suite.mockHandler.AssertNumberOfCalls(suite.T(), "Handle", 1)

	assert.Equal(suite.T(), 1, len(suite.mockHandler.handleRecords), "")
	suite.logMessageShouldBe("Success", []slog.Attr{
		slog.Any("status", iris.StatusNoContent),
		slog.Any("method", http.MethodGet),
		slog.Any("path", "/test"),
		slog.Any("query", ""),
		slog.Any("latency", checkNotEmpty),
	})
}

func (suite *MiddlewareSuite) logMessageShouldBe(message string, expectAttrs []slog.Attr) {
	record := suite.mockHandler.handleRecords[0]
	assert.Equal(suite.T(), slog.LevelInfo, record.Level, "")
	assert.Equal(suite.T(), message, record.Message, "")
	//get attrs from record
	var attrs []slog.Attr
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})
	assert.Equal(suite.T(), len(expectAttrs), len(attrs), "")
	for i := 0; i < len(expectAttrs); i++ {
		expect := expectAttrs[i]
		actual := attrs[i]
		assert.Equal(suite.T(), expect.Key, actual.Key, "")
		if expect.Value.String() == checkNotEmpty {
			assert.NotEmpty(suite.T(), actual.Value.String(), "")
		} else {
			assert.Equal(suite.T(), expect.Value.String(), actual.Value.String(), "")
		}
	}
}

func (suite *MiddlewareSuite) doRequest(method string, url string) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, url, nil)
	suite.app.ServeHTTP(w, r)
}

type observeHandler struct {
	mock.Mock
	t             *testing.T
	handleRecords []slog.Record
}

func newObserveHandler(t *testing.T) *observeHandler {
	return &observeHandler{t: t}
}

func (o *observeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	arguments := o.Called(ctx, level)
	return arguments.Bool(0)
}

func (o *observeHandler) Handle(ctx context.Context, record slog.Record) error {
	arguments := o.Called(ctx, record)
	o.handleRecords = append(o.handleRecords, record)
	return arguments.Error(0)
}

func (o *observeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return o.Called(attrs).Get(0).(slog.Handler)
}

func (o *observeHandler) WithGroup(name string) slog.Handler {
	return o.Called(name).Get(0).(slog.Handler)
}
