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
		return "", errors.New(" is must int")
	}
}

func mustInt(i interface{}) (int, error) {
	if reflect.TypeOf(i).Kind() == reflect.Float64 {
		f := i.(float64)
		d := int(f)
		return d, nil
	} else {
		if d, ok := i.(int); ok {
			return d, nil
		} else {
			return 0, errors.New(" is must int")
		}
	}
}

func toInt(str string) {

}

func mapToStruct(m map[string]interface{}, s interface{}) error {
	for k, v := range m {
		k = strings.ToUpper(string(rune(k[0]))) + k[1:]
		if err := setField(s, k, v); err != nil {
			return err
		}
	}
	return nil
}

func checkRequiredField(m map[string]interface{}, obj interface{}, name string) (bool, interface{}) {
	t := reflect.TypeOf(obj).Elem()
	field, ok := t.FieldByName(name)
	if !ok {
		return false, errors.New("cannot find field by name " + name)
	}
	for _, c := range strings.Split(field.Tag.Get("binding"), " ") {
		if c == "required" {
			for k, v := range m {
				k = strings.ToUpper(string(k[0])) + k[1:]
				if field.Name == k {
					return true, v
				}
			}
			return true, nil
		} else {
			for k, v := range m {
				k = strings.ToUpper(string(k[0])) + k[1:]
				if field.Name == k {
					return false, v
				}
			}
			return false, nil
		}
	}
	return false, nil
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

	if fieldVal.Type().Kind() == reflect.Int && val.Type().Kind() == reflect.Float64{
		f := val.Interface().(float64)
		i := int(f)
		val = reflect.ValueOf(i)
	}

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