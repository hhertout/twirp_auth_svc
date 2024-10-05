package router

import (
	"encoding/json"
	"net/http"

	"github.com/hhertout/twirp_auth/internal/hooks"
	"github.com/hhertout/twirp_auth/internal/middleware"
	"github.com/hhertout/twirp_auth/internal/repository"
	"github.com/hhertout/twirp_auth/internal/server"
	"github.com/hhertout/twirp_auth/protobuf"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
)

func GetRouter(logger *zap.Logger) *http.ServeMux {
	r, err := repository.NewUserRepository(nil)
	if err != nil {
		logger.Fatal("Error during the creation of the repository", zap.Error(err))
	}

	server := &server.AuthenticationServer{
		UserRepository: r,
	}

	handler := protobuf.NewAuthenticationServiceServer(
		server,
		twirp.WithServerPathPrefix("/api"),
		twirp.WithServerHooks(hooks.NewLoggingServerHooks(logger)),
	)

	wrapped := middleware.WithHeaders(handler)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		type Message struct {
			Message string
		}
		json.NewEncoder(w).Encode(Message{
			Message: "I'm alive",
		})
	})

	mux.Handle(handler.PathPrefix(), wrapped)

	return mux
}
