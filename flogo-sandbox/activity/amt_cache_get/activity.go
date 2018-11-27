package amt_cache_get

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/go-redis/redis"
)

const (
	HostName        = "RedisHost"
	CacheKey        = "CacheKey"
	CacheBodyOut    = "CacheBody"
	CacheHeadersOut = "CacheHeaders"
	CacheCodeOut    = "CacheCode"
	CompleteOut     = "Complete"
)

var redisConnn *redis.Client

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

	// do eval

	if redisConnn == nil {
		redisConnn = redis.NewClient(&redis.Options{
			Addr:     context.GetInput(HostName).(string),
			Password: "",
			DB:       0,
		})
	}

	cacheKey := context.GetInput(CacheKey).(string)
	codeKey := cacheKey + ".code"
	headersKey := cacheKey + ".headers"
	bodyKey := cacheKey + ".body"

	cacheData := redisConnn.MGet(codeKey, headersKey, bodyKey).Val()

	respCode := 0

	if cacheData[0] != nil {
		respCode = cacheData[0].(int)
	}

	context.SetOutput(CompleteOut, respCode != 0)
	context.SetOutput(CacheCodeOut, respCode)
	context.SetOutput(CacheHeadersOut, cacheData[1])
	context.SetOutput(CacheBodyOut, cacheData[2])

	return true, nil
}
