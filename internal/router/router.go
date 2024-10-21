package router

import (
	"encoding/json"
	"net/http"

	"github.com/hhertout/twirp_auth/internal/hooks"
	"github.com/hhertout/twirp_auth/internal/middleware"
	"github.com/hhertout/twirp_auth/internal/repository"
	"github.com/hhertout/twirp_auth/internal/server"
	"github.com/hhertout/twirp_auth/lib/crypto"
	"github.com/hhertout/twirp_auth/pkg/auth"
	"github.com/hhertout/twirp_auth/protobuf/proto_auth"
	"github.com/hhertout/twirp_auth/protobuf/proto_user"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
)

func GetRouter(logger *zap.Logger) *http.ServeMux {
	r, err := repository.NewUserRepository(nil)
	if err != nil {
		logger.Fatal("Error during the creation of the repository", zap.Error(err))
	}

	auth_server := &server.AuthenticationServer{
		Logger:          logger,
		UserRepository:  r,
		PasswordService: crypto.NewPasswordService(),
		JwtService:      crypto.NewJWTService(),
	}

	user_server := &server.UserServer{
		Logger:          logger,
		UserRepository:  r,
		PasswordService: crypto.NewPasswordService(),
		JwtService:      crypto.NewJWTService(),
		AuthManager:     auth.NewAuthManager(r),
	}

	auth_handler := proto_auth.NewAuthenticationServiceServer(
		auth_server,
		twirp.WithServerPathPrefix("/api"),
		twirp.WithServerHooks(hooks.NewLoggingServerHooks(logger)),
	)

	user_handler := proto_user.NewUserServiceServer(
		user_server,
		twirp.WithServerPathPrefix("/api"),
		twirp.WithServerHooks(hooks.NewLoggingServerHooks(logger)),
	)

	wrapped_auth := middleware.WithHeaders(auth_handler)
	wrapped_user := middleware.WithHeaders(user_handler)

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

	mux.Handle(auth_handler.PathPrefix(), wrapped_auth)
	mux.Handle(user_handler.PathPrefix(), wrapped_user)

	return mux
}
