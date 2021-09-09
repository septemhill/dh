package middleware

import (
	"context"
	"net/http"
	"strings"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/septemhill/dh/iplimiter"
)

func IPLimiterMiddleware(rps int) httptransport.RequestFunc {
	limiter := iplimiter.NewLimiter(iplimiter.RPS(rps))
	return func(ctx context.Context, req *http.Request) context.Context {
		off := strings.LastIndex(req.RemoteAddr, ":")
		l := limiter.GetLimiter(req.RemoteAddr[0:off])
		l.Take()
		return ctx
	}
}
