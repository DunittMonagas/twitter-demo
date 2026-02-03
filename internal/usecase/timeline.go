package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"twitter-demo/internal/config"
	"twitter-demo/internal/domain"
	"twitter-demo/internal/infrastructure/repository"
	"twitter-demo/pkg"
)

const (
	// MaxCachedTweets defines the maximum number of tweet IDs to keep in cache per user
	MaxCachedTweets = 1000
	// CacheExpiration defines the expiration time for the cache, set to 14 days
	CacheExpiration = 14 * 24 * time.Hour
	// CacheKey defines the key for the cache
	CacheKey = "timeline:user:%d"
)

type TimelineUsecase interface {
	GetTimeline(ctx context.Context, userID int64, limit, offset int) ([]domain.Tweet, error)
	FanOutTweet(ctx context.Context, authorID int64, tweetID int64) error
}

type Timeline struct {
	tweetRepository    repository.TweetRepository
	followerRepository repository.FollowerRepository
	cache              pkg.Cache
}

func NewTimeline(tweetRepository repository.TweetRepository, followerRepository repository.FollowerRepository, cache pkg.Cache) Timeline {
	return Timeline{
		tweetRepository:    tweetRepository,
		followerRepository: followerRepository,
		cache:              cache,
	}
}

func (t Timeline) GetTimeline(ctx context.Context, userID int64, limit, offset int) ([]domain.Tweet, error) {

	// Set default and max values for pagination
	if limit <= 0 {
		limit = config.DefaultLimit
	}

	if limit > config.MaxLimit {
		limit = config.MaxLimit
	}

	if offset < 0 {
		offset = 0
	}

	// STEP 1: Try to get tweet IDs from cache
	tweetIDs, err := t.retrieveCacheTweetIDs(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// STEP 2: Handle cache hit - we got IDs from cache
	if len(tweetIDs) == limit {
		// We have enough IDs in cache, fetch tweets from DB using these IDs
		tweets, err := t.tweetRepository.SelectTweetsByIDs(ctx, tweetIDs)
		if err == nil {
			return tweets, nil
		}
		// If DB fetch failed, fall through to STEP 3
	}

	// STEP 3: Cache miss or partial miss - fall back to database
	// This happens when:
	// - Redis doesn't have the key (err == redis.Nil)
	// - Redis returned fewer IDs than requested (partial miss)
	// - DB fetch by IDs failed
	tweets, err := t.tweetRepository.SelectTimelineTweets(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// STEP 4: Populate cache with results (only if offset is 0)
	// We only cache the "fresh" timeline (page 1) to keep cache simple
	// Deeper pages will always hit the database
	if offset == 0 && len(tweets) > 0 {
		go t.cacheTimelineTweets(context.Background(), userID, tweets)
	}

	return tweets, nil
}

// cacheTimelineTweets stores tweet IDs in Redis for future cache hits.
// This runs asynchronously to avoid blocking the request.
func (t Timeline) cacheTimelineTweets(ctx context.Context, userID int64, tweets []domain.Tweet) {
	if len(tweets) == 0 {
		return
	}

	cacheKey := t.getCacheKey(userID)

	// Extract tweet IDs in order
	tweetIDs := make([]interface{}, len(tweets))
	for i, tweet := range tweets {
		tweetIDs[i] = fmt.Sprintf("%d", tweet.ID)
	}

	// DISCLAIMER: These operations should be atomic
	// Delete old cache and push new IDs
	_ = t.cache.Delete(ctx, cacheKey)

	// Push all IDs using RPush to maintain the order from DB (newest first)
	if err := t.cache.RPush(ctx, cacheKey, tweetIDs...); err != nil {
		return
	}

	// Trim to keep only the most recent tweets
	_ = t.cache.LTrim(ctx, cacheKey, 0, MaxCachedTweets-1)

	// Set expiration time for the cache
	_ = t.cache.Expire(ctx, cacheKey, CacheExpiration)

}

// getCacheKey constructs the cache key for a user's timeline.
func (t Timeline) getCacheKey(userID int64) string {
	return fmt.Sprintf(CacheKey, userID)
}

// retrieveCacheTweetIDs retrieves tweet IDs from cache.
func (t Timeline) retrieveCacheTweetIDs(ctx context.Context, userID int64, limit, offset int) ([]int64, error) {

	cacheKey := t.getCacheKey(userID)
	startIdx := int64(offset)
	stopIdx := int64(offset + limit - 1)

	tweetIDs, err := t.cache.LRange(ctx, cacheKey, startIdx, stopIdx)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(tweetIDs))
	for _, idStr := range tweetIDs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// FanOutTweet distributes a new tweet to all followers' timeline caches.
// This implements the Fan-Out pattern for real-time timeline updates.
func (t Timeline) FanOutTweet(ctx context.Context, authorID int64, tweetID int64) error {
	// Step 1: Get all followers of the tweet author
	followerIDs, err := t.followerRepository.SelectFollowerIDsByFollowedID(ctx, authorID)
	if err != nil {
		return fmt.Errorf("failed to get followers: %w", err)
	}

	// If author has no followers, nothing to do
	if len(followerIDs) == 0 {
		return nil
	}

	// Step 2: Add tweet to each follower's timeline cache
	tweetIDStr := fmt.Sprintf("%d", tweetID)
	for _, followerID := range followerIDs {
		cacheKey := t.getCacheKey(followerID)

		// Add tweet ID to the beginning of the list (newest first) using LPUSH
		if err := t.cache.LPush(ctx, cacheKey, tweetIDStr); err != nil {
			// Log error but continue with other followers
			fmt.Printf("Failed to add tweet %d to user %d timeline: %v\n", tweetID, followerID, err)
			continue
		}

		// Trim to keep only the latest MaxCachedTweets
		if err := t.cache.LTrim(ctx, cacheKey, 0, MaxCachedTweets-1); err != nil {
			// Log error but continue
			fmt.Printf("Failed to trim cache for user %d: %v\n", followerID, err)
		}

		// Update expiration
		if err := t.cache.Expire(ctx, cacheKey, CacheExpiration); err != nil {
			// Log error but continue
			fmt.Printf("Failed to set expiration for user %d: %v\n", followerID, err)
		}
	}

	return nil
}
