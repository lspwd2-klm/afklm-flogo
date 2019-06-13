package bin_time

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math"
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

func TestAssignmentOfDate(t *testing.T) {
	dt := time.Date(2019, time.Month(6), 7, 9, 10, 11, 12, time.UTC)

	year, month, date := dt.Date()

	assert.Equal(t, 2019, year)
	assert.Equal(t, time.Month(6), month)
	assert.Equal(t, 7, date)

	hour, min, sec := dt.Clock()

	assert.Equal(t, 9, hour)
	assert.Equal(t, 10, min)
	assert.Equal(t, 11, sec)
}

func TestAssignmentFromParsedJSON(t *testing.T) {
	dtOrig, err := time.Parse(time.RFC3339, "2019-06-06T23:32:53Z")
	assert.Nil(t, err)

	dt := time.Date(dtOrig.Year(), dtOrig.Month(), dtOrig.Day(), dtOrig.Hour(), dtOrig.Minute(), dtOrig.Second(), dtOrig.Nanosecond(), dtOrig.Location())

	year, month, date := dt.Date()

	assert.Equal(t, 2019, year)
	assert.Equal(t, time.Month(6), month)
	assert.Equal(t, 6, date)

	hour, min, sec := dt.Clock()

	assert.Equal(t, 23, hour)
	assert.Equal(t, 32, min)
	assert.Equal(t, 53, sec)
}

func assertDate(t *testing.T, dt time.Time, expYear, expMonth, expDate, expHour, expMin, expSec int) {
	year, month, date := dt.Date()
	hour, min, sec := dt.Clock()

	assert.Equal(t, expYear, year)
	assert.Equal(t, time.Month(expMonth), month)
	assert.Equal(t, expDate, date)

	assert.Equal(t, expHour, hour)
	assert.Equal(t, expMin, min)
	assert.Equal(t, expSec, sec)
}

// Asserts that the duration is exactly what the caller sets. The parameters are passed
// in natural clock.
func assertDuration(t *testing.T, d time.Duration, expHours, expMinutes int) {
	assert.Equal(t, float64(expHours), math.Floor(d.Hours()))
	assert.Equal(t, float64(expHours*60+expMinutes), math.Floor(d.Minutes()))
}

func TestBasicTimeSlotBinning(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput(InputJSONTimestamp, "2019-06-06T12:32:53Z")
	tc.SetInput(InputIntervalType, "TimeSlot")
	tc.SetInput(InputIntervalSpan, "40m")

	done, err := act.Eval(tc)
	assert.Equal(t, true, done)
	assert.Nil(t, err)

	assert.Equal(t, "TimeSlot", tc.GetOutput(OutputType))
	assert.NotNil(t, tc.GetOutput(OutputAnchor))
	assertDate(t, tc.GetOutput(OutputAnchor).(time.Time), 2019, 06, 06, 0, 0, 0)

	// The index of the interval is 18
	assert.Equal(t, 18, tc.GetOutput(OutputIndex))

	// The offset of the interval is:
	// 40 minutes * 18 = 720 minutes, or 12 hours and 0 minutes.
	assert.NotNil(t, tc.GetOutput(OutputOffset))
	assertDuration(t, tc.GetOutput(OutputOffset).(time.Duration), 12, 0)

	// The 18th interval of 40 minutes starts exactly
	// as 12:00. This should be the interval start.
	assert.NotNil(t, tc.GetOutput(OutputStart))
	assertDate(t, tc.GetOutput(OutputStart).(time.Time), 2019, 06, 06, 12, 0, 0)

	// The span of the interval shall be 40 minutes.
	assert.NotNil(t, tc.GetOutput(OutputSpan))
	assertDuration(t, tc.GetOutput(OutputSpan).(time.Duration), 0, 40)

	// The 18th interval of 40 minutes should end exactly ast
	// 12:40.
	assert.NotNil(t, tc.GetOutput(OutputEnd))
	assertDate(t, tc.GetOutput(OutputEnd).(time.Time), 2019, 6, 6, 12, 40, 0)
}

func TestBasicDayBinning(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput(InputJSONTimestamp, "2019-06-06T23:32:53Z")
	tc.SetInput(InputIntervalType, "Day")
	done, err := act.Eval(tc)
	assert.Equal(t, done, true)
	assert.Nil(t, err, true)

	assert.Equal(t, "Day", tc.GetOutput(OutputType))
	assert.Nil(t, tc.GetOutput(OutputAnchor))
	assert.Equal(t, 0, tc.GetOutput(OutputIndex))
	assert.Nil(t, tc.GetOutput(OutputOffset))

	inStart := tc.GetOutput(OutputStart)
	assert.NotNil(t, inStart)

	pStart, startCastOk := tc.GetOutput(OutputStart).(time.Time)
	assert.True(t, startCastOk)

	year, month, date := pStart.Date()

	assert.Equal(t, 2019, year)
	assert.Equal(t, time.Month(6), month)
	assert.Equal(t, 6, date)

	hour, minute, second := pStart.Clock()
	assert.Equal(t, 0, hour)
	assert.Equal(t, 0, minute)
	assert.Equal(t, 0, second)

	inEnd := tc.GetOutput(OutputEnd)
	assert.NotNil(t, inEnd)

	pEnd, endCastOk := tc.GetOutput(OutputEnd).(time.Time)
	assert.True(t, endCastOk)

	year, month, date = pEnd.Date()

	assert.Equal(t, 2019, year)
	assert.Equal(t, time.Month(6), month)
	assert.Equal(t, 6, date)

	hour, minute, second = pEnd.Clock()
	assert.Equal(t, 23, hour)
	assert.Equal(t, 59, minute)
	assert.Equal(t, 59, second)
}
