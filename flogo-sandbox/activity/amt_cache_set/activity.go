package amt_cache_set

import (
	"encoding"
	"encoding/json"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/go-redis/redis"
	"time"
)

const (
	HostName   = "RedisHost"
	DurationIn = "CacheDuration"

	CacheKey       = "CacheKey"
	CacheBodyIn    = "CacheBody"
	CacheHeadersIn = "CacheHeaders"
	CacheCodeIn    = "CacheCode"
)

var redisConn *redis.Client
var log = logger.GetLogger("amt-cache-get")

// AMTCacheSetActivity is a stub for your Activity implementation
type AMTCacheSetActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AMTCacheSetActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *AMTCacheSetActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *AMTCacheSetActivity) GetRedis(opts *redis.Options) *redis.Client {
	if redisConn == nil {
		redisConn = redis.NewClient(opts)
	}

	return redisConn
}

// Eval implements activity.Activity.Eval
func (a *AMTCacheSetActivity) Eval(context activity.Context) (done bool, err error) {

	cl := a.GetRedis(&redis.Options{
		Addr:     context.GetInput(HostName).(string),
		Password: "",
		DB:       0,
	})

	duration := 10 * time.Minute
	ctxDurationIn := context.GetInput(DurationIn)
	fmt.Println("Configured duration: ", ctxDurationIn)

	// This is, perhaps, an ugly check
	/*
		if !reflect.ValueOf(ctxDurationIn).IsNil() {
			uDur, err := time.ParseDuration(ctxDurationIn.(string))
			if err != nil {
				log.Error("Cannot parse supplied duration: ", ctxDurationIn)
				duration = uDur
			}
		}
	*/

	cacheKey := context.GetInput(CacheKey).(string)
	codeKey := cacheKey + ":code"
	//headersKey := cacheKey + ":headers"
	bodyKey := cacheKey + ":body"

	codeValue := context.GetInput(CacheCodeIn).(int)

	// TODO: Re-cast map[string]string to map[string]interface{}
	// Maybe just do a conversion.
	//headersValue := context.GetInput(CacheHeadersIn).(map[string]interface{})
	cacheBody := context.GetInput(CacheBodyIn)

	log.Info("Passing in duration=", duration)
	codeSet := cl.Set(codeKey, codeValue, duration)
	log.Info("Saved response code: ", codeSet.String())

	/*
		if headersValue != nil {
			cl.HMSet(headersKey, headersValue)
			cl.Expire(headersKey, duration)
		}
	*/

	if cacheBody != nil {
		_, ok := cacheBody.(string)
		if !ok {
			_, ok = cacheBody.(encoding.BinaryMarshaler)
			if !ok {
				// Marshal the string.
				strBody, encErr := json.Marshal(cacheBody)
				if encErr != nil {
					log.Error("Could not marshal JSON ", encErr)
				} else {
					cacheBody = string(strBody)
				}
			}
		}
		cmd := cl.Set(bodyKey, cacheBody, duration)
		log.Info("Results of setting body: ", cmd.String())
	}

	return true, nil
}
