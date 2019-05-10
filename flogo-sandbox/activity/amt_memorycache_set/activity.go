package amt_memorycache_set

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/lspwd2-klm/afklm-flogo/flogo-sandbox/afklmamt/memorycache"
)

const (
	DurationIn = "CacheDuration"

	CacheKey       = "CacheKey"
	CacheBodyIn    = "CacheBody"
	CacheHeadersIn = "CacheHeaders"
	CacheCodeIn    = "CacheCode"
)

var log = logger.GetLogger("amt_memorycache_set")

// AMTCacheSetActivity is a stub for your Activity implementation
type AMTMemoryCacheSetActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AMTMemoryCacheSetActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *AMTMemoryCacheSetActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *AMTMemoryCacheSetActivity) Eval(context activity.Context) (done bool, err error) {

	cacheKey := context.GetInput(CacheKey).(string)
	cacheCode := context.GetInput(CacheCodeIn).(int)
	cacheBody := context.GetInput(CacheBodyIn)
	cacheHeaders := context.GetInput(CacheHeadersIn).(map[string]string)
	cacheDuration := context.GetInput(DurationIn).(int)

	memorychache.Set(cacheKey, memorychache.CachedHTTPResponse{cacheCode, cacheHeaders, cacheBody}, cacheDuration)

	log.Info(fmt.Sprintf("Set the value for key: %s", cacheKey))

	return true, nil
}
