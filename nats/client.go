package nats

import (
	"context"
	"fmt"
	"reflect"
)

type Client struct {
	*conn
}

// Request - a remote procedure call is created
func (n *Client) Request(ctx context.Context, subject string, request, response interface{}) (err error) {
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
	if err = n.Request(ctx, subject, reqDTO.Interface(), respDTO.Addr().Interface()); err != nil {
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
