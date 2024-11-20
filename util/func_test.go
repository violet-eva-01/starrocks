// package util @author: dongyuan @date  : 2024/10/11 @notes :
package util

import (
	"fmt"
	"testing"
)

func TestListSplit(t *testing.T) {

	var (
		starList  []string
		finalList []string
	)
	for i := 0; i <= 100000; i++ {
		starList = append(starList, "a")
	}

	length := len(starList)/5 + 1

	tmpMapList := ListSplit(starList, length)
	for _, i := range tmpMapList {
		finalList = append(finalList, i...)
	}
	fmt.Println(len(finalList))
}
