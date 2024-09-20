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
	dstStruct := reflect.ValueOf(dst).Elem()
	// engrain data into dst
	for key, value := range data {
		field := dstStruct.FieldByName(key)
		for i := 0; i < field.Type().NumField(); i++ {
			if field.Type().Field(i).Tag.Get("expki") == key {
				field = dstStruct.Field(i)
				break
			}
		}
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
		if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Struct {
			slice := reflect.MakeSlice(field.Type(), 0, field.Cap())
			if value == nil {
				continue
			}
			for idx := 0; idx < reflect.ValueOf(value).Len(); idx++ {
				itemStruct := reflect.New(field.Type().Elem()).Elem()
				if err = Engrain(reflect.ValueOf(value).Index(idx).Interface().(map[string]any), itemStruct.Addr().Interface()); err != nil {
					return err
				}
				slice = reflect.Append(slice, itemStruct)
			}
			field.Set(slice)
			continue
		}
		// todo: fix this mess
		if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Pointer {
			slice := reflect.MakeSlice(field.Type(), 0, field.Cap())
			if value == nil {
				continue
			}
			for idx := 0; idx < reflect.ValueOf(value).Len(); idx++ {
				itemStruct := reflect.New(field.Type().Elem().Elem()).Elem()
				if err = Engrain(reflect.ValueOf(value).Index(idx).Interface().(map[string]any), itemStruct.Addr().Interface()); err != nil {
					return err
				}
				ptrValue := reflect.New(itemStruct.Type())
				ptrValue.Elem().Set(itemStruct)
				slice = reflect.Append(slice, ptrValue)
			}
			field.Set(slice)
			continue
		}
		if field.Kind() == reflect.Map && field.Type().Elem().Kind() == reflect.Struct {
			mp := reflect.MakeMap(field.Type())
			for _, key := range reflect.ValueOf(value).MapKeys() {
				itemStruct := reflect.New(field.Type().Elem()).Elem()
				if err = Engrain(reflect.ValueOf(value).MapIndex(key).Interface().(map[string]any), itemStruct.Addr().Interface()); err != nil {
					return err
				}
				mp.SetMapIndex(key, itemStruct)
			}
			field.Set(mp)
			continue
		}
		if field.Kind() == reflect.Map && field.Type().Elem().Kind() == reflect.Pointer {
			mp := reflect.MakeMap(field.Type())
			for _, key := range reflect.ValueOf(value).MapKeys() {
				itemStruct := reflect.New(field.Type().Elem().Elem()).Elem()
				if err = Engrain(reflect.ValueOf(value).MapIndex(key).Interface().(map[string]any), itemStruct.Addr().Interface()); err != nil {
					return err
				}
				ptrValue := reflect.New(itemStruct.Type())
				ptrValue.Elem().Set(itemStruct)
				mp.SetMapIndex(key, ptrValue)
			}
			field.Set(mp)
			continue
		}
		if field.Kind() == reflect.Ptr {
			ptrValue := reflect.New(field.Type().Elem())
			ptrValue.Elem().Set(reflect.ValueOf(value))
			field.Set(ptrValue)
		} else {
			field.Set(reflect.ValueOf(value))
		}
	}
	return nil
}
