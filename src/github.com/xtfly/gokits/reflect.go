package gokits

import (
	"reflect"
	"strconv"
)

type RefVal reflect.Value

func (r RefVal) ToString() string {
	v := reflect.Value(r)
	kind := v.Kind()
	if kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 {
		return strconv.FormatInt(v.Int(), 10)
	} else if kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 {
		return strconv.FormatUint(v.Uint(), 10)
	} else if kind == reflect.Float32 || kind == reflect.Float64 {
		return strconv.FormatFloat(v.Float(), 'f', 2, 64)
	}
	return ""
}

func (r RefVal) IsEmpty() bool {
	v := reflect.Value(r)
	k := v.Kind()
	if k == reflect.Bool {
		return v.Bool() == false
	} else if reflect.Int < k && k < reflect.Int64 {
		return v.Int() == 0
	} else if reflect.Uint < k && k < reflect.Uintptr {
		return v.Uint() == 0
	} else if k == reflect.Float32 || k == reflect.Float64 {
		return v.Float() == 0
	} else if k == reflect.Array || k == reflect.Map || k == reflect.Slice || k == reflect.String {
		return v.Len() == 0
	} else if k == reflect.Interface || k == reflect.Ptr {
		return v.IsNil()
	}
	return false
}
