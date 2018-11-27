package amt_cache_get

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/go-redis/redis"
	"strconv"
)

const (
	HostName        = "RedisHost"
	CacheKey        = "CacheKey"
	CacheBodyOut    = "CacheBody"
	CacheHeadersOut = "CacheHeaders"
	CacheCodeOut    = "CacheCode"
	CompleteOut     = "Complete"
)

var log = logger.GetLogger("amt-cache-get")
var redisConn *redis.Client

// AmtCacheActivity is a stub for your Activity implementation
type AmtCacheActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AmtCacheActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *AmtCacheActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *AmtCacheActivity) GetRedis(opts *redis.Options) *redis.Client {
	if redisConn == nil {
		redisConn = redis.NewClient(opts)
	}

	return redisConn
}

// Eval implements activity.Activity.Eval
func (a *AmtCacheActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval

	cl := a.GetRedis(&redis.Options{
		Addr:     context.GetInput(HostName).(string),
		Password: "",
		DB:       0,
	})

	log.Info("Connected to Redis")

	cacheKey := context.GetInput(CacheKey).(string)
	codeKey := cacheKey + ":code"
	headersKey := cacheKey + ":headers"
	bodyKey := cacheKey + ":body"

	log.Info("Retrieving these keys:", codeKey, " ", headersKey, ", ", bodyKey)

	cacheData := cl.MGet(codeKey, bodyKey).Val()
	log.Info("Got length of data: ", len(cacheData))
	respCode := 0

	if len(cacheData) > 0 {
		log.Info("Okay, we have some in cache data.")
		if cacheData[0] != nil {
			parsedCode, err := strconv.Atoi(cacheData[0].(string))
			if err == nil {
				respCode = parsedCode
			}
		}

		log.Info("So, the response code is ", respCode)

		context.SetOutput(CacheCodeOut, respCode)
		context.SetOutput(CacheBodyOut, cacheData[1])

		// Headers will require fetching the map
		cachedHeaders := cl.HGetAll(headersKey).Val()
		if cachedHeaders != nil {
			context.SetOutput(CacheHeadersOut, cachedHeaders)
		}
	}

	context.SetOutput(CompleteOut, respCode != 0)

	return true, nil
}
