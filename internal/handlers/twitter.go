package handlers

import (
	"TwitterMonitor/internal/models"
	"TwitterMonitor/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GetTwitterDetail(req *models.TwitterUserSearchRequest) ([]models.TwitterUserSearchResponse, error) {
	// 调用 Twitter API 获取用户信息
	url := fmt.Sprintf("http://43.160.199.161:5189/v1/users/search?regex=%s&screen_name=%spage=%d&pageSize=%d&token=%s",
		url.QueryEscape(req.Regex),
		url.QueryEscape(req.ScreenName),
		req.Page,
		req.PageSize,
		url.QueryEscape(req.Token))

	resp, err := http.Get(url)
	if err != nil {
		utils.LogError("Failed to retrieve data from external API", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.LogError("Failed to read response from external API", err)
		return nil, err
	}

	var response []models.TwitterUserSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		utils.LogError("Failed to parse response from external API", err)
		return nil, err
	}

	return response, nil
}
