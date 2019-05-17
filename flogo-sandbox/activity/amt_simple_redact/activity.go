package amt_simple_redact

import (
	"encoding/json"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/go-redis/redis"
	"reflect"
	"regexp"
)

const (
	Pattern      = "Pattern"
	Replacements = "Replacements"
	ObjectIn     = "ObjectIn"
	ObjectOut    = "ObjectOut"
)

var log = logger.GetLogger("amt-cache-get")
var redisConn *redis.Client

// AmtCacheActivity is a stub for your Activity implementation
type AmtSimpleRedactActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AmtSimpleRedactActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *AmtSimpleRedactActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func redact(obj_in interface{}, p *regexp.Regexp, r string) interface{} {
	if obj_in == nil {
		return nil
	} else if valInt, okInt := obj_in.(int64); okInt {
		return valInt
	} else if valInt32, okInt := obj_in.(int32); okInt {
		return valInt32
	} else if floatVal, okFloat := obj_in.(float64); okFloat {
		return floatVal
	} else if floatVal32, okFloat := obj_in.(float32); okFloat {
		return floatVal32
	} else if jsonNum, okJsonNum := obj_in.(json.Number); okJsonNum {
		return jsonNum
	} else if boolVal, okBool := obj_in.(bool); okBool {
		return boolVal
	} else if strVal, okStr := obj_in.(string); okStr {
		return p.ReplaceAllString(strVal, r)
	} else if mapVal, okMap := obj_in.(map[string]interface{}); okMap {
		return redact_map(mapVal, p, r)
	} else if arrVal, okArr := obj_in.([]interface{}); okArr {
		return redact_array(arrVal, p, r)
	} else {
		return fmt.Sprintf("---no-procedure-for-%s---", reflect.TypeOf(obj_in).String())
	}
}

func redact_array(arr_in []interface{}, p *regexp.Regexp, r string) []interface{} {
	arr_out := make([]interface{}, len(arr_in))

	for i, obj := range arr_in {
		arr_out[i] = redact(obj, p, r)
	}

	return arr_out
}

func redact_map(map_in map[string]interface{}, p *regexp.Regexp, r string) map[string]interface{} {
	map_out := make(map[string]interface{})

	for k, v := range map_in {
		map_out[k] = redact(v, p, r)
	}

	return map_out
}

// Eval implements activity.Activity.Eval
func (a *AmtSimpleRedactActivity) Eval(context activity.Context) (done bool, err error) {

	pattern := context.GetInput(Pattern).(string)
	replacement := context.GetInput(Replacements).(string)
	obj_in := context.GetInput(ObjectIn)

	exec_p := regexp.MustCompile(pattern)
	context.SetOutput(ObjectOut, redact(obj_in, exec_p, replacement))

	return true, nil
}
