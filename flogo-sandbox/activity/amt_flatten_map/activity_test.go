package amt_flatten_map

import (
	"io/ioutil"
	"testing"

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

func TestEvalSorting(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	mappingTable := map[string]string{
		"Key2": "Val2",
		"Key1": "Val1",
	}

	//check result attr
	tc.SetInput(InputKey, mappingTable)

	_, err := act.Eval(tc)
	if err != nil {
		t.Error("Failed to evaluate key absence: ", err)
		t.Fail()
		return
	}

	output := tc.GetOutput(OutputKey)
	if output != "Key1=Val1;Key2=Val2;" {
		t.Error("Unexpected key result", output)
		t.Fail()
		return
	}

}

func TestEvalWithEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	_, err := act.Eval(tc)
	if err != nil {
		t.Error("Failed to evaluate key absence: ", err)
		t.Fail()
		return
	}

	output := tc.GetOutput(OutputKey)
	if output != "" {
		t.Error("Unexpected key result", output)
		t.Fail()
		return
	}

}
