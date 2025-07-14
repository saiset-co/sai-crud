package main

import (
	"context"
	"log"

	"go.uber.org/zap"

	"github.com/saiset-co/sai-service/sai"
	"github.com/saiset-co/sai-service/service"
	saiTypes "github.com/saiset-co/sai-service/types"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := service.NewService(ctx, "./config.yaml")
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	sai.Router().GET("/hello", func(ctx *saiTypes.RequestCtx) {
		_, err := ctx.SuccessJSON("{\"hello\": \"world\"}")
		if err != nil {
			sai.Logger().Error("Failed to write response", zap.Error(err))
			return
		}
	})

	if err = srv.Start(); err != nil {
		sai.Logger().Error("Failed to start service", zap.Error(err))
		cancel()
		return
	}

	cancel()
}
