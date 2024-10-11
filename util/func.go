package util

import (
	"sort"
	"strings"
)

func In(str string, strArray []string, isSort bool) bool {
	if isSort {
		sort.Strings(strArray)
	}
	index := sort.SearchStrings(strArray, str)
	if index < len(strArray) && strArray[index] == str {
		return true
	}
	return false
}

func RemoveCoincideElement(list1, list2 []string, isSort bool) []string {
	result := make([]string, 0)
	for _, i := range list1 {
		if !In(i, list2, isSort) {
			result = append(result, i)
		}
	}
	return result
}

func IfExists(str string, strArray []string) bool {
	for _, i := range strArray {
		if strings.Contains(str, i) {
			return true
		}
	}
	return false
}

func RemoveIfExistsElement(list1, list2 []string) []string {
	result := make([]string, 0)
	for _, i := range list1 {
		if !IfExists(i, list2) {
			result = append(result, i)
		}
	}
	return result
}

func RemoveRepeatElement(list []string) []string {
	temp := make(map[string]struct{})
	index := 0
	for _, v := range list {
		v = strings.TrimSpace(v)
		temp[v] = struct{}{}
	}
	tempList := make([]string, len(temp))
	for key := range temp {
		tempList[index] = key
		index++
	}
	return tempList
}

func RemoveRepeatElementAndToLower(list []string) []string {
	temp := make(map[string]struct{})
	index := 0
	for _, v := range list {
		v = strings.ToLower(strings.TrimSpace(v))
		temp[v] = struct{}{}
	}
	tempList := make([]string, len(temp))
	for key := range temp {
		tempList[index] = key
		index++
	}
	return tempList
}

func ListReplace(list []string, old, new string) []string {
	temp := make(map[string]struct{})
	index := 0
	for _, v := range list {
		v = strings.ReplaceAll(v, old, new)
		temp[v] = struct{}{}
	}
	tempList := make([]string, len(temp))
	for key := range temp {
		tempList[index] = key
		index++
	}
	return tempList
}
