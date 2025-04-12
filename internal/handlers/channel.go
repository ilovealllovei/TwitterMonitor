package handlers

import (
	"TwitterMonitor/internal/database"
	"TwitterMonitor/internal/models"
	"TwitterMonitor/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ChannelHandler handles channel-related requests
type ChannelHandler struct {
	db *database.Database
}

// NewChannelHandler creates a new channel handler
func NewChannelHandler(db *database.Database) *ChannelHandler {
	return &ChannelHandler{db: db}
}

func (h *ChannelHandler) CreateChannel(c *gin.Context) {
	var req models.CreateOrUpdateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
		})
		return
	}

	userID, done := h.checkAbnormalInfo(c, req)
	if done {
		return
	}

	// Check if the user already has a channel
	channels, err := h.db.GetChannelsByOwnerID(userID)
	if err != nil {
		utils.LogError("Error getting channels: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to check existing channels",
		})
		return
	}

	if len(channels) > 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "User already has a channel",
		})
		return
	}

	// Create a new channel
	channel := &models.Channel{
		ID:              uuid.New().String(),
		OwnerID:         userID,
		IsVerified:      false,
		Name:            req.Name,
		Description:     req.Description,
		Avatar:          req.Avatar,
		ChatLink:        req.TwitterLink,
		IsPublic:        req.IsPublic,
		IsHot:           false,
		HotExpireAt:     "",
		Watchlist:       req.Watchlist,
		Eventlist:       req.Eventlist,
		FollowerCount:   "0",
		RecentFollowers: []int{},
	}

	if err := h.db.InsertOrUpdateChannel(channel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create channel" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    10000,
		"message": "success",
		"data": gin.H{
			"channel": channel,
		},
	})
}

func (h *ChannelHandler) UpdateChannel(c *gin.Context) {
	var req models.CreateOrUpdateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
		})
		return
	}

	userID, done := h.checkAbnormalInfo(c, req)
	if done {
		return
	}

	// Check if the user already has a channel
	channels, err := h.db.GetChannelsByOwnerID(userID)
	if err != nil {
		utils.LogError("Error getting channels: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to check existing channels",
		})
		return
	}

	if len(channels) == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "User not has a channel",
		})
		return
	}

	// Get the existing channel
	existingChannel := channels[0]

	// Update only non-empty fields
	if req.Name != "" {
		existingChannel.Name = req.Name
	}
	if req.Description != "" {
		existingChannel.Description = req.Description
	}
	if req.Avatar != "" {
		existingChannel.Avatar = req.Avatar
	}
	if req.TwitterLink != "" {
		existingChannel.ChatLink = req.TwitterLink
	}
	if req.Watchlist != nil {
		existingChannel.Watchlist = req.Watchlist
	}
	if req.Eventlist != nil {
		existingChannel.Eventlist = req.Eventlist
	}
	// Always update IsPublic since it's a boolean value
	existingChannel.IsPublic = req.IsPublic

	// Update the channel
	if err := h.db.InsertOrUpdateChannel(existingChannel); err != nil {
		utils.LogError("Error updating channel: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update channel",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    10000,
		"message": "success",
		"data": gin.H{
			"channel": existingChannel,
		},
	})

}

func (h *ChannelHandler) DeleteChannel(c *gin.Context) {
	var req models.DeleteChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
		})
		return
	}

	// Check if userID is provided
	if req.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "userID is required",
		})
		return
	}

	// Check if the channel exists and belongs to the user
	channels, err := h.db.GetChannelsByOwnerID(req.UserID)
	if err != nil {
		utils.LogError("Error getting channels: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to check channel ownership",
		})
		return
	}

	// Find the channel to delete
	var channelToDelete *models.Channel
	for _, channel := range channels {
		if channel.ID == req.ID {
			channelToDelete = channel
			break
		}
	}

	if channelToDelete == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Channel not found or you don't have permission to delete it",
		})
		return
	}

	// Delete the channel
	if err := h.db.DeleteChannel(req.ID); err != nil {
		utils.LogError("Error deleting channel: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to delete channel",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    10000,
		"message": "success",
		"data": gin.H{
			"channel": channelToDelete,
		},
	})
}

func (h *ChannelHandler) checkAbnormalInfo(c *gin.Context, req models.CreateOrUpdateChannelRequest) (int, bool) {
	// Check Watchlist length
	if len(req.Watchlist) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Watchlist exceeds the maximum limit of 100 items",
		})
		return 0, true
	}

	// Extract userID from header
	userID := req.UserID
	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "userID header is required",
		})
		return 0, true
	}
	return userID, false
}

func (h *ChannelHandler) FollowChannel(c *gin.Context) {
	var req models.FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
		})
		return
	}

	// Check if userID is provided
	if req.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "userID is required",
		})
		return
	}

	// Check if the channel exists
	channels, err := h.db.GetChannelsByOwnerID(req.UserID)
	if err != nil {
		utils.LogError("Error getting channels: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to check channel existence",
		})
		return
	}

	// Check if channel exists
	channelExists := false
	for _, channel := range channels {
		if channel.ID == req.ID {
			channelExists = true
			break
		}
	}

	if !channelExists {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Channel not found",
		})
		return
	}

	// Check if already following
	isFollowing, err := h.db.IsFollowing(req.UserID, req.ID)
	if err != nil {
		utils.LogError("Error checking follow status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to check follow status",
		})
		return
	}

	if isFollowing {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Already following this channel",
		})
		return
	}

	// Create follow relationship
	follow := &models.Follow{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		ChannelID: req.ID,
	}

	if err := h.db.FollowChannel(follow); err != nil {
		utils.LogError("Error following channel: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to follow channel",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    10000,
		"message": "success",
		"data": gin.H{
			"follow": follow,
		},
	})
}

func (h *ChannelHandler) UnfollowChannel(c *gin.Context) {
	var req models.UnfollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
		})
		return
	}

	// Check if userID is provided
	if req.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "userID is required",
		})
		return
	}

	// Check if user is following the channel
	isFollowing, err := h.db.IsFollowing(req.UserID, req.ID)
	if err != nil {
		utils.LogError("Error checking follow status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to check follow status",
		})
		return
	}

	if !isFollowing {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Not following this channel",
		})
		return
	}

	// Unfollow the channel
	if err := h.db.UnfollowChannel(req.UserID, req.ID); err != nil {
		utils.LogError("Error unfollowing channel: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to unfollow channel",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    10000,
		"message": "success",
		"data": gin.H{
			"channelId": req.ID,
			"userId":    req.UserID,
		},
	})
}

func (h *ChannelHandler) GetChannelList(c *gin.Context) {
	var req models.ChannelListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
		})
		return
	}

	// Set default values if not provided
	if req.Limit == 0 {
		req.Limit = 50
	}

	// Get channels from database
	channels, err := h.db.GetChannelsByOwnerID(req.UserID)
	if err != nil {
		utils.LogError("Error getting channels: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to get channels",
		})
		return
	}

	// Apply pagination
	total := len(channels)
	start := req.Offset
	end := start + req.Limit
	if end > total {
		end = total
	}
	if start > total {
		start = total
	}
	paginatedChannels := channels[start:end]

	// Convert channels to response format
	var responseChannels []struct {
		Meta            models.ChannelMeta `json:"meta"`
		Owner           models.Owner       `json:"owner"`
		RecentFollowers []string           `json:"recentFollowers"`
	}

	for _, channel := range paginatedChannels {
		// Convert RecentFollowers from []int to []string
		recentFollowers := make([]string, len(channel.RecentFollowers))
		for i, follower := range channel.RecentFollowers {
			recentFollowers[i] = fmt.Sprintf("%d", follower)
		}

		responseChannels = append(responseChannels, struct {
			Meta            models.ChannelMeta `json:"meta"`
			Owner           models.Owner       `json:"owner"`
			RecentFollowers []string           `json:"recentFollowers"`
		}{
			Meta: models.ChannelMeta{
				Avatar:        channel.Avatar,
				ChatLink:      channel.ChatLink,
				CreatedAt:     fmt.Sprintf("%d", channel.CreatedAt),
				Description:   channel.Description,
				IsHot:         channel.IsHot,
				HotExpireAt:   channel.HotExpireAt,
				IsVerified:    channel.IsVerified,
				Eventlist:     channel.Eventlist,
				FollowerCount: channel.FollowerCount,
				ID:            channel.ID,
				Name:          channel.Name,
				OwnerID:       channel.OwnerID,
				UpdatedAt:     fmt.Sprintf("%d", channel.UpdatedAt),
				Watchlist:     channel.Watchlist,
			},
			Owner: models.Owner{
				Email:    "", // These fields should be populated from user service
				UserID:   channel.OwnerID,
				UserName: "", // Generate a default username
			},
			RecentFollowers: recentFollowers,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    10000,
		"message": "成功",
		"data": gin.H{
			"channels": responseChannels,
			"total":    fmt.Sprintf("%d", total),
		},
	})
}

//func (h *ChannelHandler) GetChannelDetail(c *gin.Context) {
//	var req models.ChannelDetailRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, models.APIResponse{
//			Success: false,
//			Error: &models.APIError{
//				Code:    "400",
//				Message: "Invalid request format",
//			},
//		})
//		return
//	}
//
//	// Get channels for the user
//	channels, err := h.db.GetChannelsByID(req.ChannelID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, models.APIResponse{
//			Success: false,
//			Error: &models.APIError{
//				Code:    "500",
//				Message: "Failed to get channels",
//			},
//		})
//		return
//	}
//
//	// Find the specific channel
//	var channel *models.Channel
//	// 遍历所有频道
//	for _, ch := range channels {
//		// 遍历当前频道的 Watchlist
//		for _, watch := range ch.Watchlist {
//			twitterId := watch.TwitterId
//			tweets := watch.Tweets
//			twitterId := watch.TwitterId
//			profileUpdate := watch.ProfileUpdate
//			follows := watch.Follows
//			ca := watch.CA
//		}
//	}
//	// Create Twitter search request
//	twitterReq := models.TwitterUserSearchRequest{
//		Regex:      req.Regex,
//		ScreenName: req.ScreenName,
//		Page:       req.Offset + 1,
//		PageSize:   req.Limit,
//		Token:      "test0623",
//	}
//
//	// Call GetTwitterDetail
//	twitterData, err := GetTwitterDetail(&twitterReq)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, models.APIResponse{
//			Success: false,
//			Error: &models.APIError{
//				Code:    "500",
//				Message: "Failed to get Twitter data",
//			},
//		})
//		return
//	}
//
//	// Return combined response
//	c.JSON(http.StatusOK, models.APIResponse{
//		Success: true,
//		Data: map[string]interface{}{
//			"channel": channel,
//			"twitter": twitterData,
//		},
//	})
//}
