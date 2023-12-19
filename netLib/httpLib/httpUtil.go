package httpLib

import (
	"fmt"
	"sort"
	"strings"
)

func GetFormHeader() map[string]string {
	return map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
}

func AddFormHeader(header map[string]string) {
	header["Content-Type"] = "application/x-www-form-urlencoded"
}

func GetJsonHeader() map[string]string {
	return map[string]string{"Content-Type": "application/json;charset=UTF-8"}
}

func AddJsonHeader(header map[string]string) {
	header["Content-Type"] = "application/json;charset=UTF-8"
}

func AssembleRequestString(paramMap map[string]string) string {
	if paramMap == nil {
		return ""
	}

	itemList := make([]string, len(paramMap))
	var index int
	for key, value := range paramMap {
		itemList[index] = fmt.Sprintf("%s=%s", key, value)
		index++
	}

	return strings.Join(itemList, "&")
}

func AssembleSortedRequestString(paramMap map[string]string, ascending bool) string {
	if paramMap == nil {
		return ""
	}

	sortedKeyList := make([]string, len(paramMap))
	var index int
	for key := range paramMap {
		sortedKeyList[index] = key
		index++
	}

	sort.Strings(sortedKeyList)
	if !ascending {
		for left, right := 0, len(sortedKeyList)-1; left < right; left, right = left+1, right-1 {
			sortedKeyList[left], sortedKeyList[right] = sortedKeyList[right], sortedKeyList[left]
		}
	}

	itemList := make([]string, len(paramMap))
	for index, key := range sortedKeyList {
		itemList[index] = fmt.Sprintf("%s=%s", key, paramMap[key])
	}

	return strings.Join(itemList, "&")
}
