package utils

import (
	"context"
	"reflect"
	"testing"
)

type QueryClient interface {
	Query(context.Context, interface{}, map[string]interface{}) error
}

type GraphqlMock struct {
	Test      *testing.T
	Object    interface{}
	Variables map[string]interface{}
}

func (gm GraphqlMock) Query(ctx context.Context, info interface{}, variables map[string]interface{}) error {
	reflect.ValueOf(info).Elem().Set(reflect.ValueOf(gm.Object))
	if !reflect.DeepEqual(gm.Variables, variables) {
		gm.Test.Errorf("Expect request variables to be %v, received %v", gm.Variables, variables)
	}
	return nil
}
