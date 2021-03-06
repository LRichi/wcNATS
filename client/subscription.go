package client

import (
	"reflect"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Subscription interface {
	GetSubject() string
}

type subscription struct {
	*nats.Subscription
	log      *zap.SugaredLogger
	process  reflect.Value
	response reflect.Value
}

// GetSubject - return subject of subscription
func (s *subscription) GetSubject() string {
	return s.Subscription.Subject
}

// notify - implements the subscriber's notify
func (s *subscription) notify(values []reflect.Value) []reflect.Value {
	if s.process.Type().In(1).String() != values[2].FieldByName("Request").Type().String() {
		s.log.Debugw("wrong notify model", "subject", s.Subscription.Subject,
			"wait", s.process.Type().In(1), "got", values[2].Type().String(),
		)
	}

	var (
		start          = time.Now()
		responseValues []reflect.Value
		ctx            = createSession(values[2].FieldByName("Session").Interface().(SessionDTO))
		err            error
	)
	defer func() {
		s.log.Debugw("Notify", "elapsed", time.Since(start).Seconds(),
			"subject", s.Subscription.Subject,
			"request", values[2].Interface(), "error", err,
		)
	}()

	// calling the subscriber
	responseValues = s.process.Call([]reflect.Value{reflect.ValueOf(ctx), values[2].FieldByName("Request")})
	if !responseValues[0].IsZero() {
		err = responseValues[0].Interface().(error)
	}

	return nil
}

// call - implements the subscriber's call and the response to the client who created the call
func (s *subscription) call(values []reflect.Value) []reflect.Value {
	if s.process.Type().In(1).String() != values[2].FieldByName("Request").Type().String() {
		s.log.Debugw("wrong calling model", "subject", s.Subscription.Subject,
			"wait", s.process.Type().In(1), "got", values[2].Type().String(),
		)
	}

	var (
		start          = time.Now()
		responseValues []reflect.Value
		ctx            = createSession(values[2].FieldByName("Session").Interface().(SessionDTO))
		err            error
	)
	defer func() {
		s.log.Debugw("Call",
			"subject", s.Subscription.Subject, "elapsed", time.Since(start).Seconds(),
			"request", values[2].Field(1).Interface(), "response", responseValues[0].Interface(),
			"error", responseValues[1].Interface(), "reply error", err,
		)
	}()

	// calling the subscriber
	responseValues = s.process.Call([]reflect.Value{reflect.ValueOf(ctx), values[2].FieldByName("Request")})

	// creating structures for the response
	dtoValue := newResponseDTO(responseValues[0].Type())
	dtoValue.Field(0).Set(responseValues[0])

	// check process error
	if !responseValues[1].IsZero() {
		t := responseValues[1].Type().String()
		m := responseValues[1].Interface().(error).Error()
		var dto = ErrorDTO{
			Type:    &t,
			Message: &m,
		}
		dtoValue.Field(1).Set(reflect.ValueOf(dto))
	}

	// create error for reply
	var respVal = s.response.Call(
		[]reflect.Value{
			values[1],
			dtoValue.Addr(),
		},
	)

	if !respVal[0].IsZero() {
		err = respVal[0].Interface().(error)
	}

	return nil
}
