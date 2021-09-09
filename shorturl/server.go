package shorturl

import (
	"context"
	"fmt"
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/septemhill/dh/shorturl/endpoints"
	"github.com/septemhill/dh/shorturl/middleware"
	"github.com/septemhill/dh/shorturl/transport"
)

type server struct {
	port    int
	handler http.Handler
}

func NewServer(port int, ep endpoints.Endpoints) *server {
	r := mux.NewRouter()

	r.Methods("POST").Path("/api/v1/urls").Handler(httptransport.NewServer(
		ep.UploadURL,
		transport.DecodeUploadRequest,
		httptransport.EncodeJSONResponse,
		httptransport.ServerBefore(middleware.IPLimiterMiddleware(1)),
	))
	r.Methods("GET").Path("/{shortkey}").Handler(httptransport.NewServer(
		ep.AccessURL,
		transport.DecodeAccessURLRequest,
		transport.EncodeAccessURLResponse,
	))

	return &server{
		port:    port,
		handler: r,
	}
}

func (s *server) Run(ctx context.Context) error {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.handler,
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(ctx)
	}()

	log.Printf("Short URL service starting on :%d ...\n", s.port)
	return srv.ListenAndServe()
}
