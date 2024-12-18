package env

import "reflect"

func ExpandStringWihEnvHook(f reflect.Type, _ reflect.Type, data any) (any, error) {
	if f.Kind() == reflect.String {
		return Expand(data.(string)), nil
	}
	return data, nil
}

func ExpandStringWithEnvHookFunc() func(f reflect.Type, t reflect.Type, data any) (any, error) {
	return ExpandStringWihEnvHook
}

func ExpandSliceWithEnvHook(f reflect.Type, _ reflect.Type, data any) (any, error) {
	if f.Kind() == reflect.Slice {
		var s []any
		data := reflect.ValueOf(data)
		length := data.Len()
		for i := 0; i < length; i++ {
			elem := reflect.ValueOf(data.Index(i).Interface())
			if elem.Kind() == reflect.String {
				s = append(s, Expand(elem.String()))
			} else if elem.Kind() == reflect.Slice {
				v, err := ExpandSliceWithEnvHook(elem.Type(), nil, elem.Interface())
				if err != nil {
					return nil, err
				}
				s = append(s, v)
			} else if elem.Kind() == reflect.Map {
				v, err := ExpandStringKeyMapWithEnvHook(elem.Type(), nil, elem.Interface())
				if err != nil {
					return nil, err
				}
				s = append(s, v)
			} else {
				s = append(s, elem.Interface())
			}
		}
		return s, nil
	}
	return data, nil
}

func ExpandSliceWithEnvHookFunc() func(f reflect.Type, t reflect.Type, data any) (any, error) {
	return ExpandSliceWithEnvHook
}

func ExpandStringKeyMapWithEnvHook(f reflect.Type, _ reflect.Type, data any) (any, error) {
	if f.Kind() == reflect.Map {
		m := make(map[string]any)
		iter := reflect.ValueOf(data).MapRange()
		for iter.Next() {
			k := iter.Key().String()
			v := reflect.ValueOf(iter.Value().Interface())
			if v.Kind() == reflect.String {
				m[k] = Expand(v.String())
			} else if v.Kind() == reflect.Slice {
				v, err := ExpandSliceWithEnvHook(v.Type(), nil, v.Interface())
				if err != nil {
					return nil, err
				}
				m[k] = v
			} else if v.Kind() == reflect.Map {
				v, err := ExpandStringKeyMapWithEnvHook(v.Type(), nil, v.Interface())
				if err != nil {
					return nil, err
				}
				m[k] = v
			} else {
				m[k] = v.Interface()
			}
		}
		return m, nil
	}
	return data, nil
}

func ExpandStringKeyMapWithEnvHookFunc() func(f reflect.Type, t reflect.Type, data any) (any, error) {
	return ExpandStringKeyMapWithEnvHook
}
