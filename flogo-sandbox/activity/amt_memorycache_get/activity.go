package amt_memorycache_get

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/go-redis/redis"
	"github.com/lspwd2-klm/afklm-flogo/flogo-sandbox/afklmamt/memorycache"
)

const (
	CacheKey       = "CacheKey"
	StaleTolerance = "StaleTolerance"
	// TODO: add the stale tolerance.
	CacheBodyOut    = "CacheBody"
	CacheHeadersOut = "CacheHeaders"
	CacheCodeOut    = "CacheCode"
	CompleteOut     = "Complete"
	StaleOut        = "Stale"
)

var log = logger.GetLogger("amt-cache-get")
var redisConn *redis.Client

// AmtCacheActivity is a stub for your Activity implementation
type AmtMemoryCacheGetActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AmtMemoryCacheGetActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *AmtMemoryCacheGetActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *AmtMemoryCacheGetActivity) Eval(context activity.Context) (done bool, err error) {

	cacheKey := context.GetInput(CacheKey).(string)
	obj_holder, found := memorychache.Get(cacheKey)

	if found {
		obj := obj_holder.(memorychache.CachedHTTPResponse)

		context.SetOutput(CacheCodeOut, obj.Code)
		context.SetOutput(CacheBodyOut, obj.Body)
		context.SetOutput(CacheHeadersOut, obj.Headers)

		context.SetOutput(CompleteOut, true)
	} else {
		context.SetOutput(CompleteOut, false)
	}

	return true, nil
}
