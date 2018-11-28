package amt_cache_set

import (
	"encoding/json"
	"fmt"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"io/ioutil"
	"testing"

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

type SampleMessage struct {
	FieldA string
}

func TestSettingTheValue(t *testing.T) {

	message := &SampleMessage{
		FieldA: "valueb",
	}
	//slcD := []string{"apple", "peach", "pear"}

	//message := &response1{
	//	Page:   1,
	//	Fruits: []string{"apple", "peach", "pear"},
	//}

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput(HostName, "192.168.99.100:6379")
	tc.SetInput(CacheKey, "MyWriteKey")

	tc.SetInput(CacheCodeIn, 200)
	tc.SetInput(CacheHeadersIn, map[string]interface{}{
		"Content-Type": "application/json",
		"AFKLM-Market": "NL",
		"Foo":          5,
	})

	output, encErr := json.Marshal(message)
	if encErr != nil {
		t.Error("Could not marshal message")
		t.Error(encErr)
		t.Fail()
		return
	}

	fmt.Println("Formatted JSON object:")
	fmt.Println(string(output))

	tc.SetInput(CacheBodyIn, message)

	//check result attr
	_, err := act.Eval(tc)
	if err != nil {
		t.Error("Failed to set values")
		t.Error(err)
		t.Fail()
		return
	}

	fmt.Println("Check that the keys are set in Redis")

	amtAct := act.(*AMTCacheSetActivity)
	cl := amtAct.GetRedis(nil)
	fmt.Println(cl.TTL("MyWriteKey:body").String())
	fmt.Println(cl.TTL("ref").String())
}
