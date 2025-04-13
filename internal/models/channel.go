package models

// ChannelType represents the type of channel
type ChannelType string

const (
	ChannelTypeTwitter ChannelType = "X_TWITTER"
)

// AccountConfig represents the configuration for a monitored Twitter account
type AccountConfig struct {
	ShowTweets         bool `json:"showTweets"`
	ShowProfileUpdates bool `json:"showProfileUpdates"`
	ShowFollows        bool `json:"showFollows"`
	FilterCA           bool `json:"filterCA"`
}

// MonitoredAccount represents a Twitter account being monitored
type MonitoredAccount struct {
	Username    string        `json:"username"`
	DisplayName string        `json:"displayName"`
	IsVerified  bool          `json:"isVerified"`
	Config      AccountConfig `json:"config"`
}

// Channel represents a Twitter monitoring channel
type Channel struct {
	ID              string      `json:"id" gorm:"primaryKey"`
	OwnerID         int         `json:"ownerId"`
	IsVerified      bool        `json:"isVerified"`
	Name            string      `json:"name" gorm:"not null"`
	Description     string      `json:"description"`
	Avatar          string      `json:"avatar"`
	ChatLink        string      `json:"chatLink"`
	IsPublic        bool        `json:"isPublic"`
	IsHot           bool        `json:"isHot"`
	HotExpireAt     string      `json:"hotExpireAt"`
	CreatedAt       int64       `json:"createdAt"`
	UpdatedAt       int64       `json:"updatedAt"`
	Watchlist       []Watchlist `json:"watchlist" gorm:"type:jsonb"`
	Eventlist       []EventList `json:"eventlist" gorm:"type:jsonb"`
	FollowerCount   string      `json:"followerCount"`
	RecentFollowers []int       `json:"recentFollowers" gorm:"type:jsonb"`
}

// Watchlist represents a watched address in a channel
type Watchlist struct {
	TwitterName   string `json:"twitterName"`
	TwitterId     string `json:"twitterId"`
	Tweets        bool   `json:"tweets"`
	ProfileUpdate bool   `json:"profileUpdate"`
	Follows       bool   `json:"follows"`
	CA            string `json:"ca"` // CA address
	Risk          int    `json:"risk"`
}

// EventList represents an event filter in a channel
type EventList struct {
	FilterType   string        `json:"filterType"`
	OrConditions []OrCondition `json:"orConditions"`
}

// OrCondition represents a condition in an event filter
type OrCondition struct {
	AndConditions []AndCondition `json:"andConditions"`
}

// AndCondition represents a condition in an event filter
type AndCondition struct {
	Compare string   `json:"compare"`
	Field   string   `json:"field"`
	Values  []string `json:"values"`
}

// Owner represents the owner of a channel
type Owner struct {
	Email    string `json:"email"`
	UserID   int    `json:"userId"`
	UserName string `json:"userName"`
}

// ChannelDetail represents a channel with owner information
type ChannelDetail struct {
	Meta            ChannelMeta `json:"meta"`
	Owner           Owner       `json:"owner"`
	RecentFollowers []int       `json:"recentFollowers"`
}

// ChannelMeta represents the metadata of a channel
type ChannelMeta struct {
	Avatar        string      `json:"avatar"`
	ChainID       string      `json:"chainId"`
	ChatLink      string      `json:"chatLink"`
	CreatedAt     string      `json:"createdAt"`
	Description   string      `json:"description"`
	IsHot         bool        `json:"isHot"`
	HotExpireAt   string      `json:"hotExpireAt"`
	IsVerified    bool        `json:"isVerified"`
	Eventlist     []EventList `json:"eventlist"`
	FollowerCount string      `json:"followerCount"`
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	OwnerID       int         `json:"ownerId"`
	UpdatedAt     string      `json:"updatedAt"`
	Watchlist     []Watchlist `json:"watchlist"`
}

// CreateOrUpdateChannelRequest represents the request to create a channel
type CreateOrUpdateChannelRequest struct {
	UserID      int         `form:"userId"`
	Name        string      `json:"name" binding:"required"`
	Avatar      string      `json:"avatar" binding:"required"`
	Description string      `json:"description" binding:"required"`
	TwitterLink string      `json:"twitter" binding:"required"`
	IsPublic    bool        `json:"isPublic"`
	Watchlist   []Watchlist `json:"watchlist" binding:"required"`
	Eventlist   []EventList `json:"eventlist" binding:"required"`
}

// ChannelListResponse represents the response for channel list
type ChannelListResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Channels []struct {
			Meta            ChannelMeta `json:"meta"`
			Owner           Owner       `json:"owner"`
			RecentFollowers []string    `json:"recentFollowers"`
		} `json:"channels"`
		Total string `json:"total"`
	} `json:"data"`
}

// ChannelDetailResponse represents the response for channel detail
type ChannelDetailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Channel    ChannelDetail `json:"channel"`
		IsFollowed bool          `json:"isFollowed"`
	} `json:"data"`
}

// FollowRequest represents the request to follow a channel
type FollowRequest struct {
	ID     string `json:"id" binding:"required"`
	UserID int    `json:"userId" binding:"required"`
}

// UnfollowRequest represents the request to unfollow a channel
type UnfollowRequest struct {
	ID     string `json:"id" binding:"required"`
	UserID int    `json:"userId" binding:"required"`
}

// DeleteChannelRequest represents the request to delete a channel
type DeleteChannelRequest struct {
	ID     string `json:"id" binding:"required"`
	UserID int    `json:"userId" binding:"required"`
}

// ChannelListRequest represents the request to get channel list
type ChannelListRequest struct {
	UserID int    `form:"userId"`
	Type   string `form:"type"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
}

// APIResponse represents the common API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents the error response format
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Follow represents a user following a channel
type Follow struct {
	ID        string `json:"id" gorm:"primaryKey"`
	UserID    int    `json:"userId" gorm:"index"`
	ChannelID string `json:"channelId" gorm:"index"`
	CreatedAt int64  `json:"createdAt"`
}

// TwitterInfo represents a Twitter information record
type TwitterInfo struct {
	ID         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	TwitterId  string `json:"twitterId" gorm:"not null"`
	Content    string `json:"content" gorm:"type:longtext"`
	ChainId    string `json:"chainId"`
	Address    string `json:"address"`
	CreateTime int64  `json:"createTime"`
	Type       int    `json:"type" gorm:"not null"`
}

// ChannelContentRequest represents the request to get channel content
type ChannelContentRequest struct {
	ChannelID   string `form:"channelId" binding:"required"`
	Limit       int    `form:"limit"`
	Offset      int    `form:"offset"`
	Type        string `form:"type"`
	ContentType int    `form:"contentType" binding:"required"`
}
