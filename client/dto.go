package client

import (
	"context"
	"fmt"
	"reflect"
)

const (
	contextKeySession = "session"
	contextKeyService = "service"
	contextKeyMethod  = "method"
)

type SessionDTO struct {
	Session *string
	Service *string
	Method  *string
}

type ErrorDTO struct {
	Type    *string
	Message *string
}

func (e ErrorDTO) Error() string {
	if e.Message == nil {
		panic("is nil")
	}

	return fmt.Sprintf("%s", *e.Message)
}

// newRequestDTO - creates a transport data structure in memory to receive a request
func newRequestDTO(req reflect.Type) reflect.Value {
	//
	// schema:
	//  type DTO struct {
	//      Context TypeOf()
	//      Request TypeOf()
	// }
	//
	var (
		sessionDTO SessionDTO
		s          = reflect.TypeOf(sessionDTO)
	)

	fieldsOfDtoStruct := reflect.StructOf([]reflect.StructField{
		{
			Name:      "Session",
			PkgPath:   "",
			Type:      s,
			Tag:       "",
			Offset:    0,
			Index:     nil,
			Anonymous: false,
		},
		{
			Name:      "Request",
			PkgPath:   "",
			Type:      req,
			Tag:       "",
			Offset:    1,
			Index:     nil,
			Anonymous: false,
		},
	})

	dtoStruct := reflect.StructOf([]reflect.StructField{
		{
			Name:      "DTO",
			PkgPath:   "",
			Type:      fieldsOfDtoStruct,
			Tag:       "",
			Offset:    0,
			Index:     nil,
			Anonymous: false,
		},
	})

	dto := reflect.New(dtoStruct)

	return dto.Elem().Field(0)
}

// newResponseDTO - creates a transport data structure in memory to receive a response
func newResponseDTO(resp reflect.Type) reflect.Value {
	//
	// schema:
	//  type DTO struct {
	//      Response TypeOf()
	//      Error TypeOf()
	// }
	//

	var (
		dtoError ErrorDTO
		err      = reflect.TypeOf(dtoError)
	)

	fieldsOfDtoStruct := reflect.StructOf([]reflect.StructField{
		{
			Name:      "Response",
			PkgPath:   "",
			Type:      resp,
			Tag:       "",
			Offset:    0,
			Index:     nil,
			Anonymous: false,
		},
		{
			Name:      "Error",
			PkgPath:   "",
			Type:      err,
			Tag:       "",
			Offset:    1,
			Index:     nil,
			Anonymous: false,
		},
	})

	dtoStruct := reflect.StructOf([]reflect.StructField{
		{
			Name:      "DTO",
			PkgPath:   "",
			Type:      fieldsOfDtoStruct,
			Tag:       "",
			Offset:    0,
			Index:     nil,
			Anonymous: false,
		},
	})

	dto := reflect.New(dtoStruct)

	return dto.Elem().Field(0)
}

// getSession - return session context keys
func getSession(ctx context.Context) (dto SessionDTO) {
	if ctx != nil {
		var v interface{}
		if v = ctx.Value(contextKeySession); v != nil {
			if s, ok := v.(string); ok {
				dto.Session = &s
			}
		}
		if v = ctx.Value(contextKeyService); v != nil {
			if s, ok := v.(string); ok {
				dto.Service = &s
			}
		}
		if v = ctx.Value(contextKeyMethod); v != nil {
			if s, ok := v.(string); ok {
				dto.Method = &s
			}
		}
	}

	return
}

// createSession - create context with value if DTO's data exists, or empty context
func createSession(dto SessionDTO) (ctx context.Context) {
	ctx = context.Background()
	if dto.Session != nil {
		ctx = context.WithValue(ctx, contextKeySession, *dto.Session)
	}

	if dto.Service != nil {
		ctx = context.WithValue(ctx, contextKeyService, *dto.Service)
	}

	if dto.Method != nil {
		ctx = context.WithValue(ctx, contextKeyMethod, *dto.Method)
	}

	return
}
