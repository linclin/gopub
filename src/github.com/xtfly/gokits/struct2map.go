package gokits

import "reflect"

func Struct2Map(s interface{}) map[string]interface{} {
	return Struct2MapTag(s, defaultTag)
}

func Struct2MapTag(s interface{}, tagName string) map[string]interface{} {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
		t = t.Elem()
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	m := make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		if !fv.CanInterface() {
			continue
		}

		if ft.PkgPath != "" { // unexported
			continue
		}

		var name string
		var option string

		name, option = parseTag(ft.Tag.Get(tagName))

		if name == "-" {
			continue // ignore "-"
		}

		if name == "" {
			name = ft.Name // use field name
		}

		if option == "omitempty" {
			if RefVal(fv).IsEmpty() {
				continue // skip empty field
			}
		}

		// ft.Anonymous means embedded field
		if ft.Anonymous {
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				continue // nil
			}

			if (fv.Kind() == reflect.Struct) ||
				(fv.Kind() == reflect.Ptr && fv.Elem().Kind() == reflect.Struct) {

				// embedded struct
				embedded := Struct2MapTag(fv.Interface(), tagName)

				for embName, embValue := range embedded {
					m[embName] = embValue
				}
			}
			continue
		}

		if option == "string" {
			s := RefVal(fv).ToString()
			if s != "" {
				m[name] = s
				continue
			}
		}

		m[name] = fv.Interface()
	}

	return m
}
