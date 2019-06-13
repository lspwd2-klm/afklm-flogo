package parse_ecle

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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

func TestSuccessfulParse(t *testing.T) {
	line := "api.host,185.46.213.92,0d2bdcc5-5fe8-44a0-b751-1b5a5327297a,GET,/a/b/c/d/0/products,HTTP/1.1,2,200,-,Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:67.0) Gecko/20100101 Firefox/67.0,1559780515.387_c7mes8e46wqcmbpm44t4ja3h_dd6yx8s3apv7zdeeg7n4zcvj,2019-06-06T00:21:55Z,packageKey,serviceOd,prod-j-worker-eu-west-1a-06.mashery.com,apiMethod,0,-,0.145,0.09,0.0,0.0,accessToken,1,12190,7,0.0,LE Travel API,200_OK_API,Default,planUUID,le-travelapi-shops,LE Travel API,package-UUID,serviceDefUUID"

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput(Input, line)
	done, err := act.Eval(tc)

	assert.True(t, done)
	assert.Nil(t, err)

	ecle_any := tc.GetOutput(OutputECLE)
	assert.NotNil(t, ecle_any)

	if ecle, castOk := ecle_any.(ParsedECLE); castOk {
		assert.Equal(t, "api.host", ecle.HostName)
		assert.Equal(t, "185.46.213.92", ecle.SourceIp)
		assert.Equal(t, "0d2bdcc5-5fe8-44a0-b751-1b5a5327297a", ecle.RequestUUID)
		assert.Equal(t, "GET", ecle.HttpMethod)
		assert.Equal(t, "/a/b/c/d/0/products", ecle.RequestURI)
		assert.Equal(t, "HTTP/1.1", ecle.HttpVersion)
		assert.Equal(t, 2, ecle.Bytes)
		assert.Equal(t, 200, ecle.ResponseStatusCode)
		assert.Equal(t, "-", ecle.Referrer)
		assert.Equal(t, "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:67.0) Gecko/20100101 Firefox/67.0", ecle.UserAgent)
		assert.Equal(t, "1559780515.387_c7mes8e46wqcmbpm44t4ja3h_dd6yx8s3apv7zdeeg7n4zcvj", ecle.RequestId)

		assert.NotNil(t, ecle.RequestTime)

		// Call time: 2019-06-06T00:21:55Z
		year, month, date := ecle.RequestTime.Date()
		hour, min, sec := ecle.RequestTime.Clock()

		assert.Equal(t, 2019, year)
		assert.Equal(t, time.Month(6), month)
		assert.Equal(t, 6, date)

		assert.Equal(t, 0, hour)
		assert.Equal(t, 21, min)
		assert.Equal(t, 55, sec)

		// AT this point, we are looking at:
		// packageKey,serviceOd,prod-j-worker-eu-west-1a-06.mashery.com,apiMethod,0,-,0.145,0.09,0.0,0.0,accessToken,1,12190,7,0.0,LE Travel API,200_OK_API,Default,planUUID,le-travelapi-shops,LE Travel API,package-UUID,serviceDefUUID"

		assert.Equal(t, "packageKey", ecle.PackageKey)
		assert.Equal(t, "serviceOd", ecle.ServiceId)
		assert.Equal(t, "prod-j-worker-eu-west-1a-06.mashery.com", ecle.TrafficManagerName)
		assert.Equal(t, "apiMethod", ecle.ApiMethodName)
		assert.Equal(t, false, ecle.CacheHit)
		assert.Equal(t, "-", ecle.TrafficManagerErrorCode)
		assert.Equal(t, float64(0.145), ecle.TotalRequestTime)
		assert.Equal(t, float64(0.09), ecle.RemoteTotalTime)
		assert.Equal(t, float64(0.0), ecle.ConnectTime)
		assert.Equal(t, float64(0.0), ecle.PreTransferTime)
		assert.Equal(t, "accessToken", ecle.AccessToken)
		assert.Equal(t, true, ecle.SslEnabled)
		assert.Equal(t, 12190, ecle.ActualQuota)
		assert.Equal(t, 7, ecle.ActualQPS)
		assert.Equal(t, float64(0.0), ecle.ClientTransferTime)
		assert.Equal(t, float64(0.0), ecle.ClientTransferTime)

		// LE Travel API,200_OK_API,Default,planUUID,le-travelapi-shops,LE Travel API,package-UUID,serviceDefUUID"
		assert.Equal(t, "LE Travel API", ecle.PackageName)
		assert.Equal(t, "200_OK_API", ecle.ResponseMessage)
		assert.Equal(t, "Default", ecle.PlanName)
		assert.Equal(t, "planUUID", ecle.PlanUUID)
		assert.Equal(t, "le-travelapi-shops", ecle.EndpointName)
		assert.Equal(t, "LE Travel API", ecle.PackageName)
		assert.Equal(t, "package-UUID", ecle.PackageUUID)
		assert.Equal(t, "serviceDefUUID", ecle.ServiceDefinitionEndpointUUID)

	} else {
		t.Error("ECLE is not expected type")
		t.Fail()
	}

	// Checking the outputs for the individual fields
	assert.Equal(t, "api.host", tc.GetOutput(OutputHostName))
	assert.Equal(t, "185.46.213.92", tc.GetOutput(OutputSourceIp))
	assert.Equal(t, "0d2bdcc5-5fe8-44a0-b751-1b5a5327297a", tc.GetOutput(OutputRequestUUID))
	assert.Equal(t, "GET", tc.GetOutput(OutputHttpMethod))
	assert.Equal(t, "/a/b/c/d/0/products", tc.GetOutput(OutputRequestURI))
	assert.Equal(t, "HTTP/1.1", tc.GetOutput(OutputHttpVersion))
	assert.Equal(t, 2, tc.GetOutput(OutputBytes))
	assert.Equal(t, 200, tc.GetOutput(OutputResponseStatusCode))
	assert.Equal(t, "-", tc.GetOutput(OutputReferrer))
	assert.Equal(t, "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:67.0) Gecko/20100101 Firefox/67.0", tc.GetOutput(OutputUserAgent))
	assert.Equal(t, "1559780515.387_c7mes8e46wqcmbpm44t4ja3h_dd6yx8s3apv7zdeeg7n4zcvj", tc.GetOutput(OutputRequestId))

	callTime_any := tc.GetOutput(OutputRequestTime)
	assert.NotNil(t, callTime_any)

	if callTime, ctCastOk := callTime_any.(time.Time); ctCastOk {
		// Call time: 2019-06-06T00:21:55Z
		year, month, date := callTime.Date()
		hour, min, sec := callTime.Clock()

		assert.Equal(t, 2019, year)
		assert.Equal(t, time.Month(6), month)
		assert.Equal(t, 6, date)

		assert.Equal(t, 0, hour)
		assert.Equal(t, 21, min)
		assert.Equal(t, 55, sec)
	}

	// AT this point, we are looking at:
	// packageKey,serviceOd,prod-j-worker-eu-west-1a-06.mashery.com,apiMethod,0,-,0.145,0.09,0.0,0.0,accessToken,1,12190,7,0.0,LE Travel API,200_OK_API,Default,planUUID,le-travelapi-shops,LE Travel API,package-UUID,serviceDefUUID"

	assert.Equal(t, "packageKey", tc.GetOutput(OutputPackageKey))
	assert.Equal(t, "serviceOd", tc.GetOutput(OutputServiceId))
	assert.Equal(t, "prod-j-worker-eu-west-1a-06.mashery.com", tc.GetOutput(OutputTrafficManagerName))
	assert.Equal(t, "apiMethod", tc.GetOutput(OutputAPIMethodName))
	assert.Equal(t, false, tc.GetOutput(OutputCacheHit))
	assert.Equal(t, "-", tc.GetOutput(OutputTrafficManagerErrorCode))
	assert.Equal(t, float64(0.145), tc.GetOutput(OutputTotalRequestTime))
	assert.Equal(t, float64(0.09), tc.GetOutput(OutputRemoteTotalTime))
	assert.Equal(t, float64(0.0), tc.GetOutput(OutputConnectTime))
	assert.Equal(t, float64(0.0), tc.GetOutput(OutputPreTransferTime))
	assert.Equal(t, "accessToken", tc.GetOutput(OutputAccessToken))
	assert.Equal(t, true, tc.GetOutput(OutputSslEnabled))
	assert.Equal(t, 12190, tc.GetOutput(OutputActualQuota))
	assert.Equal(t, 7, tc.GetOutput(OutputActualQPS))
	assert.Equal(t, float64(0.0), tc.GetOutput(OutputClientTransferTime))

	// LE Travel API,200_OK_API,Default,planUUID,le-travelapi-shops,LE Travel API,package-UUID,serviceDefUUID"
	assert.Equal(t, "LE Travel API", tc.GetOutput(OutputPackageName))
	assert.Equal(t, "200_OK_API", tc.GetOutput(OutputResponseMessage))
	assert.Equal(t, "Default", tc.GetOutput(OutputPlanName))
	assert.Equal(t, "planUUID", tc.GetOutput(OutputPlanUUID))
	assert.Equal(t, "le-travelapi-shops", tc.GetOutput(OutputEndpointName))
	assert.Equal(t, "LE Travel API", tc.GetOutput(OutputPackageName))
	assert.Equal(t, "package-UUID", tc.GetOutput(OutputPackageUUID))
	assert.Equal(t, "serviceDefUUID", tc.GetOutput(OutputServiceDefinitionEndpointUUID))
}
