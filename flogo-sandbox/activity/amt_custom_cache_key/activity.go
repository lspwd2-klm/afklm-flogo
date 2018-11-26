package amt_custom_cache_key

import (
	"bytes"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"reflect"
)

const (
	OutputKey         = "CacheKey"
	InputCacheHeaders = "CacheHeaders"
	InputHeaders      = "headers_in"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	key := "any"

	headersToCache := context.GetInput(InputCacheHeaders)
	fmt.Println(headersToCache)
	if headersToCache != nil && reflect.ValueOf(headersToCache).Kind() == reflect.Slice {
		headerCacheArr := headersToCache.([]string)

		rawHeadersIn := context.GetInput(InputHeaders)
		fmt.Println(rawHeadersIn)
		if rawHeadersIn != nil && reflect.ValueOf(rawHeadersIn).Kind() == reflect.Map {

			headersMap := rawHeadersIn.(map[string]string)
			var sb bytes.Buffer

			for _, header := range headerCacheArr {
				sb.WriteString("/")
				passedHeader, present := headersMap[header]
				if !present {
					passedHeader = "*"
				}

				sb.WriteString(passedHeader)
			}

			key = sb.String()
		}
	}

	context.SetOutput(OutputKey, key)

	return true, nil
}
