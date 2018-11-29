package amt_map

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

const (
	InputKey       = "Key"
	MappingTable   = "MappingTable"
	MappedOut      = "Mapped"
	MappedValueOut = "MappedValue"
)

//var log = logger.GetLogger("amt-custom-cache-key")

// AmtMap is a stub for your Activity implementation
type AmtMap struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AmtMap{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *AmtMap) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *AmtMap) Eval(context activity.Context) (done bool, err error) {

	// do eval
	inTable := context.GetInput(MappingTable).(map[string]string)
	inKey := context.GetInput(InputKey).(string)

	val, found := inTable[inKey]
	if found {
		context.SetOutput(MappedValueOut, val)
	}

	context.SetOutput(MappedOut, found)
	return true, nil
}
