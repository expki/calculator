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
		//fmt.Printf("key: %s, value: %+v, kind: %d\n", key, value, dstStruct.Kind())
		field := dstStruct.FieldByName(key)
		if !field.IsValid() {
			continue
		}
		if !field.CanSet() {
			continue
		}
		if field.Kind() == reflect.Struct {
			//fmt.Println("struct")
			if err = Engrain(value.(map[string]any), field.Addr().Interface()); err != nil {
				return err
			}
			continue
		}
		if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Struct {
			//fmt.Println("array struct")
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
			//fmt.Println("array pointer struct")
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
			//fmt.Println("map struct")
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
			//fmt.Println("map pointer struct")
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
			//fmt.Println("pointer value")
			ptrValue := reflect.New(field.Type().Elem())
			ptrValue.Elem().Set(reflect.ValueOf(value))
			field.Set(ptrValue)
		} else {
			//fmt.Println("value")
			field.Set(reflect.ValueOf(value))
		}
	}
	return nil
}
