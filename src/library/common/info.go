package common

import (
	"encoding/json"
	_ "fmt"
	_ "reflect"
	"strconv"
	"strings"
)

type Info map[string]interface{}

func (info Info) String(key string) string {
	if info[key] == nil {
		return ""
	}
	// 其他类型取值时转换类型，这里只取值
	return ToString(info[key])
}

func (info Info) Bool(key string) bool {
	if info[key] == nil {
		return false
	}
	info[key] = ToBool(info[key])
	return info[key].(bool)
}

func (info Info) Int(key string) int {
	if info[key] == nil {
		return 0
	}
	info[key] = ToInt(info[key])
	return info[key].(int)
}

func (info Info) Float(key string) float64 {
	if info[key] == nil {
		return float64(0)
	}
	info[key] = ToFloat(info[key])
	return info[key].(float64)
}

func (info Info) ToString() string {
	if info == nil {
		return ""
	}
	ret, _ := json.Marshal(info)
	return string(ret)
}

/**
 * 支持类型：[]Info,[]interface{}
 * 其他类型报错
 */
func (info Info) InfoList(key string) []Info {
	if info[key] == nil {
		info[key] = []Info{}
		return info[key].([]Info)
	}

	switch info[key].(type) {
	case []Info:
		return info[key].([]Info)
	case []interface{}:
		// 这个是可以转换的结构
		ret := []Info{}
		for _, item := range info[key].([]interface{}) {
			ret = append(ret, Info(item.(map[string]interface{})))
		}
		info[key] = ret
		return info[key].([]Info)
	default:
		// 不能识别的结构，直接报错
		panic("InfoList出错")
	}
}

/**
 * 判断指定key是否为list
 */
func (info Info) IsList(key string) bool {
	if info[key] == nil {
		return false
	}

	switch info[key].(type) {
	case []Info:
		return true
	case []interface{}:
		return true
	default:
		return false
	}
}

/**
 * 支持类型：Info,interface{}
 * 其他类型报错
 */
func (info Info) Info(key string) Info {
	if info[key] == nil {
		return nil
	}

	switch info[key].(type) {
	case Info:
		return info[key].(Info)
	case interface{}:
		// 这个是可以转换的结构
		info[key] = Info(info[key].(map[string]interface{}))
		return info[key].(Info)
	default:
		// 不能识别的结构，直接报错
		panic("InfoList出错")
	}
}

/**
 * 使用info中的值覆盖source
 */
func (info Info) Merge(source Info) {
	for key, value := range source {
		if info[key] == nil {
			info[key] = value
		}
	}
}

/**
 * 对象转换为string
 * 支持类型：int,float64,string,bool(true:"1";false:"0")
 * 其他类型报错
 */
func ToString(obj interface{}) string {
	switch obj.(type) {
	case int:
		return strconv.Itoa(obj.(int))
	case float64:
		return strconv.FormatFloat(obj.(float64), 'f', -1, 64)
	case string:
		return obj.(string)
	case bool:
		if obj.(bool) {
			return "1"
		} else {
			return "0"
		}
	default:
		panic("ToString出错")
	}
}

/**
 * 对象转换为bool
 * 支持类型：int,float64,string,bool
 * 其他类型报错
 */
func ToBool(obj interface{}) bool {
	switch obj.(type) {
	case int:
		if obj.(int) == 0 {
			return false
		} else {
			return true
		}
	case float64:
		if obj.(float64) == 0 {
			return false
		} else {
			return true
		}
	case string:
		trues := map[string]int{"true": 1, "是": 1, "1": 1, "真": 1}
		if _, ok := trues[strings.ToLower(obj.(string))]; ok {
			return true
		} else {
			return false
		}
	case bool:
		return obj.(bool)
	default:
		panic("ToBool出错")
	}
}

/**
 * 对象转换为int
 * 支持类型：int,float64,string,bool(true:1;false:0)
 * 其他类型报错
 */
func ToInt(obj interface{}) int {
	switch obj.(type) {
	case int:
		return obj.(int)
	case float64:
		return int(obj.(float64))
	case string:
		ret, _ := strconv.Atoi(obj.(string))
		return ret
	case bool:
		if obj.(bool) {
			return 1
		} else {
			return 0
		}
	default:
		panic("ToInt出错")
	}
}

/**
 * 对象转换为float64
 * 支持类型：int,float64,string,bool(true:1;false:0)
 * 其他类型报错
 */
func ToFloat(obj interface{}) float64 {
	switch obj.(type) {
	case int:
		return float64(obj.(int))
	case float64:
		return obj.(float64)
	case string:
		ret, _ := strconv.ParseFloat(obj.(string), 64)
		return ret
	case bool:
		if obj.(bool) {
			return float64(1)
		} else {
			return float64(0)
		}
	default:
		panic("ToFloat出错")
	}
}
