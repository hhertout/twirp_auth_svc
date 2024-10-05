package router

import (
	"encoding/json"
	"net/http"

	pb "github.com/hhertout/twirp_example/generated"
	"github.com/hhertout/twirp_example/internal/hooks"
	"github.com/hhertout/twirp_example/internal/server"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
)

func GetRouter(logger *zap.Logger) *http.ServeMux {
	server := &server.Server{}
	twirpHandler := pb.NewHaberdasherServer(
		server,
		twirp.WithServerPathPrefix("/api"),
		twirp.WithServerHooks(hooks.NewLoggingServerHooks(logger)),
	)

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

	mux.Handle(twirpHandler.PathPrefix(), twirpHandler)

	return mux
}
