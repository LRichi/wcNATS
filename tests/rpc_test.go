package tests

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/zap"

	"github.com/LRichi/wcNATS/client"
)

func requestTestWithoutError(ctx context.Context, log *zap.SugaredLogger, cli *client.Client) error {
	var (
		r = &right{log: log}
	)
	ctx = context.WithValue(ctx, "session", "111111")
	ctx = context.WithValue(ctx, "service", "goTest")
	ctx = context.WithValue(ctx, "method", "requestTest")

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

func requestTestWithError(ctx context.Context, log *zap.SugaredLogger, cli *client.Client) error {
	var (
		r = &right{log: log}
	)
	ctx = context.WithValue(ctx, "session", "2222222")
	ctx = context.WithValue(ctx, "service", "goTest")
	ctx = context.WithValue(ctx, "method", "requestTest")

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
		req  = Request{Message: ""}
		resp Response
	)

	if err = cli.Request(ctx, subjectRequest, &req, &resp); err != nil {
		return fmt.Errorf("failed to request: %w", err)
	}

	return nil
}

func TestNotifies_Request(t *testing.T) {
	lc := zap.NewDevelopmentConfig()
	lc.DisableStacktrace = true
	lc.DisableCaller = true

	log, err := lc.Build()
	if err != nil {
		t.Fatal(err)
	}

	var (
		ctx = context.Background()
		cli = client.New(log.Sugar().Named("NATS").Named("CLIENT"), "127.0.0.1:1222", "test", 100)
	)
	defer cli.Close()

	tests := []struct {
		name    string
		handle  func(ctx context.Context, log *zap.SugaredLogger, cli *client.Client) error
		wantErr bool
	}{
		{
			name:    "TEST_WITH_ERROR",
			handle:  requestTestWithoutError,
			wantErr: false,
		},
		{
			name:    "TEST_WITHOUT_ERROR",
			handle:  requestTestWithError,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = tt.handle(ctx, log.Sugar().Named(tt.name), cli); (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
