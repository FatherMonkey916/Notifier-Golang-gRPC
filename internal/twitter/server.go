package twitter

import (
	l "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
	"notifier/grpc/notifier/twitter"
	"notifier/internal/notify"
	"os"
)

const (
	logfile = "/var/log/notifier/twitter.log"
)

var (
	log = l.New()
)

func init() {
	log.Formatter = new(l.JSONFormatter)

	log.Level = l.InfoLevel

	logfile, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Info(err)
	}

	log.Out = logfile
}

type TweetService struct {
	notify.Connection
}

func (ts *TweetService) Stream(srv twitter.TweetService_StreamServer) error {
	println("Started streaming data")
	ch := ts.Subscribe()
	for {
		select {
		case data := <-ch:
			_ = srv.Send(data.(*twitter.Tweet))
			logTweet(data.(*twitter.Tweet))
		}
	}
}

func (ts *TweetService) Receive(srv twitter.TweetService_ReceiveServer) error {

	packet, err := srv.Recv()
	if err == io.EOF {
		log.Info("Completed receiving data!")
	}

	if err != nil {
		log.Info("Error while receiving data:  %v", err)
	}

	ts.In(packet)
	return nil
}

func RegisterTwitterServer(server *grpc.Server) {
	tweetService := &TweetService{Connection: *notify.NewNotifier()}
	twitter.RegisterTweetServiceServer(server, tweetService)
	log.Info("Started Twitter server!")

	go tweetService.Notify()
}
