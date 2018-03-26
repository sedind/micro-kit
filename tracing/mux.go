package tracing

import (
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
)

// NewServeMux creates a new TracedServeMux
func NewServeMux(tracer opentracing.Tracer) *TracedServeMux {
	return &TracedServeMux{
		mux:    http.NewServeMux(),
		tracer: tracer,
	}
}

// TracedServeMux created a new TracedServeMux
type TracedServeMux struct {
	mux    *http.ServeMux
	tracer opentracing.Tracer
}

// Handle implements http.ServeMux#Handle
func (tm *TracedServeMux) Handle(pattern string, handler http.Handler) {
	middleware := nethttp.Middleware(
		tm.tracer,
		handler,
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return "HTTP " + r.Method + " " + pattern
		}))
	tm.mux.Handle(pattern, middleware)
}

func (tm *TracedServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm.mux.ServeHTTP(w, r)
}
