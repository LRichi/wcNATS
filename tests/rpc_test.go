package tests

import (
	"context"
	"fmt"
	"sync"
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

func notifyTestWithoutError(ctx context.Context, log *zap.SugaredLogger, cli *client.Client) error {
	var (
		r = &right{
			wg:  &sync.WaitGroup{},
			log: log,
		}
	)
	ctx = context.WithValue(ctx, "session", "222222")
	ctx = context.WithValue(ctx, "service", "goTest")
	ctx = context.WithValue(ctx, "method", "notifyTest")

	sub, err := cli.Subscribe(subjectRequest, r.receiveNotify)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}
	defer func() {
		if err = cli.Unsubscribe(sub); err != nil {
			panic(err)
		}
	}()

	var (
		req = Request{Message: "send data"}
	)

	r.wg.Add(1)
	if err = cli.Publish(ctx, subjectRequest, &req); err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}
	r.wg.Wait()

	return nil
}

func notifyTestWithError(ctx context.Context, log *zap.SugaredLogger, cli *client.Client) error {
	var (
		r = &right{
			wg:  &sync.WaitGroup{},
			log: log,
		}
	)
	ctx = context.WithValue(ctx, "session", "222222")
	ctx = context.WithValue(ctx, "service", "goTest")
	ctx = context.WithValue(ctx, "method", "notifyTest")

	sub, err := cli.Subscribe(subjectRequest, r.receiveNotify)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}
	defer func() {
		if err = cli.Unsubscribe(sub); err != nil {
			panic(err)
		}
	}()

	var (
		req = Request{Message: ""}
	)

	r.wg.Add(1)
	if err = cli.Publish(ctx, subjectRequest, &req); err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}
	r.wg.Wait()

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
			name:    "TEST_WITHOUT_ERROR",
			handle:  requestTestWithoutError,
			wantErr: false,
		},
		{
			name:    "TEST_WITh_ERROR",
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

func TestNotifies_Notify(t *testing.T) {
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
			handle:  notifyTestWithoutError,
			wantErr: false,
		},
		{
			name:    "TEST_WITHOUT_ERROR",
			handle:  notifyTestWithError,
			wantErr: false,
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
