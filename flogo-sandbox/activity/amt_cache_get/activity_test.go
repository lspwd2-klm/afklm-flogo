package amt_cache_get

// Ensure that you've run docker run --name flogo-redis -d redis prior running thees
// tests

import (
	"github.com/go-redis/redis"
	"io/ioutil"
	"testing"
	"time"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
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

func TestCacheMiss(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput(HostName, "192.168.99.100:6379")
	tc.SetInput(CacheKey, "NonExistingKey")

	_, err := act.Eval(tc)
	if err != nil {
		t.Error("Unable to run test miss")
		t.Error(err)
		t.Fail()
		return
	}

	complete := tc.GetOutput(CompleteOut).(bool)
	if complete {
		t.Error("Missing cache must be empty")
		t.Fail()
		return
	}

}

func TestCacheHit(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput(HostName, "192.168.99.100:6379")
	tc.SetInput(CacheKey, "MyKey")

	amtAct := act.(*AmtCacheActivity)
	redisOpt := &redis.Options{
		Addr:     "192.168.99.100",
		Password: "",
		DB:       0,
	}

	cl := amtAct.GetRedis(redisOpt)
	duration := time.Duration(10) * time.Minute
	headers := map[string]interface{}{
		"Content-Type": "application/json",
		"AFKLM-Market": "NL",
	}

	cl.Set("MyKey:code", 201, duration)
	cl.HMSet("MyKey:headers", headers)
	cl.Expire("MyKey:headers", duration)

	cl.Set("MyKey:body", "{a=b}", duration)

	_, err := act.Eval(tc)
	if err != nil {
		t.Error("Failed to evaluate action on all parameters present")
		t.Error(err)
		t.Fail()
		return
	}

	gotOutput := tc.GetOutput(CompleteOut).(bool)
	if !gotOutput {
		t.Error("Failed to retrieve data from Redis: data is not there")
		t.Fail()
		return
	}

	// Now, let's check the results.
	outCode := tc.GetOutput(CacheCodeOut)
	if outCode != 201 {
		t.Error("Did not receive expected code from Redis, got this: ", outCode)
		t.Fail()
		return
	}

	outHeaders := tc.GetOutput(CacheHeadersOut).(map[string]string)
	if outHeaders["Content-Type"] != "application/json" {
		t.Error("Missing / unexpected content type")
		t.Fail()
	}
	if outHeaders["AFKLM-Market"] != "NL" {
		t.Error("Missing / unexpected AFKLM market")
		t.Fail()
	}

	if tc.GetOutput(CacheBodyOut) != "{a=b}" {
		t.Error("Missing / unexpected body")
		t.Fail()
	}
}

/*
func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs

	act.Eval(tc)

	//check result attr
}
*/
