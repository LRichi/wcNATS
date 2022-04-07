package client

import (
	"fmt"
	"reflect"
)

const (
	gotFuncCallDesc   = "func(context.Context,*struct)(*struct,error)"
	gotFuncNotifyDesc = "func(context.Context,*struct)(error)"
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

func validateHandleOfCall(handle interface{}) error {
	var t = reflect.TypeOf(handle)

	switch {
	// check num parameters
	case t.NumIn() != 2:
		return fmt.Errorf("incoming parameters != 2, use %s", gotFuncCallDesc)
	case t.NumOut() != 2:
		return fmt.Errorf("returned parameters != 2, use %s", gotFuncCallDesc)

	// check first incoming parameter
	case t.In(0).Kind() != reflect.Interface:
		return fmt.Errorf("firts parameter is not interface, use %s", gotFuncCallDesc)
	case t.In(0).String() != "context.Context":
		return fmt.Errorf("firts parameter is not context.Context, use %s", gotFuncCallDesc)

	// check second parameter
	case t.In(1).Kind() != reflect.Ptr:
		return fmt.Errorf("second parameter is not PTR, use %s", gotFuncCallDesc)
	case t.In(1).Elem().Kind() != reflect.Struct:
		return fmt.Errorf("second parameter is not *struct, use %s", gotFuncCallDesc)

	// Check first returned parameter
	case t.Out(0).Kind() != reflect.Ptr:
		return fmt.Errorf("firts parameter is not PTR, use %s", gotFuncCallDesc)
	case t.Out(0).Elem().Kind() != reflect.Struct:
		return fmt.Errorf("firts parameter is not *struct, use %s", gotFuncCallDesc)

	// Check second returned parameter
	case t.Out(1).Kind() != reflect.Interface:
		return fmt.Errorf("firts parameter is not interface, use %s", gotFuncCallDesc)
	case t.Out(1).String() != "error":
		return fmt.Errorf("second parameter is not interface(error), use %s", gotFuncCallDesc)
	}

	return nil
}

func isRequestHandle(handle interface{}) (bool, error) {
	var t = reflect.TypeOf(handle)
	switch {
	case t.Kind() != reflect.Func:
		return false, fmt.Errorf("is not func, use %s or %s", gotFuncCallDesc, gotFuncNotifyDesc)
	case t.NumIn() == 2 && t.NumOut() == 2:
		return true, nil
	case t.NumIn() == 2 && t.NumOut() == 1:
		return true, nil
	default:
		return false, fmt.Errorf("unable to determine subscription type, use %s or %s", gotFuncCallDesc, gotFuncNotifyDesc)
	}
}

func validateHandleOfNotify(handle interface{}) error {
	var t = reflect.TypeOf(handle)
	switch {
	// check num parameters
	case t.NumIn() != 2:
		return fmt.Errorf("incoming parameters != 2, use %s", gotFuncNotifyDesc)
	case t.NumOut() != 1:
		return fmt.Errorf("returned parameters != 1, use %s", gotFuncNotifyDesc)

	// check first incoming parameter
	case t.In(0).Kind() != reflect.Interface:
		return fmt.Errorf("firts parameter is not interface, use %s", gotFuncNotifyDesc)
	case t.In(0).String() != "context.Context":
		return fmt.Errorf("firts parameter is not context.Context, use %s", gotFuncNotifyDesc)

	// check second parameter
	case t.In(1).Kind() != reflect.Ptr:
		return fmt.Errorf("second parameter is not PTR, use %s", gotFuncNotifyDesc)
	case t.In(1).Elem().Kind() != reflect.Struct:
		return fmt.Errorf("second parameter is not *struct, use %s", gotFuncNotifyDesc)

	// Check first returned parameter
	case t.Out(0).Kind() != reflect.Interface:
		return fmt.Errorf("firts parameter is not interface, use %s", gotFuncNotifyDesc)
	case t.Out(0).String() != "error":
		return fmt.Errorf("second parameter is not interface(error), use %s", gotFuncNotifyDesc)
	}

	return nil
}
