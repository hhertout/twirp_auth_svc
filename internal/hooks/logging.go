package hooks

import (
	"context"
	"time"

	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
)

type ServerContextKey string

func NewLoggingServerHooks(l *zap.Logger) *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			ctx = context.WithValue(ctx, ServerContextKey("time_taken"), time.Now())
			return ctx, nil
		},
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			return ctx, nil
		},
		Error: func(ctx context.Context, twerr twirp.Error) context.Context {
			l.Sugar().Errorf("Error code=%s; message= %s", string(twerr.Code()), twerr.Msg())
			return ctx
		},
		ResponseSent: func(ctx context.Context) {
			clientIP, _ := ctx.Value(ServerContextKey("x-forwarded-for")).(string)
			if clientIP == "" {
				clientIP, _ = ctx.Value(ServerContextKey("remote-addr")).(string)
			}
			if clientIP == "" {
				clientIP = "unknown"
			}

			userAgent, _ := ctx.Value(ServerContextKey("user-agent")).(string)
			route, _ := ctx.Value(ServerContextKey("route")).(string)
			method, _ := twirp.MethodName(ctx)
			service, _ := twirp.ServiceName(ctx)
			status, _ := twirp.StatusCode(ctx)
			start, _ := ctx.Value(ServerContextKey("time_taken")).(time.Time)

			logFields := []zap.Field{
				zap.String("route", route),
				zap.String("user_agent", userAgent),
				zap.String("client_ip", clientIP),
				zap.String("status", status),
				zap.String("service", service),
				zap.String("method", method),
				zap.Duration("time_taken", time.Since(start)),
			}

			l.Info("GRPC Request", logFields...)
		},
	}
}
