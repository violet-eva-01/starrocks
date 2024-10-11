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

func Match(str string, strArray []string) bool {
	for _, i := range strArray {
		if strings.Contains(str, i) {
			return true
		}
	}
	return false
}

func RemoveMatchElement(list1, list2 []string) []string {
	result := make([]string, 0)
	for _, i := range list1 {
		if !Match(i, list2) {
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

func ListSplit(input []string, length int) map[int][]string {

	times := len(input) / length // 10001 / 2001 = 4
	output := make(map[int][]string, times+1)
	residual := len(input) % length // 10001 % 2001 = 1997

	if times == 0 || (times == 1 && residual == 0) {
		output[0] = input
	} else {
		if residual == 0 {
			times -= 1
		}

		starLen := 0
		endLen := length
		for index := 0; index <= times; index++ {
			output[index] = input[starLen:endLen]
			starLen += length
			if residual != 0 && index == times-1 {
				endLen += residual
			} else {
				endLen += length
			}
		}
	}

	return output
}
