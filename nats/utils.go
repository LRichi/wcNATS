package nats

import (
	"fmt"
	"reflect"
)

// validateModel - checks a model for a request or response
func validateModel(model interface{}) error {
	if model == nil {
		return fmt.Errorf("is not may be nil, use *struct")
	}

	var t = reflect.TypeOf(model)
	switch {
	case t.Kind() != reflect.Ptr:
		return fmt.Errorf("value is not ptr)")
	case t.Elem().Kind() != reflect.Struct:
		return fmt.Errorf("value is not struct")
	}

	return nil
}
