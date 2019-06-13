package bin_time

import (
	"errors"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"time"
)

var log = logger.GetLogger("bin-date-action")

const (
	InputJSONTimestamp = "json"
	InputTimestamp     = "timestamp"
	InputIntervalType  = "interval"
	InputIntervalSpan  = "intervalSpan"
	OutputType         = "type"
	OutputAnchor       = "anchor"
	OutputIndex        = "index"
	OutputOffset       = "offset"
	OutputStart        = "start"
	OutputSpan         = "span"
	OutputEnd          = "end"
)

var the24h = 24 * time.Hour

type AmtBinTime struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AmtBinTime{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *AmtBinTime) Metadata() *activity.Metadata {
	return a.metadata
}

// Bins to the point-in-time, using the specified interval span in minutes.
func binTimeSlot(point time.Time, span time.Duration, ctx activity.Context) {
	ctx.SetOutput(OutputType, "TimeSlot")

	hours, minutes, seconds := point.Clock()
	year, month, date := point.Date()

	// Getting the interval start
	intervalAnchor := time.Date(year, month, date, 0, 0, 0, 0, point.Location())
	ctx.SetOutput(OutputAnchor, intervalAnchor)

	secondsAtPoint := seconds + minutes*60 + hours*3600

	intervalIndex := secondsAtPoint / int(span.Seconds())
	ctx.SetOutput(OutputIndex, intervalIndex)
	ctx.SetOutput(OutputSpan, span)

	offset := time.Duration(intervalIndex) * span
	ctx.SetOutput(OutputOffset, offset)

	intervalStart := intervalAnchor.Add(offset)
	ctx.SetOutput(OutputStart, intervalStart)

	intervalEnd := intervalStart.Add(span)
	ctx.SetOutput(OutputEnd, intervalEnd)

}

// Bins to the point-in-time, using the specified interval span in minutes.
func binDate(point time.Time, ctx activity.Context) {
	ctx.SetOutput(OutputType, "Day")

	year, month, date := point.Date()

	// Getting the interval start
	intervalStart := time.Date(year, month, date, 0, 0, 0, 0, point.Location())
	ctx.SetOutput(OutputStart, intervalStart)

	log.Debug("Calculated interval start: %s", intervalStart)

	ctx.SetOutput(OutputSpan, the24h)

	intervalEnd := time.Date(year, month, date, 23, 59, 59, 999, point.Location())
	ctx.SetOutput(OutputEnd, intervalEnd)
}

// Eval implements activity.Activity.Eval
func (a *AmtBinTime) Eval(context activity.Context) (done bool, err error) {
	var uDate time.Time

	jsonInput := context.GetInput(InputJSONTimestamp)
	timestampInput := context.GetInput(InputTimestamp)

	if jsonInput != nil {
		jsonInputStr, jsonIsStr := jsonInput.(string)
		if !jsonIsStr {
			log.Error(fmt.Sprintf("String expected for JSON input"))
			return false, errors.New("JSON input date must be string")
		}

		var parseErr error

		uDate, parseErr = time.Parse(time.RFC3339, jsonInputStr)
		if parseErr != nil {
			log.Error(fmt.Sprintf("Cannot parse string %s: %s", jsonInputStr, parseErr))
			return false, parseErr
		}
	} else if timestampInput != nil {
		var castOk bool

		if uDate, castOk = timestampInput.(time.Time); !castOk {
			log.Error("Cannot cast timestamp input to time.Time")
			return false, errors.New("Cannot cast timestamp input to time.Time")
		}
	}

	uInterval := context.GetInput(InputIntervalType).(string)
	switch uInterval {
	case "TimeSlot":
		if uSpan, castOk := context.GetInput(InputIntervalSpan).(string); castOk {
			uSpanDuration, intervalErr := time.ParseDuration(uSpan)
			if intervalErr == nil {
				binTimeSlot(uDate, uSpanDuration, context)
				return true, nil
			} else {
				log.Error(fmt.Sprintf("Failed to parse supplied time interval: %s.", uSpan))
				return false, errors.New("time slot format is not recognized")
			}
		} else {
			log.Error("Span must be present when interval is set to TimeSlot")
			return false, errors.New("time slot span is not present")
		}
	case "Day":
		binDate(uDate, context)
		return true, nil
	default:
		return false, errors.New(fmt.Sprintf("Unknown interval type: %s", uInterval))
	}
}
