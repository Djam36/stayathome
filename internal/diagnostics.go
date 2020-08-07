package internal

import (
	"net"
	"net/http"

	muxtrace "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux"
	oteltrace "go.opentelemetry.io/otel/api/trace"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

// Diagnostics responsible for diagnostics logic of the app
func Diagnostics(logger *zap.SugaredLogger, tracer oteltrace.Tracer, port string, shutdown chan<- error) *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/health", handleHealth(logger.With("handler", "health")))

	mw := muxtrace.Middleware("diags", muxtrace.WithTracer(tracer))
	r.Use(mw)

	server := http.Server{
		Addr:    net.JoinHostPort("", port),
		Handler: r,
	}

	logger.Info("Ready to start the server...")
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			shutdown <- err
		}
	}()

	return &server
}

func handleHealth(logger *zap.SugaredLogger) func(http.ResponseWriter, *http.Request) {
	return func(
		w http.ResponseWriter, r *http.Request) {
		logger.Info("Received a call")
		w.WriteHeader(http.StatusOK)
	}
}
