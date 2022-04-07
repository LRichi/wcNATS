package client

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.uber.org/zap"

	"github.com/nats-io/nats.go"
)

type Client struct {
	*conn
	log *zap.SugaredLogger
}

// Request - a remote procedure call is created
//
// use handle: func(context.Context,*struct)(*struct,error)
func (c *Client) Request(ctx context.Context, subject string, request, response interface{}) (err error) {
	start := time.Now()
	defer func() {
		c.log.Debugw("Request",
			"subject", subject, "elapsed", time.Since(start).Seconds(),
			"request", request, "response", response, "error", err,
		)
	}()

	// validate request and response
	if err = validateModel(request); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	if err = validateModel(response); err != nil {
		return fmt.Errorf("invalid response: %w", err)
	}

	// create a DTO in memory
	var (
		reqDTO  = newRequestDTO(reflect.TypeOf(request))
		respDTO = newResponseDTO(reflect.TypeOf(response))
	)
	reqDTO.FieldByName("Session").Set(reflect.ValueOf(getSession(ctx)))
	reqDTO.FieldByName("Request").Set(reflect.ValueOf(request))

	// call
	if err = c.request(ctx, subject, reqDTO.Interface(), respDTO.Addr().Interface()); err != nil {
		return convertErr(err)
	}

	// create response
	valResp := reflect.ValueOf(response)
	if !respDTO.FieldByName("Response").IsZero() {
		valResp.Elem().Set(respDTO.FieldByName("Response").Elem())
	}

	// check error
	if respDTO.FieldByName("Error").Interface().(ErrorDTO).Type != nil {
		return respDTO.FieldByName("Error").Interface().(ErrorDTO)
	}

	return nil
}

// Publish - for notify subscribers
//
// use handle: func(context.Context,*struct)error
func (c *Client) Publish(ctx context.Context, subject fmt.Stringer, value interface{}) (err error) {
	var start = time.Now()
	defer func() {
		c.log.Debugw("Publish", "subject", subject, "elapsed", time.Since(start).Seconds(),
			"value", value, "error", err,
		)
	}()

	if err = validateModel(value); err != nil {
		return fmt.Errorf("invalid value: %w", err)
	}

	// create a DTO in memory
	var reqDTO = newRequestDTO(reflect.TypeOf(value))
	reqDTO.FieldByName("Session").Set(reflect.ValueOf(getSession(ctx)))
	reqDTO.FieldByName("Request").Set(reflect.ValueOf(value))

	return c.conn.publish(subject.String(), reqDTO.Interface())
}

// Subscribe - subscribe handle for remote call or notify
//
// use handle:
//	for rpc:    func(context.Context,*struct)(*struct,error)
//	for notify: func(context.Context,*struct)(error)
func (c *Client) Subscribe(subject string, handle interface{}) (Subscription, error) {

	var (
		start = time.Now()
		err   error
	)
	defer func() {
		c.log.Debugw("Subscribe", "subject", subject, "elapsed", time.Since(start).Seconds(),
			"handle", reflect.TypeOf(handle).String(), "error", err,
		)
	}()

	isRequest, err := isRequestHandle(handle)
	if err != nil {
		return nil, err
	}

	var (
		sub = &subscription{
			log:          c.log,
			Subscription: nil,
			process:      reflect.ValueOf(handle),
			response:     reflect.ValueOf(c.publish),
		}
		reqDTO           = newRequestDTO(reflect.TypeOf(handle).In(1))
		typeOfDtoReqIn   = []reflect.Type{reflect.TypeOf(""), reflect.TypeOf(""), reqDTO.Type()}
		typeDtoReqHandle = reflect.FuncOf(typeOfDtoReqIn, nil, false)
	)

	var valueDtoReqHandle reflect.Value
	if isRequest {
		if err = validateHandleOfCall(handle); err != nil {
			return nil, fmt.Errorf("invalid handle: %w", err)
		}
		valueDtoReqHandle = reflect.MakeFunc(typeDtoReqHandle, sub.call)
	} else {
		if err = validateHandleOfNotify(handle); err != nil {
			return nil, fmt.Errorf("invalid handle: %w", err)
		}
		valueDtoReqHandle = reflect.MakeFunc(typeDtoReqHandle, sub.notify)
	}

	sub.Subscription, err = c.subscribe(subject, valueDtoReqHandle.Interface())
	if err != nil {
		return nil, convertErr(err)
	}

	return sub, nil
}

// Unsubscribe - delete subscription
func (c *Client) Unsubscribe(sub Subscription) (err error) {
	var (
		start  = time.Now()
		handle interface{}
	)
	defer func() {
		c.log.Debugw("Unsubscribe",
			"subject", sub.GetSubject(), "elapsed", time.Since(start).Seconds(),
			"handle", handle, "error", err,
		)
	}()

	if s, ok := sub.(*subscription); ok {
		return convertErr(s.Subscription.Unsubscribe())
	}

	return fmt.Errorf("invalid subscription type %s", reflect.TypeOf(sub).String())
}

// New - return new 'NATS' client for rpc and broadcast notifications
func New(log *zap.SugaredLogger, url, name string, maxReconnects int) *Client {
	return &Client{
		conn: newConn(url, nats.JSON_ENCODER, name, maxReconnects),
		log:  log,
	}
}
