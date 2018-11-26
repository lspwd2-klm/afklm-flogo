package amt_custom_cache_key

import (
	"encoding/json"
	"fmt"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"io/ioutil"
	"testing"
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

func TestCreatingOnNoInputs(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	_, err := act.Eval(tc)
	if err != nil {
		t.Error("Unable to run on no inputs")
		t.Fail()
	}

	output := tc.GetOutput(OutputKey)
	if output != "any" {
		t.Error("Key not set any on no inputs")
		t.Fail()
	}
}

func TestCreatingOnFullConfig(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	jsonText := "{\"headers\":[\"Accept\",\"AFKLM-Market\",\"AFKLM-Channel\"]}"

	var jsObj interface{}

	jsErr := json.Unmarshal([]byte(jsonText), &jsObj)
	if jsErr != nil {
		t.Error("Cannot parse JSON", jsErr)
		t.Fail()
	}

	//fmt.Println(jsObj)
	tc.SetInput(InputCacheHeaders, jsObj)

	tc.SetInput(InputHeaders, map[string]string{
		"Accept":       "application/json",
		"Dummy":        "Dummy headers",
		"AFKLM-Market": "NL",
	})

	_, err := act.Eval(tc)
	if err != nil {
		t.Error("Unable to run on full config")
		t.Fail()
	}

	output := tc.GetOutput(OutputKey)
	if output != "/application/json/NL/*" {
		t.Error(fmt.Sprintf("Unexpected value received: %s", output))
		t.Fail()
	}
}

/*
func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

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
