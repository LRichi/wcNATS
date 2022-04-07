package tests

import (
	"context"
	"fmt"
	"strings"

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
	var caller = []string{"service", ctx.Value("service").(string), "session", ctx.Value("method").(string)}
	r.log.Infow("received call",
		"session", ctx.Value("session"), "caller", strings.Join(caller, "."), "request", req)

	if req.Message == "" {
		return &Response{}, fmt.Errorf("no messge")
	}

	return &Response{Message: "Yes, i'm fine"}, nil
}

func (r *right) receiveNotify(ctx context.Context, req *Request) error {
	var caller = []string{"service", ctx.Value("service").(string), "session", ctx.Value("method").(string)}
	r.log.Infow("received notify",
		"session", ctx.Value("session"), "caller", strings.Join(caller, "."), "request", req)

	// if use return error, search in NATS client debug
	//return fmt.Errorf("test error")

	return nil
}
