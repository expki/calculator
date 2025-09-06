package encoding

import (
	"errors"
	"fmt"
	"reflect"
)

// Engrain recursively maps data from a map[string]any into a destination struct pointer.
// It supports nested structs, slices, maps, pointers, and basic types.
func Engrain(data map[string]any, dst any) error {
	// Validate input data
	if data == nil {
		return nil
	}

	// Validate destination is a pointer
	dstType := reflect.TypeOf(dst)
	if dstType == nil || dstType.Kind() != reflect.Ptr {
		return errors.New("engrain dst must be a non-nil pointer")
	}

	// Validate pointer is not nil
	dstValue := reflect.ValueOf(dst)
	if dstValue.IsNil() {
		return errors.New("engrain dst must not be nil")
	}

	// Get the element the pointer points to
	dstElem := dstValue.Elem()
	if dstElem.Kind() != reflect.Struct {
		return errors.New("engrain dst must be a pointer to a struct")
	}

	return engrainStruct(data, dstElem)
}

// engrainStruct handles the actual mapping of data to struct fields
func engrainStruct(data map[string]any, dstStruct reflect.Value) error {
	dstType := dstStruct.Type()

	for key, value := range data {
		// Find the struct field by name
		field := dstStruct.FieldByName(key)
		if !field.IsValid() {
			// Field doesn't exist, skip it
			continue
		}

		if !field.CanSet() {
			// Field cannot be set (unexported), skip it
			continue
		}

		// Handle nil values
		if value == nil {
			// Set zero value for the field type
			field.Set(reflect.Zero(field.Type()))
			continue
		}

		// Set the field value based on its type
		if err := setFieldValue(field, reflect.ValueOf(value)); err != nil {
			// Find field info for better error message
			fieldInfo, _ := dstType.FieldByName(key)
			return fmt.Errorf("failed to set field %s: %w", fieldInfo.Name, err)
		}
	}

	return nil
}

// setFieldValue sets a struct field value, handling type conversions
func setFieldValue(field reflect.Value, value reflect.Value) error {
	fieldType := field.Type()
	valueType := value.Type()

	// Handle exact type match
	if valueType.AssignableTo(fieldType) {
		field.Set(value)
		return nil
	}

	// Handle different field types
	switch field.Kind() {
	case reflect.Struct:
		return handleStructField(field, value)
	case reflect.Ptr:
		return handlePointerField(field, value)
	case reflect.Slice:
		return handleSliceField(field, value)
	case reflect.Map:
		return handleMapField(field, value)
	case reflect.Interface:
		return handleInterfaceField(field, value)
	default:
		return handleBasicField(field, value)
	}
}

// handleStructField handles setting struct fields
func handleStructField(field reflect.Value, value reflect.Value) error {
	// Value must be a map for struct fields
	if value.Kind() != reflect.Map {
		return fmt.Errorf("expected map for struct field, got %v", value.Kind())
	}

	// Convert to map[string]any
	data, ok := value.Interface().(map[string]any)
	if !ok {
		return errors.New("failed to convert value to map[string]any")
	}

	return engrainStruct(data, field)
}

// handlePointerField handles setting pointer fields
func handlePointerField(field reflect.Value, value reflect.Value) error {
	elemType := field.Type().Elem()

	// Create a new instance of the pointed-to type
	newPtr := reflect.New(elemType)

	// Set the value to the new pointer's element
	if err := setFieldValue(newPtr.Elem(), value); err != nil {
		return err
	}

	field.Set(newPtr)
	return nil
}

// handleSliceField handles setting slice fields
func handleSliceField(field reflect.Value, value reflect.Value) error {
	// Handle []any as source
	if value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	if value.Kind() != reflect.Slice {
		return fmt.Errorf("expected slice for slice field, got %v", value.Kind())
	}

	elemType := field.Type().Elem()
	slice := reflect.MakeSlice(field.Type(), 0, value.Len())

	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)

		// If elem is interface{}, get the actual value
		if elem.Kind() == reflect.Interface {
			if elem.IsNil() {
				// Handle nil elements in slice
				var newElem reflect.Value
				switch elemType.Kind() {
				case reflect.Ptr:
					// For pointer types, append nil pointer
					newElem = reflect.Zero(elemType)
				default:
					// For non-pointer types, append zero value
					newElem = reflect.Zero(elemType)
				}
				slice = reflect.Append(slice, newElem)
				continue
			}
			elem = reflect.ValueOf(elem.Interface())
		}

		// Handle nil reflect.Value
		if !elem.IsValid() {
			var newElem reflect.Value
			switch elemType.Kind() {
			case reflect.Ptr:
				// For pointer types, append nil pointer
				newElem = reflect.Zero(elemType)
			default:
				// For non-pointer types, append zero value
				newElem = reflect.Zero(elemType)
			}
			slice = reflect.Append(slice, newElem)
			continue
		}

		// Create new element
		var newElem reflect.Value

		switch elemType.Kind() {
		case reflect.Struct:
			newElem = reflect.New(elemType).Elem()
			if err := setFieldValue(newElem, elem); err != nil {
				return fmt.Errorf("failed to set slice element %d: %w", i, err)
			}
		case reflect.Ptr:
			if elemType.Elem().Kind() == reflect.Struct {
				// Handle []*Struct
				newElem = reflect.New(elemType.Elem())
				if err := setFieldValue(newElem.Elem(), elem); err != nil {
					return fmt.Errorf("failed to set slice pointer element %d: %w", i, err)
				}
			} else {
				// Handle []*BasicType
				newElem = reflect.New(elemType.Elem())
				if err := setFieldValue(newElem.Elem(), elem); err != nil {
					return fmt.Errorf("failed to set slice pointer element %d: %w", i, err)
				}
			}
		default:
			// Handle []BasicType
			newElem = reflect.New(elemType).Elem()
			if err := setFieldValue(newElem, elem); err != nil {
				return fmt.Errorf("failed to set slice element %d: %w", i, err)
			}
		}

		slice = reflect.Append(slice, newElem)
	}

	field.Set(slice)
	return nil
}

// handleMapField handles setting map fields
func handleMapField(field reflect.Value, value reflect.Value) error {
	if value.Kind() != reflect.Map {
		return fmt.Errorf("expected map for map field, got %v", value.Kind())
	}

	keyType := field.Type().Key()
	elemType := field.Type().Elem()
	newMap := reflect.MakeMap(field.Type())

	for _, key := range value.MapKeys() {
		mapValue := value.MapIndex(key)

		// If mapValue is interface{}, get the actual value
		if mapValue.Kind() == reflect.Interface {
			if mapValue.IsNil() {
				// Handle nil values in map
				var newValue reflect.Value
				switch elemType.Kind() {
				case reflect.Ptr:
					// For pointer types, set nil pointer
					newValue = reflect.Zero(elemType)
				case reflect.Interface:
					// For interface types, set nil
					newValue = reflect.Zero(elemType)
				default:
					// For non-pointer types, set zero value
					newValue = reflect.Zero(elemType)
				}
				// Convert key if necessary
				var newKey reflect.Value
				if key.Type().AssignableTo(keyType) {
					newKey = key
				} else if key.Type().ConvertibleTo(keyType) {
					newKey = key.Convert(keyType)
				} else {
					return fmt.Errorf("cannot convert map key of type %v to %v", key.Type(), keyType)
				}
				newMap.SetMapIndex(newKey, newValue)
				continue
			}
			mapValue = reflect.ValueOf(mapValue.Interface())
		}

		// Handle nil reflect.Value
		if !mapValue.IsValid() {
			var newValue reflect.Value
			switch elemType.Kind() {
			case reflect.Ptr:
				// For pointer types, set nil pointer
				newValue = reflect.Zero(elemType)
			case reflect.Interface:
				// For interface types, set nil
				newValue = reflect.Zero(elemType)
			default:
				// For non-pointer types, set zero value
				newValue = reflect.Zero(elemType)
			}
			// Convert key if necessary
			var newKey reflect.Value
			if key.Type().AssignableTo(keyType) {
				newKey = key
			} else if key.Type().ConvertibleTo(keyType) {
				newKey = key.Convert(keyType)
			} else {
				return fmt.Errorf("cannot convert map key of type %v to %v", key.Type(), keyType)
			}
			newMap.SetMapIndex(newKey, newValue)
			continue
		}

		// Convert key if necessary
		var newKey reflect.Value
		if key.Type().AssignableTo(keyType) {
			newKey = key
		} else if key.Type().ConvertibleTo(keyType) {
			newKey = key.Convert(keyType)
		} else {
			return fmt.Errorf("cannot convert map key of type %v to %v", key.Type(), keyType)
		}

		// Create new value
		var newValue reflect.Value

		switch elemType.Kind() {
		case reflect.Struct:
			newValue = reflect.New(elemType).Elem()
			if err := setFieldValue(newValue, mapValue); err != nil {
				return fmt.Errorf("failed to set map value for key %v: %w", key.Interface(), err)
			}
		case reflect.Ptr:
			if elemType.Elem().Kind() == reflect.Struct {
				// Handle map[K]*Struct
				newValue = reflect.New(elemType.Elem())
				if err := setFieldValue(newValue.Elem(), mapValue); err != nil {
					return fmt.Errorf("failed to set map pointer value for key %v: %w", key.Interface(), err)
				}
			} else {
				// Handle map[K]*BasicType
				newValue = reflect.New(elemType.Elem())
				if err := setFieldValue(newValue.Elem(), mapValue); err != nil {
					return fmt.Errorf("failed to set map pointer value for key %v: %w", key.Interface(), err)
				}
			}
		case reflect.Interface:
			// Handle map[K]interface{}
			newValue = mapValue
		default:
			// Handle map[K]BasicType
			newValue = reflect.New(elemType).Elem()
			if err := setFieldValue(newValue, mapValue); err != nil {
				return fmt.Errorf("failed to set map value for key %v: %w", key.Interface(), err)
			}
		}

		newMap.SetMapIndex(newKey, newValue)
	}

	field.Set(newMap)
	return nil
}

// handleInterfaceField handles setting interface{} fields
func handleInterfaceField(field reflect.Value, value reflect.Value) error {
	field.Set(value)
	return nil
}

// handleBasicField handles setting basic type fields with type conversion
func handleBasicField(field reflect.Value, value reflect.Value) error {
	fieldType := field.Type()
	valueType := value.Type()

	// Try direct conversion
	if valueType.ConvertibleTo(fieldType) {
		field.Set(value.Convert(fieldType))
		return nil
	}

	// Handle numeric conversions
	if isNumeric(field.Kind()) && isNumeric(value.Kind()) {
		field.Set(value.Convert(fieldType))
		return nil
	}

	// Handle string to basic type conversions
	if value.Kind() == reflect.String {
		str := value.String()
		switch field.Kind() {
		case reflect.String:
			field.SetString(str)
			return nil
		case reflect.Bool:
			switch str {
			case "true", "1":
				field.SetBool(true)
			case "false", "0":
				field.SetBool(false)
			default:
				return fmt.Errorf("cannot convert string %q to bool", str)
			}
			return nil
		}
	}

	return fmt.Errorf("cannot convert %v to %v", valueType, fieldType)
}

// isNumeric checks if a kind is numeric
func isNumeric(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}
