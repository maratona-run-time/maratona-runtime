package utils

import "context"

type QueryClient interface {
	Query(context.Context, interface{}, map[string]interface{}) error
}
