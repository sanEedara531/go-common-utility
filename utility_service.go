package common

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

type utilityService struct {
}

func GetUUID() string {
	uuid1, err := uuid.NewUUID()
	if err != nil {
		fmt.Println(err)
	}
	return uuid1.String()
}

func GetUTCTimeString() string {
	return time.Now().UTC().String()
}

func PrepareUpdateExpression(request interface{}, keysToExclude []string) map[string]interface{} {
	var UpdateExpression string
	response := make(map[string]interface{})
	expressionAttributeValues := make(map[string]interface{})
	v := reflect.ValueOf(request)

	typeOfS := v.Type()
	numOfFields := v.NumField()

	for i := 0; i < numOfFields; i++ {
		if !Contains(keysToExclude, typeOfS.Field(i).Name) {
			var keyName = ":" + typeOfS.Field(i).Name
			var data = v.Field(i).Interface()
			expressionAttributeValues[keyName] = data
			UpdateExpression = UpdateExpression + LowerFirst(typeOfS.Field(i).Name) + "= :" + typeOfS.Field(i).Name

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

func LowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

func Contains(arr []string, str string) bool {
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


func AppendQueryParams(url string, query_params map[string]string) string {
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

func IsEmpty(str string) bool {
	if len(str) > 0 && len(strings.TrimSpace(str)) > 0 {
		return false
	}
	return true
}