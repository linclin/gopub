package gokits

import "strings"

//----------------------------------------------------------
func SliceStrTo(sts []string) []interface{} {
	uns := make([]interface{}, len(sts))
	for _, v := range sts {
		uns = append(uns, v)
	}
	return uns
}

//----------------------------------------------------------
func SliceInt64To(sts []int64) []interface{} {
	uns := make([]interface{}, len(sts))
	for _, v := range sts {
		uns = append(uns, v)
	}
	return uns

}

//----------------------------------------------------------
func IsSliceContainsStr(sl []string, str string) bool {
	str = strings.ToLower(str)
	for _, s := range sl {
		if strings.ToLower(s) == str {
			return true
		}
	}
	return false
}

//----------------------------------------------------------
func IsSliceContainsInt64(sl []int64, i int64) bool {
	for _, s := range sl {
		if s == i {
			return true
		}
	}
	return false
}
