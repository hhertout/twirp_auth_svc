package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/hhertout/twirp_example/internal/router"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	if os.Getenv("GO_ENV") == "development" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	router := router.GetRouter(logger)

	if os.Getenv("GO_ENV") == "development" {
		logger.Sugar().Info("⚠️ Caution : The server will be running under development mode 🔨🔨")
	}

	logger.Sugar().Info("🚀 Server running on port %d ! \n", port)
	http.ListenAndServe(fmt.Sprintf("%v:%v", "0.0.0.0", port), router)
}
