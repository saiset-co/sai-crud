package main

import (
	"context"
	"log"

	"go.uber.org/zap"

	"github.com/saiset-co/sai-crud/internal"
	"github.com/saiset-co/sai-crud/types"
	"github.com/saiset-co/sai-service/sai"
	"github.com/saiset-co/sai-service/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv, err := service.NewService(ctx, "./config.yaml")
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	if err := initializeComponents(); err != nil {
		sai.Logger().Error("Failed to initialize components", zap.Error(err))
		cancel()
		return
	}

	if err := srv.Start(); err != nil {
		sai.Logger().Error("Failed to start service", zap.Error(err))
		cancel()
		return
	}

	cancel()
}

func initializeComponents() error {
	config := sai.Config()
	logger := sai.Logger()

	var serviceConfig = new(types.ServiceConfig)

	err := config.GetAs("crud", serviceConfig)
	if err != nil {
		sai.Logger().Error("failed to get service config", zap.Error(err))
		return err
	}

	is := internal.NewService(serviceConfig)
	handler := internal.NewHandler(is)

	setupRoutes(handler)

	logger.Info("Components initialized successfully")
	return nil
}

func setupRoutes(handler *internal.Handler) {
	crud := sai.Router().Group("/api/v1")

	crud.POST("/", handler.Create).
		WithDoc("Create Documents", "Create multiple documents in a collection", "documents", &types.CreateRequest{}, &types.CreateResponse{})

	crud.GET("/", handler.Read).
		WithDoc("Get Documents", "Get documents with filtering and pagination. Add ?count=1 to include total count", "documents", &types.ReadRequest{}, &types.ReadResponse{})

	crud.PUT("/", handler.Update).
		WithDoc("Update Documents", "Update multiple documents by filter", "documents", &types.UpdateRequest{}, &types.UpdateResponse{})

	crud.DELETE("/", handler.Delete).
		WithDoc("Delete Documents", "Delete multiple documents by filter", "documents", &types.DeleteRequest{}, &types.DeleteResponse{})

	sai.Logger().Info("Routes configured successfully")
}
