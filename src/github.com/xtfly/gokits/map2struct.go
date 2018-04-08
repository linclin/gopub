package gokits

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const defaultTag = "map"

func Map2Struct(vals map[string]interface{}, dst interface{}) (err error) {
	return Map2StructTag(vals, dst, defaultTag)
}

func Map2StructTag(vals map[string]interface{}, dst interface{}, tagName string) (err error) {
	defer func() {
		e := recover()
		if e != nil {
			if v, ok := e.(error); ok {
				err = fmt.Errorf("Panic: %v", v.Error())
			} else {
				err = fmt.Errorf("Panic: %v", e)
			}
		}
	}()

	pt := reflect.TypeOf(dst)
	pv := reflect.ValueOf(dst)

	if pv.Kind() != reflect.Ptr || pv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("not a pointer of struct")
	}

	var f reflect.StructField
	var fv reflect.Value

	for i := 0; i < pt.Elem().NumField(); i++ {
		f = pt.Elem().Field(i)
		fv = pv.Elem().Field(i)

		if f.Anonymous {
			continue
		}

		if !fv.CanSet() {
			continue
		}

		tag := f.Tag.Get(tagName)
		name, option := parseTag(tag)

		if name == "" {
			// tag name is not set, use field name
			name = f.Name
		}

		err = map2Field(vals, fv, name, option)
		if err != nil {
			return fmt.Errorf("field %v(%v) error: %v", name, fv.Type().Kind(), err.Error())
		}

		continue
	}

	return nil
}

func Map2Field(vals map[string]interface{}, dst interface{}, tag string) error {
	fv := reflect.ValueOf(dst)
	if fv.Kind() != reflect.Ptr {
		return fmt.Errorf("dst must be a pointer")
	}

	name, option := parseTag(tag)

	return map2Field(vals, fv, name, option)
}

func map2Field(vals map[string]interface{}, fv reflect.Value, name, option string) error {
	if name == "-" || name == "" {
		return nil // ignore "-"
	}

	// value from map
	val, ok := vals[name]
	if !ok {
		val, ok = vals[strings.ToLower(name)]
	}

	if !ok { // value not found
		if option == "required" {
			return fmt.Errorf("'%v' is required", name)
		}

		if option == "string" {
			return nil // ignore it when 'option' is string
		}

		if len(option) != 0 {
			val = option // 'option' means 'default value' here
		} else {
			return nil // ignore it
		}
	}

	return convert(val, fv)
}

// Convert varies types of value to a certain type.
// Value can be type of string or json or whatever type which is convertable to the target type.
func Convert(dst interface{}, val interface{}) error {
	fv := reflect.ValueOf(dst)
	if fv.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("dst must be a pointer")
	}
	return convert(val, fv)
}

func convert(val interface{}, fv reflect.Value) (err error) {
	// assign or convert value to field
	if assignToField(val, fv) == nil {
		return nil
	}

	switch v := val.(type) {
	case string:
		// parse string to value
		s := strings.TrimSpace(v)
		err = convertStringToValue(s, fv, fv.Type().Kind())

	case json.RawMessage:
		// unmarshal json
		err = convertJsonToValue(v, fv)

	default:
		err = fmt.Errorf("value type is not supported: value=%v", val)
	}

	return err
}

func assignToField(val interface{}, fv reflect.Value) error {
	vv := reflect.ValueOf(val)
	vt := reflect.TypeOf(val)
	ft := fv.Type()

	// assign or convert value to field
	if vt.AssignableTo(ft) {
		fv.Set(vv)
		return nil
	}
	if vt.ConvertibleTo(ft) {
		fv.Set(vv.Convert(ft))
		return nil
	}
	return fmt.Errorf("can not assign: value=%v(%v)", val, vt.Kind())
}

func convertJsonToValue(data json.RawMessage, fv reflect.Value) error {
	var err error

	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			fv.Set(reflect.New(fv.Type().Elem()))
		}
	} else {
		fv = fv.Addr()
	}

	err = json.Unmarshal(data, fv.Interface())

	if err != nil {
		return fmt.Errorf("invalid json: %v, %v", err.Error(), string(data))
	}

	return nil
}

func convertStringToValue(s string, fv reflect.Value, kind reflect.Kind) error {
	if !fv.CanAddr() {
		return fmt.Errorf("target can not addr")
	}

	if assignToField(s, fv) == nil {
		return nil
	}

	if kind == reflect.String {
		fv.SetString(s)
		return nil
	}

	if kind == reflect.Slice {
		return convertStringToSlice(s, fv)
	}

	if kind == reflect.Ptr || kind == reflect.Struct {
		return convertJsonToValue(json.RawMessage(s), fv)
	}

	if kind == reflect.Bool {
		switch strings.ToLower(s) {
		case "true":
			fv.SetBool(true)
		case "false":
			fv.SetBool(false)
		case "1":
			fv.SetBool(true)
		case "0":
			fv.SetBool(false)
		case "y":
			fv.SetBool(true)
		case "n":
			fv.SetBool(false)
		default:
			return fmt.Errorf("invalid bool: value=%v", s)
		}
		return nil
	}

	if reflect.Int <= kind && kind <= reflect.Int64 {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int: value=%v, err=%s", s, err.Error())
		}
		fv.SetInt(i)

	} else if reflect.Uint <= kind && kind <= reflect.Uint64 {
		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid uint: value=%v, err=%s", s, err.Error())
		}
		fv.SetUint(i)

	} else if reflect.Float32 == kind || kind == reflect.Float64 {
		i, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("invalid float: value=%v, err=%s", s, err.Error())
		}
		fv.SetFloat(i)

	} else {
		// type not support
		return fmt.Errorf("type not support: value=%v", s)
	}
	return nil
}

func convertStringToSlice(s string, fv reflect.Value) error {
	var err error
	ft := fv.Type()
	et := ft.Elem()

	if len(s) == 0 {
		return nil
	}

	data := json.RawMessage(s)
	if data[0] == '[' && data[len(data)-1] == ']' {
		return convertJsonToValue(data, fv)
	}

	ss := strings.Split(s, ",")
	fs := reflect.MakeSlice(ft, 0, len(ss))

	for _, si := range ss {
		ev := reflect.New(et).Elem()

		err = convertStringToValue(si, ev, et.Kind())
		if err != nil {
			return err
		}
		fs = reflect.Append(fs, ev)
	}

	fv.Set(fs)

	return nil
}

func parseTag(tag string) (string, string) {
	tags := strings.Split(tag, ",")

	if len(tags) <= 0 {
		return "", ""
	}

	if len(tags) == 1 {
		return tags[0], ""
	}

	return tags[0], tags[1]
}
