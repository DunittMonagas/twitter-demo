package config

// Kafka topic names
const (
	TopicTweets  = "tweets"
	TopicFollows = "follows"
)

// Kafka message key formats
const (
	KeyFormatTweet = "tweet-%d" // tweet-{tweetID}
)
