package parse_ecle

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"strconv"
	"strings"
	"time"
)

var log = logger.GetLogger("parse-ecle-action")

const (
	Input                               = "input"
	OutputECLE                          = "ecle"
	OutputHostName                      = "HostName"
	OutputSourceIp                      = "SourceIp"
	OutputRequestUUID                   = "RequestUUID"
	OutputHttpMethod                    = "HttpMethod"
	OutputRequestURI                    = "RequestURI"
	OutputHttpVersion                   = "HttpVersion"
	OutputBytes                         = "Bytes"
	OutputResponseStatusCode            = "ResponseStatusCode"
	OutputReferrer                      = "Referrer"
	OutputUserAgent                     = "UserAgent"
	OutputRequestId                     = "RequestId"
	OutputRequestTime                   = "RequestTime"
	OutputPackageKey                    = "PackageKey"
	OutputServiceId                     = "ServiceId"
	OutputTrafficManagerName            = "TrafficManagerName"
	OutputAPIMethodName                 = "APIMethodName"
	OutputCacheHit                      = "CacheHit"
	OutputTrafficManagerErrorCode       = "TrafficManagerErrorCode"
	OutputTotalRequestTime              = "TotalRequestTime"
	OutputRemoteTotalTime               = "RemoteTotalTime"
	OutputConnectTime                   = "ConnectTime"
	OutputPreTransferTime               = "PreTransferTime"
	OutputAccessToken                   = "AccessToken"
	OutputSslEnabled                    = "SslEnabled"
	OutputActualQuota                   = "ActualQuota"
	OutputActualQPS                     = "ActualQPS"
	OutputClientTransferTime            = "ClientTransferTime"
	OutputServiceName                   = "ServiceName"
	OutputResponseMessage               = "ResponseMessage"
	OutputPlanName                      = "PlanName"
	OutputPlanUUID                      = "PlanUUID"
	OutputEndpointName                  = "EndpointName"
	OutputPackageName                   = "PackageName"
	OutputPackageUUID                   = "PackageUUID"
	OutputServiceDefinitionEndpointUUID = "ServiceDefinitionEndpointUUID"
)

type ParsedECLE struct {
	HostName                      string
	SourceIp                      string
	RequestUUID                   string
	HttpMethod                    string
	RequestURI                    string
	HttpVersion                   string
	Bytes                         int
	ResponseStatusCode            int
	Referrer                      string
	UserAgent                     string
	RequestId                     string
	RequestTime                   time.Time
	PackageKey                    string
	ServiceId                     string
	TrafficManagerName            string
	ApiMethodName                 string
	CacheHit                      bool
	TrafficManagerErrorCode       string
	TotalRequestTime              float64
	RemoteTotalTime               float64
	ConnectTime                   float64
	PreTransferTime               float64
	AccessToken                   string
	SslEnabled                    bool
	ActualQuota                   int
	ActualQPS                     int
	ClientTransferTime            float64
	ServiceName                   string
	ResponseMessage               string
	PlanName                      string
	PlanUUID                      string
	EndpointName                  string
	PackageName                   string
	PackageUUID                   string
	ServiceDefinitionEndpointUUID string
}

type AmtParseECLE struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AmtParseECLE{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *AmtParseECLE) Metadata() *activity.Metadata {
	return a.metadata
}

func writeOutput(ecle ParsedECLE, ctx activity.Context) {
	ctx.SetOutput(OutputECLE, ecle)

	ctx.SetOutput(OutputHostName, ecle.HostName)
	ctx.SetOutput(OutputSourceIp, ecle.SourceIp)
	ctx.SetOutput(OutputRequestUUID, ecle.RequestUUID)
	ctx.SetOutput(OutputHttpMethod, ecle.HttpMethod)
	ctx.SetOutput(OutputRequestURI, ecle.RequestURI)
	ctx.SetOutput(OutputHttpVersion, ecle.HttpVersion)
	ctx.SetOutput(OutputBytes, ecle.Bytes)
	ctx.SetOutput(OutputResponseStatusCode, ecle.ResponseStatusCode)
	ctx.SetOutput(OutputReferrer, ecle.Referrer)
	ctx.SetOutput(OutputUserAgent, ecle.UserAgent)
	ctx.SetOutput(OutputRequestId, ecle.RequestId)
	ctx.SetOutput(OutputRequestTime, ecle.RequestTime)
	ctx.SetOutput(OutputPackageKey, ecle.PackageKey)
	ctx.SetOutput(OutputServiceId, ecle.ServiceId)
	ctx.SetOutput(OutputTrafficManagerName, ecle.TrafficManagerName)
	ctx.SetOutput(OutputAPIMethodName, ecle.ApiMethodName)
	ctx.SetOutput(OutputCacheHit, ecle.CacheHit)
	ctx.SetOutput(OutputTrafficManagerErrorCode, ecle.TrafficManagerErrorCode)
	ctx.SetOutput(OutputTotalRequestTime, ecle.TotalRequestTime)
	ctx.SetOutput(OutputRemoteTotalTime, ecle.RemoteTotalTime)
	ctx.SetOutput(OutputConnectTime, ecle.ConnectTime)
	ctx.SetOutput(OutputPreTransferTime, ecle.PreTransferTime)
	ctx.SetOutput(OutputAccessToken, ecle.AccessToken)
	ctx.SetOutput(OutputSslEnabled, ecle.SslEnabled)
	ctx.SetOutput(OutputActualQuota, ecle.ActualQuota)
	ctx.SetOutput(OutputActualQPS, ecle.ActualQPS)
	ctx.SetOutput(OutputClientTransferTime, ecle.ClientTransferTime)
	ctx.SetOutput(OutputServiceName, ecle.ServiceName)
	ctx.SetOutput(OutputResponseMessage, ecle.ResponseMessage)
	ctx.SetOutput(OutputPlanName, ecle.PlanName)
	ctx.SetOutput(OutputPlanUUID, ecle.PlanUUID)
	ctx.SetOutput(OutputEndpointName, ecle.EndpointName)
	ctx.SetOutput(OutputPackageName, ecle.PackageName)
	ctx.SetOutput(OutputPackageUUID, ecle.PackageUUID)
	ctx.SetOutput(OutputServiceDefinitionEndpointUUID, ecle.ServiceDefinitionEndpointUUID)
}

func to_i(dest *int, source string, errCnt *int) {
	var perr error

	*dest, perr = strconv.Atoi(source)
	if perr != nil {
		log.Error(fmt.Sprintf("Can't parse %s as int", source))
		*errCnt += 1
	}
}

func to_f(dest *float64, source string, errCnt *int) {
	var perr error

	*dest, perr = strconv.ParseFloat(source, 64)
	if perr != nil {
		log.Error(fmt.Sprintf("Can't parse %s as float", source))
		*errCnt += 1
	}
}

func to_time(dest *time.Time, source string, errCnt *int) {
	var perr error

	*dest, perr = time.Parse(time.RFC3339, source)
	if perr != nil {
		log.Error(fmt.Sprintf("Can't parse %s as time", source))
		*errCnt += 1
	}
}

func (a *AmtParseECLE) Eval(ctx activity.Context) (done bool, err error) {
	input := ctx.GetInput(Input).(string)

	r := csv.NewReader(strings.NewReader(input))
	rec, err := r.Read()

	if err == nil {
		retVal := ParsedECLE{}

		errCnt := 0

		if len(rec) == 35 {
			retVal.HostName = rec[0]
			retVal.SourceIp = rec[1]
			retVal.RequestUUID = rec[2]
			retVal.HttpMethod = rec[3]
			retVal.RequestURI = rec[4]
			retVal.HttpVersion = rec[5]

			to_i(&retVal.Bytes, rec[6], &errCnt)
			to_i(&retVal.ResponseStatusCode, rec[7], &errCnt)

			retVal.Referrer = rec[8]
			retVal.UserAgent = rec[9]
			retVal.RequestId = rec[10]
			retVal.RequestId = rec[10]

			to_time(&retVal.RequestTime, rec[11], &errCnt)

			retVal.PackageKey = rec[12]
			retVal.ServiceId = rec[13]
			retVal.TrafficManagerName = rec[14]
			retVal.ApiMethodName = rec[15]
			retVal.CacheHit = rec[16] == "true"
			retVal.TrafficManagerErrorCode = rec[17]

			to_f(&retVal.TotalRequestTime, rec[18], &errCnt)
			to_f(&retVal.RemoteTotalTime, rec[19], &errCnt)
			to_f(&retVal.ConnectTime, rec[20], &errCnt)
			to_f(&retVal.PreTransferTime, rec[21], &errCnt)

			retVal.AccessToken = rec[22]
			retVal.SslEnabled = rec[23] == "1"

			to_i(&retVal.ActualQuota, rec[24], &errCnt)
			to_i(&retVal.ActualQPS, rec[25], &errCnt)

			to_f(&retVal.ClientTransferTime, rec[26], &errCnt)

			retVal.ServiceName = rec[27]
			retVal.ResponseMessage = rec[28]
			retVal.PlanName = rec[29]
			retVal.PlanUUID = rec[30]
			retVal.EndpointName = rec[31]
			retVal.PackageName = rec[32]
			retVal.PackageUUID = rec[33]
			retVal.ServiceDefinitionEndpointUUID = rec[34]

			if errCnt == 0 {
				writeOutput(retVal, ctx)
				return true, nil
			} else {
				log.Error(fmt.Sprintf("There were %d errors while trying to parse ECLE string: >|%s|<", errCnt, input))
				return false, errors.New("field conversion errors")
			}
		} else {
			return false, errors.New("wrong number of ELCE fields")
		}

	} else {
		return false, err
	}
}
