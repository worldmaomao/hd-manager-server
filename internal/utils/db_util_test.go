package utils

import (
	"fmt"
	"regexp"
	"testing"
)

func TestFilterSpecChar(t *testing.T) {
	str1 := "a\\sd*23-23*4}zzsd*()&^&%f{"

	reg, _ := regexp.Compile(`([\*\.\?\+\$\^\[\]\(\)\{\}\|\\\/]+)`)
	result := reg.ReplaceAllString(str1, `\${1}`)

	fmt.Println(result)
}
