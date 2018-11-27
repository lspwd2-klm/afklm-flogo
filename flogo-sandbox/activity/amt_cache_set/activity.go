package amt_cache_set

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/go-redis/redis"
)

const (
	HostName       = "RedisHost"
	CacheKey       = "CacheKey"
	CacheBodyIn    = "CacheBody"
	CacheHeadersIn = "CacheHeaders"
	CacheCodeIn    = "CacheCode"
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

	cacheCode := context.GetInput(CacheCodeIn).(int)
	cacheHeaders := context.GetInput(CacheHeadersIn)
	cacheBody := context.GetInput(CacheBodyIn)

	redisConnn.MSet(codeKey, cacheCode, headersKey, cacheHeaders, bodyKey, cacheBody)

	return true, nil
}
