package encoding

import (
	"errors"
	"reflect"
)

// Engrain a map into a destination struct pointer
func Engrain(data map[string]any, dst any) (err error) {
	// validate input data
	if data == nil {
		return nil
	}
	// validate input dst
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return errors.New("engrain dst must be a pointer")
	}
	// validate pointer is not nil
	if reflect.ValueOf(dst).IsNil() {
		return errors.New("engrain dst must not be nil")
	}
	// validate pointer points to a struct
	dstStruct := reflect.ValueOf(dst).Elem()
	if dstStruct.Kind() != reflect.Struct {
		return errors.New("engrain dst must be a struct pointer")
	}
	// engrain data into dst
	for key, value := range data {
		field := dstStruct.FieldByName(key)
		if !field.IsValid() {
			continue
		}
		if !field.CanSet() {
			continue
		}
		if field.Kind() == reflect.Struct {
			if err = Engrain(value.(map[string]any), field.Addr().Interface()); err != nil {
				return err
			}
			continue
		}
		field.Set(reflect.ValueOf(value))
	}
	return nil
}
