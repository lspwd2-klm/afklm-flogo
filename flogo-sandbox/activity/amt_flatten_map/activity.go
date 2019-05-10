package amt_flatten_map

import (
	"bytes"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"sort"
)

const (
	InputKey  = "InputMap"
	OutputKey = "OutputValue"
)

var log = logger.GetLogger("amt-flatten-map")

// AmtFlattenMap is a stub for your Activity implementation
type AmtFlattenMap struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AmtFlattenMap{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *AmtFlattenMap) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *AmtFlattenMap) Eval(context activity.Context) (done bool, err error) {

	var buffer bytes.Buffer
	buffer.WriteString("")

	// do eval
	inTable := context.GetInput(InputKey).(map[string]string)
	if inTable != nil {
		// Sort keys in alphabetical order
		var keys []string
		for k := range inTable {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			buffer.WriteString(fmt.Sprintf("%s=%s;", k, inTable[k]))
		}
	}

	context.SetOutput(OutputKey, buffer.String())
	return true, nil
}
