package nats

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	*conn
	log *zap.SugaredLogger
}

// Subscribe - subscribe handle for remote call or notify
//
// use:
// for rpc:    func(context.Context,*struct)(*struct,error)
// fpr notify: func(context.Context,*struct)(error)
func (c *Client) Subscribe(subject fmt.Stringer, handle interface{}) (Subscription, error) {
	var (
		start = time.Now()
		err   error
	)
	defer func() {
		c.log.Debug("Subscribe", "elapsed", time.Since(start).Seconds(), "subject", subject,
			"handle", handle, "error", err,
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
		reqDTO            = newRequestDTO(reflect.TypeOf(handle).In(1))
		typeOfNewHandleIn = []reflect.Type{reflect.TypeOf(""), reflect.TypeOf(""), reqDTO.Type()}
		typeOfNewHandle   = reflect.FuncOf(typeOfNewHandleIn, nil, false)
	)

	var valueOfNewHandle reflect.Value
	if isRequest {
		if err = validateHandleOfCall(handle); err != nil {
			return nil, fmt.Errorf("invalid handle: %w", err)
		}
		valueOfNewHandle = reflect.MakeFunc(typeOfNewHandle, sub.call)
	} else {
		if err = validateHandleOfNotify(handle); err != nil {
			return nil, fmt.Errorf("invalid handle: %w", err)
		}
		valueOfNewHandle = reflect.MakeFunc(typeOfNewHandle, sub.notify)
	}

	sub.Subscription, err = c.subscribe(subject.String(), valueOfNewHandle.Interface())
	if err != nil {
		return nil, convertErr(err)
	}

	return sub, nil
}

// Request - a remote procedure call is created
//
// use: func(context.Context,*struct)(*struct,error)
func (c *Client) Request(ctx context.Context, subject string, request, response interface{}) (err error) {
	start := time.Now()
	defer func() {
		c.log.Debug("Subscribe", "elapsed", time.Since(start).Seconds(),
			"subject", subject, "request", request, "response", response, "error", err,
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
	if err = c.Request(ctx, subject, reqDTO.Interface(), respDTO.Addr().Interface()); err != nil {
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
