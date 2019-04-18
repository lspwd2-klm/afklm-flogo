package amt_memorycache_get

// Ensure that you've run docker run --name flogo-redis -d redis prior running thees
// tests

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/lspwd2-klm/afklm-flogo/flogo-sandbox/afklmamt/memorycache"
	"io/ioutil"
	"testing"
	"time"
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

func TestBasicMissingGet(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput(CacheKey, "missing-by-definition")
	act.Eval(tc)

	if tc.GetOutput(CompleteOut).(bool) {
		t.Error("Evaluation shouldn't be successfully completed")
		t.Fail()
	}
}

func TestBasicGet(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	key := fmt.Sprintf("testkey-%d", time.Now().Unix())
	headers := make(map[string]string)
	headers["Accept"] = "application/json"

	memorychache.Set(key, memorychache.CachedHTTPResponse{201, headers, "Okay, buddy!"})

	tc.SetInput("CacheKey", key)
	act.Eval(tc)

	if tc.GetOutput(CacheCodeOut) != 201 {
		t.Error("Unexpeted code")
		t.Fail()
	}

	if tc.GetOutput(CacheBodyOut) != "Okay, buddy!" {
		t.Error("Unexpected body")
		t.Fail()
	}

	// Ignore headers for now.
}
