package handlers

import (
	"TwitterMonitor/internal/database"
	"TwitterMonitor/internal/models"
	"TwitterMonitor/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

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
		ChatLink:        req.ChatLink,
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

	if userID != channels[0].OwnerID {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "User not permitted to change channel",
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
	if req.ChatLink != "" {
		existingChannel.ChatLink = req.ChatLink
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

	if req.UserID != channels[0].OwnerID {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "User not permitted to change channel",
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
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "400",
				Message: "Invalid request format: " + err.Error(),
			},
		})
		return
	}

	// Set default values for limit and offset
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	var channels []*models.Channel
	var err error

	switch req.Type {
	case "1":
		// Get channels by owner ID
		channels, err = h.db.GetChannelsByOwnerID(req.UserID)
	case "2":
		// Get followed channels
		follows, err := h.db.GetFollowedChannels(req.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "500",
					Message: "Failed to get followed channels",
				},
			})
			return
		}

		// Collect channel IDs
		var channelIDs []string
		for _, follow := range follows {
			channelIDs = append(channelIDs, follow.ChannelID)
		}

		// Get all channels in a single query
		channels, err = h.db.GetChannelByIDs(channelIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "500",
					Message: "Failed to get channels",
				},
			})
			return
		}
	default:
		// Get all channels
		channels, err = h.db.GetAllChannels(req.Limit, req.Offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "500",
				Message: "Failed to get channels",
			},
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
		RecentFollowers []int              `json:"recentFollowers"`
	}

	for _, channel := range paginatedChannels {

		responseChannels = append(responseChannels, struct {
			Meta            models.ChannelMeta `json:"meta"`
			Owner           models.Owner       `json:"owner"`
			RecentFollowers []int              `json:"recentFollowers"`
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
			RecentFollowers: channel.RecentFollowers,
		})
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"channels": responseChannels,
			"total":    total,
		},
	})
}

// fetchMarketInfo fetches market info for a given chain ID and token CA
func fetchMarketInfo(chainId, tokenCa string) (interface{}, error) {
	if chainId == "" || tokenCa == "" {
		return nil, nil
	}

	url := fmt.Sprintf("https://api.litrocket.io/v1/market/market_info?chain_id=%s&token_ca=%s", chainId, tokenCa)
	log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add required headers
	req.Header.Set("X-Language", "zh")
	req.Header.Set("X-Source", "ios")
	req.Header.Set("Qlbl69aq2dxo4t", "1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Data interface{} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (h *ChannelHandler) GetChannelContent(c *gin.Context) {
	var req models.ChannelContentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "400",
				Message: "Invalid request format: " + err.Error(),
			},
		})
		return
	}

	// Set default values for limit and offset
	if req.Limit <= 0 {
		req.Limit = 50 // Default limit
	}
	if req.Offset < 0 {
		req.Offset = 0 // Default offset
	}

	// Get channels for the user
	channels, err := h.db.GetChannelsByID(req.ChannelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "500",
				Message: "Failed to get channels",
			},
		})
		return
	}

	if len(channels) == 0 {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "404",
				Message: "Channel not found",
			},
		})
		return
	}

	channel := channels[0]
	var twitterInfos []*models.TwitterInfo

	if req.ContentType == 1 {
		// Build conditions for tweets
		var conditions []string
		for _, watch := range channel.Watchlist {
			log.Printf("watch: %+v", watch)
			if !watch.Tweets {
				continue
			}
			condition := fmt.Sprintf("(twitterId = '%s'", watch.TwitterId)
			if watch.CA != "" {
				condition += fmt.Sprintf(" AND address = '%s'", watch.CA)
			}
			condition += ")"
			conditions = append(conditions, condition)
		}

		twitterInfos, err = h.db.GetTwitterInfoByWatchlist(conditions, req.ContentType, req.Limit, req.Offset)
		if err != nil {
			utils.LogError("Failed to get Twitter info: %v", err)
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "500",
					Message: "Failed to get Twitter info",
				},
			})
			return
		}
	} else if req.ContentType == 2 {
		// Build conditions for profile updates and follows
		var twitterIds []string
		for _, watch := range channel.Watchlist {
			if watch.ProfileUpdate && watch.Follows {
				twitterIds = append(twitterIds, watch.TwitterId)
			}
		}

		twitterInfos, err = h.db.GetTwitterInfoByProfileAndFollow(twitterIds, req.ContentType, req.Limit, req.Offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "500",
					Message: "Failed to get Twitter info",
				},
			})
			return
		}
	}

	// Get market info for each Twitter info
	marketInfos := make([]interface{}, len(twitterInfos))
	for i, info := range twitterInfos {
		marketInfo, err := fetchMarketInfo(info.ChainId, info.Address)

		if err != nil {
			utils.LogError("Failed to fetch market info: %v", err)
			marketInfos[i] = nil
			continue
		}
		marketInfos[i] = marketInfo
	}

	// Return combined response
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"twitter": twitterInfos,
			"market":  marketInfos,
			"total":   len(channels),
		},
	})
}

func (h *ChannelHandler) TwitterInfo(c *gin.Context) {
	user := c.Query("user")
	if user == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "400",
				Message: "user parameter is required",
			},
		})
		return
	}

	url := fmt.Sprintf("http://43.160.199.161:5188/tw_user_info?user=%s&token=test0623", user)
	resp, err := http.Get(url)
	if err != nil {
		utils.LogError("Failed to call Twitter info API: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "500",
				Message: "Failed to fetch Twitter info",
			},
		})
		return
	}
	defer resp.Body.Close()

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		utils.LogError("Failed to decode response: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "500",
				Message: "Failed to parse response",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
	})
}
