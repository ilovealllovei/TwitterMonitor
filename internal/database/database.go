package database

import (
	"TwitterMonitor/internal/models"
	"TwitterMonitor/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Database represents the MySQL connection
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new MySQL connection
func NewDatabase(uri string) (*Database, error) {
	db, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	// Check the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %v", err)
	}

	return &Database{db: db}, nil
}

// InsertOrUpdateChannel inserts a channel into the MySQL database or updates it if it already exists
func (db *Database) InsertOrUpdateChannel(channel *models.Channel) error {
	now := time.Now().UnixMilli()
	query := `INSERT INTO channels (id, ownerId, isVerified, name, description, avatar, chatLink, isPublic, isHot, hotExpireAt, createdAt, updatedAt, watchlist, eventlist, followerCount, recentFollowers) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	          ON DUPLICATE KEY UPDATE
	          isVerified = VALUES(isVerified),
	          name = VALUES(name),
	          description = VALUES(description),
	          avatar = VALUES(avatar),
	          chatLink = VALUES(chatLink),
	          isPublic = VALUES(isPublic),
	          isHot = VALUES(isHot),
	          hotExpireAt = VALUES(hotExpireAt),
	          updatedAt = VALUES(updatedAt),
	          watchlist = VALUES(watchlist),
	          eventlist = VALUES(eventlist),
	          followerCount = VALUES(followerCount),
	          recentFollowers = VALUES(recentFollowers)`

	// 对 Watchlist 进行 JSON 编码
	watchlistJSON, err := json.Marshal(channel.Watchlist)
	if err != nil {
		return fmt.Errorf("failed to marshal watchlist: %v", err)
	}
	watchlistStr := string(watchlistJSON)

	// 对 recentFollowers 进行 JSON 编码
	recentFollowersJSON, err := json.Marshal(channel.RecentFollowers)
	if err != nil {
		return fmt.Errorf("failed to marshal RecentFollowers: %v", err)
	}
	recentFollowersStr := string(recentFollowersJSON)

	// 对 Eventlist 进行 JSON 编码
	eventlistJSON, err := json.Marshal(channel.Eventlist)
	if err != nil {
		return fmt.Errorf("failed to marshal eventlist: %v", err)
	}
	eventlistStr := string(eventlistJSON)

	_, err = db.db.Exec(query,
		channel.ID,
		channel.OwnerID,
		channel.IsVerified,
		channel.Name,
		channel.Description,
		channel.Avatar,
		channel.ChatLink,
		channel.IsPublic,
		channel.IsHot,
		channel.HotExpireAt,
		now,
		now,
		watchlistStr,
		eventlistStr,
		channel.FollowerCount,
		recentFollowersStr,
	)
	if err != nil {
		return fmt.Errorf("failed to insert or update channel: %v", err)
	}
	return nil
}

// GetChannelsByOwnerID retrieves channels by OwnerID
func (db *Database) GetChannelsByOwnerID(ownerID int) ([]*models.Channel, error) {
	query := `SELECT * FROM channels WHERE ownerId = ?`
	rows, err := db.db.Query(query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query channels: %v", err)
	}
	defer rows.Close()

	var channels []*models.Channel
	for rows.Next() {
		channel := &models.Channel{}
		var watchlistStr, eventlistStr, recentFollowersStr string
		err := rows.Scan(
			&channel.ID,
			&channel.OwnerID,
			&channel.IsVerified,
			&channel.Name,
			&channel.Description,
			&channel.Avatar,
			&channel.ChatLink,
			&channel.IsPublic,
			&channel.IsHot,
			&channel.HotExpireAt,
			&channel.CreatedAt,
			&channel.UpdatedAt,
			&watchlistStr,
			&eventlistStr,
			&channel.FollowerCount,
			&recentFollowersStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan channel: %v", err)
		}

		// Unmarshal Watchlist
		if err := json.Unmarshal([]byte(watchlistStr), &channel.Watchlist); err != nil {
			return nil, fmt.Errorf("failed to unmarshal watchlist: %v", err)
		}

		// Unmarshal Eventlist
		if err := json.Unmarshal([]byte(eventlistStr), &channel.Eventlist); err != nil {
			return nil, fmt.Errorf("failed to unmarshal eventlist: %v", err)
		}

		// Unmarshal RecentFollowers
		if err := json.Unmarshal([]byte(recentFollowersStr), &channel.RecentFollowers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal recentFollowers: %v", err)
		}

		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return channels, nil
}

// GetChannelsByID retrieves channels by OwnerID
func (db *Database) GetChannelsByID(id string) ([]*models.Channel, error) {
	query := `SELECT * FROM channels WHERE id = ?`
	rows, err := db.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query channels: %v", err)
	}
	defer rows.Close()

	var channels []*models.Channel
	for rows.Next() {
		channel := &models.Channel{}
		var watchlistStr, eventlistStr, recentFollowersStr string
		err := rows.Scan(
			&channel.ID,
			&channel.OwnerID,
			&channel.IsVerified,
			&channel.Name,
			&channel.Description,
			&channel.Avatar,
			&channel.ChatLink,
			&channel.IsPublic,
			&channel.IsHot,
			&channel.HotExpireAt,
			&channel.CreatedAt,
			&channel.UpdatedAt,
			&watchlistStr,
			&eventlistStr,
			&channel.FollowerCount,
			&recentFollowersStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan channel: %v", err)
		}

		// Unmarshal Watchlist
		if err := json.Unmarshal([]byte(watchlistStr), &channel.Watchlist); err != nil {
			return nil, fmt.Errorf("failed to unmarshal watchlist: %v", err)
		}

		// Unmarshal Eventlist
		if err := json.Unmarshal([]byte(eventlistStr), &channel.Eventlist); err != nil {
			return nil, fmt.Errorf("failed to unmarshal eventlist: %v", err)
		}

		// Unmarshal RecentFollowers
		if err := json.Unmarshal([]byte(recentFollowersStr), &channel.RecentFollowers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal recentFollowers: %v", err)
		}

		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return channels, nil
}

// DeleteChannel deletes a channel by its ID
func (db *Database) DeleteChannel(channelID string) error {
	query := `DELETE FROM channels WHERE id = ?`
	result, err := db.db.Exec(query, channelID)
	if err != nil {
		return fmt.Errorf("failed to delete channel: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("channel not found")
	}

	return nil
}

// UpdateFollowerCount updates the follower count and recent followers of a channel
func (db *Database) UpdateFollowerCount(channelID string, userID int, increment bool) error {
	// Start transaction
	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Update follower count
	operator := "-"
	if increment {
		operator = "+"
	}
	query := fmt.Sprintf("UPDATE channels SET followerCount = CAST(followerCount AS SIGNED) %s 1 WHERE id = ?", operator)
	_, err = tx.Exec(query, channelID)
	if err != nil {
		return fmt.Errorf("failed to update follower count: %v", err)
	}

	if increment {
		// Get current recentFollowersStr
		var currentFollowersStr string
		err = tx.QueryRow("SELECT recentFollowersStr FROM channels WHERE id = ?", channelID).Scan(&currentFollowersStr)
		if err != nil {
			return fmt.Errorf("failed to get current followers: %v", err)
		}

		// Parse current followers
		var currentFollowers []int
		if currentFollowersStr != "" {
			if err := json.Unmarshal([]byte(currentFollowersStr), &currentFollowers); err != nil {
				return fmt.Errorf("failed to parse current followers: %v", err)
			}
		}

		// Add new follower and ensure uniqueness
		exists := false
		for _, id := range currentFollowers {
			if id == userID {
				exists = true
				break
			}
		}
		if !exists {
			// Add new follower at the beginning
			currentFollowers = append([]int{userID}, currentFollowers...)
			// Keep only the most recent 50 followers
			if len(currentFollowers) > 50 {
				currentFollowers = currentFollowers[:50]
			}
		}

		// Convert back to JSON
		newFollowersStr, err := json.Marshal(currentFollowers)
		if err != nil {
			return fmt.Errorf("failed to marshal followers: %v", err)
		}

		// Update recentFollowersStr
		_, err = tx.Exec("UPDATE channels SET recentFollowersStr = ? WHERE id = ?", string(newFollowersStr), channelID)
		if err != nil {
			return fmt.Errorf("failed to update recent followers: %v", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// FollowChannel creates a follow relationship between a user and a channel
func (db *Database) FollowChannel(follow *models.Follow) error {
	now := time.Now().UnixMilli()

	// Start a transaction
	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Insert follow relationship
	query := `INSERT INTO follows (id, userId, channelId, createdAt) VALUES (?, ?, ?, ?)`
	_, err = tx.Exec(query,
		follow.ID,
		follow.UserID,
		follow.ChannelID,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to follow channel: %v", err)
	}

	// Update follower count
	if err := db.UpdateFollowerCount(follow.ChannelID, follow.UserID, true); err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// UnfollowChannel removes a follow relationship between a user and a channel
func (db *Database) UnfollowChannel(userID int, channelID string) error {
	// Start a transaction
	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Delete follow relationship
	query := `DELETE FROM follows WHERE userId = ? AND channelId = ?`
	result, err := tx.Exec(query, userID, channelID)
	if err != nil {
		return fmt.Errorf("failed to unfollow channel: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("follow relationship not found")
	}

	// Update follower count
	if err := db.UpdateFollowerCount(channelID, userID, false); err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// IsFollowing checks if a user is following a channel
func (db *Database) IsFollowing(userID int, channelID string) (bool, error) {
	query := `SELECT COUNT(*) FROM follows WHERE userId = ? AND channelId = ?`
	var count int
	err := db.db.QueryRow(query, userID, channelID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check follow status: %v", err)
	}
	return count > 0, nil
}

// GetTwitterInfoByTypeAndCA gets Twitter info by type and CA address
func (db *Database) GetTwitterInfoByTypeAndCA(type_ int, caAddresses []string, limit, offset int) ([]models.TwitterInfo, error) {
	var twitterInfos []models.TwitterInfo
	var query string
	var args []interface{}

	query = "SELECT id, twitterId, content, chainId, address, createTime, type FROM twitter_info WHERE type = ?"
	args = append(args, type_)

	if len(caAddresses) > 0 {
		placeholders := make([]string, len(caAddresses))
		for i := range caAddresses {
			placeholders[i] = "?"
			args = append(args, caAddresses[i])
		}
		query += " AND address IN (" + strings.Join(placeholders, ",") + ")"
	}

	query += " ORDER BY createTime DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var info models.TwitterInfo
		err := rows.Scan(&info.ID, &info.TwitterId, &info.Content, &info.ChainId, &info.Address, &info.CreateTime, &info.Type)
		if err != nil {
			return nil, err
		}
		twitterInfos = append(twitterInfos, info)
	}

	return twitterInfos, nil
}

// GetTwitterInfoByType gets Twitter info by type only
func (db *Database) GetTwitterInfoByType(type_ int, limit, offset int) ([]models.TwitterInfo, error) {
	var twitterInfos []models.TwitterInfo

	query := "SELECT id, twitterId, content, chainId, address, createTime, type FROM twitter_info WHERE type = ? ORDER BY createTime DESC LIMIT ? OFFSET ?"
	rows, err := db.db.Query(query, type_, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var info models.TwitterInfo
		err := rows.Scan(&info.ID, &info.TwitterId, &info.Content, &info.ChainId, &info.Address, &info.CreateTime, &info.Type)
		if err != nil {
			return nil, err
		}
		twitterInfos = append(twitterInfos, info)
	}

	return twitterInfos, nil
}

// GetTwitterInfoByWatchlist gets Twitter info based on watchlist conditions
func (db *Database) GetTwitterInfoByWatchlist(conditions []string, contentType int, limit, offset int) ([]*models.TwitterInfo, error) {
	query := `
		SELECT id, twitterId, content, COALESCE(chainId, ''), COALESCE(address, ''), createTime, type
		FROM twitter_info
		WHERE type = ?
	`

	if len(conditions) > 0 {
		query += " AND (" + strings.Join(conditions, " OR ") + ")"
	}

	query += " ORDER BY createTime DESC LIMIT ? OFFSET ?"

	rows, err := db.db.Query(query, contentType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query Twitter info: %v", err)
	}
	defer rows.Close()

	var twitterInfos []*models.TwitterInfo
	for rows.Next() {
		var info models.TwitterInfo
		if err := rows.Scan(&info.ID, &info.TwitterId, &info.Content, &info.ChainId, &info.Address, &info.CreateTime, &info.Type); err != nil {
			return nil, fmt.Errorf("failed to scan Twitter info: %v", err)
		}
		twitterInfos = append(twitterInfos, &info)
	}

	return twitterInfos, nil
}

// GetTwitterInfoByProfileAndFollow gets Twitter info based on profile and follow conditions
func (db *Database) GetTwitterInfoByProfileAndFollow(twitterIds []string, contentType int, limit, offset int) ([]*models.TwitterInfo, error) {
	if len(twitterIds) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(twitterIds))
	args := make([]interface{}, len(twitterIds)+3)
	for i := range twitterIds {
		placeholders[i] = "?"
		args[i] = twitterIds[i]
	}
	args[len(twitterIds)] = contentType
	args[len(twitterIds)+1] = limit
	args[len(twitterIds)+2] = offset

	query := fmt.Sprintf(`
		SELECT id, twitter_id, content, COALESCE(chain_id, ''), COALESCE(address, ''), create_time, type
		FROM twitter_info
		WHERE twitter_id IN (%s) AND type = ?
		ORDER BY create_time DESC
		LIMIT ? OFFSET ?
	`, strings.Join(placeholders, ","))

	rows, err := db.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query Twitter info: %v", err)
	}
	defer rows.Close()

	var twitterInfos []*models.TwitterInfo
	for rows.Next() {
		var info models.TwitterInfo
		if err := rows.Scan(&info.ID, &info.TwitterId, &info.Content, &info.ChainId, &info.Address, &info.CreateTime, &info.Type); err != nil {
			return nil, fmt.Errorf("failed to scan Twitter info: %v", err)
		}
		twitterInfos = append(twitterInfos, &info)
	}

	return twitterInfos, nil
}

// GetFollowedChannels gets all channels followed by a user
func (db *Database) GetFollowedChannels(userID int) ([]*models.Follow, error) {
	var follows []*models.Follow
	rows, err := db.db.Query("SELECT * FROM follows WHERE userId = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var follow models.Follow
		err := rows.Scan(&follow.ID, &follow.UserID, &follow.ChannelID, &follow.CreatedAt)
		if err != nil {
			return nil, err
		}
		follows = append(follows, &follow)
	}

	return follows, nil
}

// GetChannelByID gets a channel by its ID
func (db *Database) GetChannelByID(channelID string) (*models.Channel, error) {
	var channel models.Channel
	row := db.db.QueryRow("SELECT * FROM channels WHERE id = ?", channelID)
	var watchlistStr, eventlistStr, recentFollowersStr string
	err := row.Scan(
		&channel.ID,
		&channel.OwnerID,
		&channel.IsVerified,
		&channel.Name,
		&channel.Description,
		&channel.Avatar,
		&channel.ChatLink,
		&channel.IsPublic,
		&channel.IsHot,
		&channel.HotExpireAt,
		&channel.CreatedAt,
		&channel.UpdatedAt,
		&watchlistStr,
		&eventlistStr,
		&channel.FollowerCount,
		&recentFollowersStr,
	)

	// Unmarshal Watchlist
	if err := json.Unmarshal([]byte(watchlistStr), &channel.Watchlist); err != nil {
		return nil, fmt.Errorf("failed to unmarshal watchlist: %v", err)
	}

	// Unmarshal Eventlist
	if err := json.Unmarshal([]byte(eventlistStr), &channel.Eventlist); err != nil {
		return nil, fmt.Errorf("failed to unmarshal eventlist: %v", err)
	}

	// Unmarshal RecentFollowers
	if err := json.Unmarshal([]byte(recentFollowersStr), &channel.RecentFollowers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recentFollowers: %v", err)
	}

	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// GetAllChannels gets all channels with pagination
func (db *Database) GetAllChannels(limit, offset int) ([]*models.Channel, error) {
	var channels []*models.Channel
	query := "SELECT * FROM channels"
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var channel models.Channel
		var watchlistStr, eventlistStr, recentFollowersStr string
		err := rows.Scan(
			&channel.ID,
			&channel.OwnerID,
			&channel.IsVerified,
			&channel.Name,
			&channel.Description,
			&channel.Avatar,
			&channel.ChatLink,
			&channel.IsPublic,
			&channel.IsHot,
			&channel.HotExpireAt,
			&channel.CreatedAt,
			&channel.UpdatedAt,
			&watchlistStr,
			&eventlistStr,
			&channel.FollowerCount,
			&recentFollowersStr,
		)
		if err != nil {
			utils.LogError("Error scanning channel", "error", err)
			return nil, err
		}

		// Unmarshal Watchlist
		if err := json.Unmarshal([]byte(watchlistStr), &channel.Watchlist); err != nil {
			return nil, fmt.Errorf("failed to unmarshal watchlist: %v", err)
		}

		// Unmarshal Eventlist
		if err := json.Unmarshal([]byte(eventlistStr), &channel.Eventlist); err != nil {
			return nil, fmt.Errorf("failed to unmarshal eventlist: %v", err)
		}

		// Unmarshal RecentFollowers
		if err := json.Unmarshal([]byte(recentFollowersStr), &channel.RecentFollowers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal recentFollowers: %v", err)
		}

		channels = append(channels, &channel)
	}

	return channels, nil
}

// GetChannelByIDs gets multiple channels by their IDs in a single query
func (db *Database) GetChannelByIDs(channelIDs []string) ([]*models.Channel, error) {
	if len(channelIDs) == 0 {
		return nil, nil
	}

	// Build the query with placeholders
	placeholders := make([]string, len(channelIDs))
	args := make([]interface{}, len(channelIDs))
	for i, id := range channelIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := "SELECT * FROM channels WHERE id IN (" + strings.Join(placeholders, ",") + ")"
	rows, err := db.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query channels: %v", err)
	}
	defer rows.Close()

	var channels []*models.Channel
	for rows.Next() {
		channel := &models.Channel{}
		var watchlistStr, eventlistStr, recentFollowersStr string
		err := rows.Scan(
			&channel.ID,
			&channel.OwnerID,
			&channel.IsVerified,
			&channel.Name,
			&channel.Description,
			&channel.Avatar,
			&channel.ChatLink,
			&channel.IsPublic,
			&channel.IsHot,
			&channel.HotExpireAt,
			&channel.CreatedAt,
			&channel.UpdatedAt,
			&watchlistStr,
			&eventlistStr,
			&channel.FollowerCount,
			&recentFollowersStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan channel: %v", err)
		}

		// Unmarshal Watchlist
		if err := json.Unmarshal([]byte(watchlistStr), &channel.Watchlist); err != nil {
			return nil, fmt.Errorf("failed to unmarshal watchlist: %v", err)
		}

		// Unmarshal Eventlist
		if err := json.Unmarshal([]byte(eventlistStr), &channel.Eventlist); err != nil {
			return nil, fmt.Errorf("failed to unmarshal eventlist: %v", err)
		}

		// Unmarshal RecentFollowers
		if err := json.Unmarshal([]byte(recentFollowersStr), &channel.RecentFollowers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal recentFollowers: %v", err)
		}

		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return channels, nil
}

// GetRecentFollowers gets the most recent followers for a channel
func (db *Database) GetRecentFollowers(channelID string, limit int) ([]int, error) {
	if limit <= 0 {
		limit = 100 // Default to 100 if not specified
	}

	query := `
		SELECT userId 
		FROM follows 
		WHERE channelId = ? 
		ORDER BY createdAt DESC 
		LIMIT ?
	`

	rows, err := db.db.Query(query, channelID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent followers: %v", err)
	}
	defer rows.Close()

	var followers []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("failed to scan follower: %v", err)
		}
		followers = append(followers, userID)
	}

	return followers, nil
}
