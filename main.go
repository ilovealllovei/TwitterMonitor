package main

import (
	"TwitterMonitor/config"
	"TwitterMonitor/internal/database"
	"TwitterMonitor/internal/handlers"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting the application...")
	cfg := config.LoadConfig()
	log.Println("Config loaded:", cfg)

	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	log.Println("Database connected successfully")

	// Initialize handlers
	channelHandler := handlers.NewChannelHandler(db)
	log.Println("Channel handler initialized")

	// Initialize Gin router
	router := gin.Default()
	log.Println("Gin router initialized")

	// API routes
	api := router.Group("/v1")
	{
		channel := api.Group("/channel")
		{
			channel.POST("/create", channelHandler.CreateChannel)
			channel.POST("/update", channelHandler.UpdateChannel)
			channel.POST("/delete", channelHandler.DeleteChannel)
			channel.POST("/follow", channelHandler.FollowChannel)
			channel.POST("/unfollow", channelHandler.UnfollowChannel)
			channel.GET("/channel_list", channelHandler.GetChannelList)
			channel.GET("/channel_content", channelHandler.GetChannelContent)
			channel.GET("/channel_ws", channelHandler.ChannelWSHandler)
		}
	}

	// Start server
	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
