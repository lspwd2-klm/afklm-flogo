package amt_memorycache_set

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/lspwd2-klm/afklm-flogo/flogo-sandbox/afklmamt/memorycache"
	"io/ioutil"
	"testing"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestSetting(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	cacheKey := fmt.Sprintf("cache-set-%d", time.Now().Unix())

	tc.SetInput(CacheKey, cacheKey)
	tc.SetInput(CacheCodeIn, 202)
	tc.SetInput(CacheHeadersIn, nil)
	tc.SetInput(CacheBodyIn, "Setting Value")

	done, error := act.Eval(tc)
	if !done {
		t.Error("Activity must be completed")
		t.Fail()
	}

	if error != nil {
		t.Error("Activity must complete without an error")
		t.Fail()
	}

	obj, present := memorychache.Get(cacheKey)
	if !present {
		t.Error("Value must be pesent in the memory cache")
		t.Fail()
	}

	value := obj.(memorychache.CachedHTTPResponse)

	if value.Body != "Setting Value" {
		t.Error(fmt.Sprintf("Wrong value got cached: %s", value.Body))
		t.Fail()
	}

	if value.Code != 202 {
		t.Error(fmt.Sprintf("Wrong code got cached: %d", value.Code))
		t.Fail()
	}
}
