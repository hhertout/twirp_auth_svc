package middleware

import (
	"context"
	"net/http"

	"github.com/hhertout/twirp_auth/internal/hooks"
)

func WithHeaders(base http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx = context.WithValue(ctx, hooks.ServerContextKey("route"), r.URL.Path)

		ctx = context.WithValue(ctx, hooks.ServerContextKey("user-agent"), r.Header.Get("User-Agent"))
		ctx = context.WithValue(ctx, hooks.ServerContextKey("x-forwarded-for"), r.Header.Get("X-Forwarded-For"))
		ctx = context.WithValue(ctx, hooks.ServerContextKey("remote-addr"), r.Header.Get("Remote-Addr"))

		r = r.WithContext(ctx)

		base.ServeHTTP(w, r)
	})
}
