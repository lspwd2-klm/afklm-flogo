package amt_map

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

func TestEvalMissing(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	mappingTable := map[string]string{
		"Key1": "Val1",
		"Key2": "Val2",
	}

	//check result attr
	tc1 := test.NewTestActivityContext(getActivityMetadata())
	tc1.SetInput(MappingTable, mappingTable)
	tc.SetInput(InputKey, "MissingKey")

	_, err := act.Eval(tc)
	if err != nil {
		t.Error("Failed to evaluate key absence: ", err)
		t.Fail()
		return
	}

	if tc.GetOutput(MappedOut).(bool) {
		t.Error("Mapping DID succeed, but was NOT expected to")
		t.Fail()
		return
	}

	if len(tc.GetOutput(MappedValueOut).(string)) > 0 {
		t.Error("Mapping should be nil.")
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

	mappingTable := map[string]string{
		"Key1": "Val1",
		"Key2": "Val2",
	}
	tc.SetInput(MappingTable, mappingTable)
	tc.SetInput(InputKey, "Key1")

	//setup attrs

	_, err := act.Eval(tc)
	if err != nil {
		t.Error("Failed to evaluate direct presence: ", err)
		t.Fail()
		return
	}

	if !tc.GetOutput(MappedOut).(bool) {
		t.Error("Mapping did not succeed, but was expected to")
		t.Fail()
		return
	}

	if tc.GetOutput(MappedValueOut) != "Val1" {
		t.Error("Mapping does not contain the expected value.")
	}

}
