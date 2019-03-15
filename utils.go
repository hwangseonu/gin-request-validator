package gin_validator

import (
	"errors"
	"reflect"
	"strings"
)

func mustString(i interface{}) (string, error) {
	if s, ok := i.(string); ok {
		return s, nil
	} else {
		return "", errors.New("interface is not string")
	}
}

func mapToStruct(i interface{}, s interface{}) error {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	for k, v := range m {
		k = strings.ToUpper(string(rune(k[0]))) + k[1:]
		err := setField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	fieldVal := structValue.FieldByName(name)

	if !fieldVal.IsValid() {
		return errors.New("no such field: " + name + " in object")
	}
	if !fieldVal.CanSet() {
		return errors.New("cannot set " + name + " field value")
	}

	val := reflect.ValueOf(value)

	if fieldVal.Type() != val.Type() {
		if m, ok := value.(map[string]interface{}); ok {
			if fieldVal.Kind() == reflect.Struct {
				return mapToStruct(m, fieldVal.Addr().Interface())
			}
			if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Struct {
				if fieldVal.IsNil() {
					fieldVal.Set(reflect.New(fieldVal.Type()).Elem())
				}
				return mapToStruct(m, fieldVal.Interface())
			}
		}
		return errors.New("provided value type didn't match obj field type")
	}
	fieldVal.Set(val)
	return nil
}