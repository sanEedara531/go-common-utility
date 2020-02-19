package common

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

type utilityService struct {
}

func getUUID() string {
	uuid1, err := uuid.NewUUID()
	if err != nil {
		fmt.Println(err)
	}
	return uuid1.String()
}

func getUTCTimeString() string {
	return time.Now().UTC().String()
}

func prepareUpdateExpression(request interface{}, keysToExclude []string) map[string]interface{} {
	var UpdateExpression string
	response := make(map[string]interface{})
	expressionAttributeValues := make(map[string]interface{})
	v := reflect.ValueOf(request)

	typeOfS := v.Type()
	numOfFields := v.NumField()

	for i := 0; i < numOfFields; i++ {
		if !contains(keysToExclude, typeOfS.Field(i).Name) {
			var keyName = ":" + typeOfS.Field(i).Name
			var data = v.Field(i).Interface()
			expressionAttributeValues[keyName] = data
			UpdateExpression = UpdateExpression + lowerFirst(typeOfS.Field(i).Name) + "= :" + typeOfS.Field(i).Name

			if i != numOfFields-1 {
				UpdateExpression = UpdateExpression + ", "
			}
		}
	}
	UpdateExpression = "set " + UpdateExpression
	response["updateExpression"] = UpdateExpression
	response["expressionAttributeValues"] = expressionAttributeValues
	return response
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

//InterfaceToMap function converts an interface to map
func InterfaceToMap(request interface{}) map[string]interface{} {

	response := make(map[string]interface{})
	v := reflect.ValueOf(request)
	typeOfS := v.Type()

	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)
			response[key.String()] = strct.Interface()

		}
	} else {
		for i := 0; i < v.NumField(); i++ {
			var keyName = typeOfS.Field(i).Name
			var data = v.Field(i).Interface()
			response[keyName] = data
		}
	}
	return response
}

func removeEmptyFields(vehicleCreated VmsVehicle) []string {
	var result []string
	v := reflect.ValueOf(vehicleCreated)

	typeOfS := v.Type()
	numOfFields := v.NumField()
	for i := 0; i < numOfFields; i++ {
		if len(strings.TrimSpace(v.Field(i).String())) == 0 {
			result = append(result, typeOfS.Field(i).Name)
		}
	}
	return result
}

func getCall(url string, query_params map[string]string, header http.Header) ([]byte, error){
	var endpointFms = appendQueryParams(url, query_params)
	httpClient := http.Client{
		Timeout: time.Second * ConfigurationObj.Server.TimeOut,
	}

	req, err := http.NewRequest(http.MethodGet, endpointFms, nil)
	if err != nil {
		ZapLoggerObj.Error("Unable to make request to endpoint " + endpointFms)
		return nil, err
	}
	req.Header = header

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		ZapLoggerObj.Error("get request to endpoint failed " + endpointFms)
		return nil, getErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		ZapLoggerObj.Error("unable to convert to an object " + endpointFms)
		return nil, readErr
	}
	return body, readErr
}

func appendQueryParams(url string, query_params map[string]string) string {
	var result strings.Builder
	result.WriteString(url + "?")

	for query_param_key, query_param_value := range query_params {
		if len(strings.TrimSpace(query_param_value)) > 0 {
			result.WriteString(query_param_key + "=" + query_param_value + "&")
		}
	}
	var resultString = result.String()
	return resultString[:len(resultString)-1]
}

func isEmpty(str string) bool {
	if len(str) > 0 && len(strings.TrimSpace(str)) > 0 {
		return false
	}
	return true
}