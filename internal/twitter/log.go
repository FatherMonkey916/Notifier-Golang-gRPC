package twitter

import (
	l "github.com/sirupsen/logrus"
	"notifier/grpc/notifier/twitter"
)

func logTweet(tweet *twitter.Tweet) {
	log.WithFields(l.Fields{
		"evcreated_at": tweet.CreatedAt,
		"id":           tweet.Id,
		"text":         tweet.Text,
		"source":       tweet.Source,
		"name":         tweet.User.Name,
	}).Info("Streaming successfully")
}
