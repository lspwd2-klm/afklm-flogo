package amt_simple_redact

// Ensure that you've run docker run --name flogo-redis -d redis prior running thees
// tests

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

func TestCreate(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestBasicRedaction(t *testing.T) {
	var dat map[string]interface{}
	byt, _ := ioutil.ReadFile("./ref.json")

	if err := json.Unmarshal(byt, &dat); err != nil {
		t.Fail()
	}

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput(ObjectIn, dat)
	tc.SetInput(Pattern, "_REPL_")
	tc.SetInput(Replacements, "*****")

	done, err := act.Eval(tc)
	if !done {
		t.Error("Should be done")
		t.Fail()
	}

	if err != nil {
		t.Error("Error not expected")
		t.Fail()
	}

	out := tc.GetOutput(ObjectOut).(map[string]interface{})

	fmt.Print(out)
}

func TestRootArrayRedaction(t *testing.T) {
	var dat interface{}
	byt, _ := ioutil.ReadFile("./complex.json")

	if err := json.Unmarshal(byt, &dat); err != nil {
		t.Error(fmt.Sprintf("Could not read json %s", err))
		t.Fail()
	}

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput(ObjectIn, dat)
	tc.SetInput(Pattern, "mega.secret.domain.com")
	tc.SetInput(Replacements, "*****.com")

	done, err := act.Eval(tc)
	if !done {
		t.Error("Should be done")
		t.Fail()
	}

	if err != nil {
		t.Error("Error not expected")
		t.Fail()
	}

	out := tc.GetOutput(ObjectOut)

	fmt.Print(out)
}

func TestRootArrayRedaction2(t *testing.T) {
	var dat interface{}
	byt, _ := ioutil.ReadFile("./complex2.json")

	if err := json.Unmarshal(byt, &dat); err != nil {
		t.Error(fmt.Sprintf("Could not read json %s", err))
		t.Fail()
	}

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput(ObjectIn, dat)
	tc.SetInput(Pattern, "mega.secret.domain.com")
	tc.SetInput(Replacements, "*****.com")

	done, err := act.Eval(tc)
	if !done {
		t.Error("Should be done")
		t.Fail()
	}

	if err != nil {
		t.Error("Error not expected")
		t.Fail()
	}

	out := tc.GetOutput(ObjectOut)

	fmt.Print(out)
}
