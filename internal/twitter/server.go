package twitter

import (
	"google.golang.org/grpc"
	"io"
	"log"
	"notifier/grpc/notifier/twitter"
	"notifier/internal/notify"
)

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
		}
	}
}

func (ts *TweetService) Receive(srv twitter.TweetService_ReceiveServer) error {

	packet, err := srv.Recv()
	if err == io.EOF {
		println("Completed receiving data!")
	}

	if err != nil {
		log.Fatalf("Error while receiving data:  %v", err)
	}

	ts.In(packet)
	return nil
}

func RegisterTwitterServer(server *grpc.Server) {
	tweetService := &TweetService{Connection: *notify.NewNotifier()}
	twitter.RegisterTweetServiceServer(server, tweetService)
	log.Println("Started Twitter server!")

	go tweetService.Notify()
}
