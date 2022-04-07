package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LRichi/wcNATS/client"
	"go.uber.org/zap"
)

const subjectRequest = "test.subject.request"

type Request struct {
	Message string
}

type Response struct {
	Message string
}

type right struct {
	log *zap.SugaredLogger
}

func (r *right) receiveCall(ctx context.Context, req *Request) (*Response, error) {
	r.log.Infow("received call", "session", ctx.Value("session"), "request", req)

	if req.Message == "" {
		return &Response{}, fmt.Errorf("no messge")
	}

	return &Response{Message: "Yes, i'm fine"}, nil
}

func (r *right) receiveNotify(ctx context.Context, req *Request) error {
	r.log.Infow("received notify", "session", ctx.Value("session"), "request", req)

	// if use return error, search in NATS client debug
	if req.Message == "" {
		return fmt.Errorf("no message")
	}

	return nil
}

func executeRequest(ctx context.Context, log *zap.SugaredLogger, cli *client.Client) error {
	var (
		r = &right{log: log}
	)

	sub, err := cli.Subscribe(subjectRequest, r.receiveCall)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}
	defer func() {
		if err = cli.Unsubscribe(sub); err != nil {
			panic(err)
		}
	}()

	var (
		req  = Request{Message: "The one on the right, are you alive?"}
		resp Response
	)

	if err = cli.Request(ctx, subjectRequest, &req, &resp); err != nil {
		return fmt.Errorf("failed to request: %w", err)
	}

	return nil
}

func main() {
	lc := zap.NewDevelopmentConfig()
	lc.DisableStacktrace = true
	lc.DisableCaller = true

	log, err := lc.Build()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var (
		ctx = context.Background()
		cli = client.New(log.Sugar().Named("CLIENT"), "127.0.0.1:4222", "test", 100)
	)
	defer cli.Close()

	if err = executeRequest(ctx, log.Sugar(), cli); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
