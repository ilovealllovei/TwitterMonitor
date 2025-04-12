package models

type TwitterUserSearchRequest struct {
	Regex      string `json:"regex" binding:"required"`
	ScreenName string `json:"screen_name" binding:"required"`
	PageSize   int    `json:"pageSize" binding:"required"`
	Page       int    `json:"page" binding:"required"`
	Token      string `json:"token" binding:"required"`
}

type TwitterUserSearchResponse struct {
	ID                  string   `json:"_id"`
	ConversationIDStr   string   `json:"conversation_id_str"`
	ContentImg          string   `json:"content_Img"`
	ContentVideo        string   `json:"content_video"`
	IDCar               string   `json:"idCar"`
	LastUpdateTime      string   `json:"lastupdatetime"`
	MyFavoriteCount     int      `json:"myFavoriteCount"`
	MyFullText          string   `json:"myFullText"`
	MyQuoteCount        int      `json:"myQuoteCount"`
	MyReplyCount        int      `json:"myReplyCount"`
	MyRetweetCount      int      `json:"myRetweetCount"`
	MyTwTime            string   `json:"myTwTime"`
	OtherAccount        string   `json:"other_Account"`
	OtherFavoriteCount  string   `json:"other_FavoriteCount"`
	OtherNameId         string   `json:"other_NameId"`
	OtherQuoteCount     string   `json:"other_QuoteCount"`
	OtherReplyCount     string   `json:"other_ReplyCount"`
	OtherRetweetCount   string   `json:"other_RetweetCount"`
	OtherText           string   `json:"other_Text"`
	OtherHeadesPhoto    string   `json:"other_heades_photo"`
	OtherJoinTime       string   `json:"other_joinTime"`
	OtherName           string   `json:"other_name"`
	OtherProfileImg     string   `json:"other_profile_Img"`
	OtherTextID         string   `json:"other_text_id"`
	QuotedAccount       string   `json:"quoted_Account"`
	QuotedFavoriteCount int      `json:"quoted_FavoriteCount"`
	QuotedQuoteCount    int      `json:"quoted_QuoteCount"`
	QuotedReplyCount    int      `json:"quoted_ReplyCount"`
	QuotedRetweetCount  int      `json:"quoted_RetweetCount"`
	QuotedText          string   `json:"quoted_Text"`
	QuotedHeadesPhoto   []string `json:"quoted_heades_photo"`
	QuotedJoinTime      string   `json:"quoted_joinTime"`
	QuotedName          string   `json:"quoted_name"`
	QuotedProfileImg    []string `json:"quoted_profile_Img"`
	QuotedTextID        string   `json:"quoted_text_id"`
	ScreenName          string   `json:"screen_name"`
	User                struct {
		UserID                string   `json:"user_id"`
		Description           string   `json:"description"`
		JoinTime              string   `json:"join_time"`
		FriendsCount          int      `json:"friends_count"`
		FollowersCount        int      `json:"followers_count"`
		Location              string   `json:"location"`
		Name                  string   `json:"name"`
		ScreenName            string   `json:"screen_name"`
		HeaderPhoto           []string `json:"header_photo"`
		StatusesCount         int      `json:"statuses_count"`
		FavouritesCount       int      `json:"favourites_count"`
		ProfileBannerImageURL []string `json:"profile_banner_image_url"`
	} `json:"user"`
	HasTokens map[string][]string `json:"has_tokens"`
}
