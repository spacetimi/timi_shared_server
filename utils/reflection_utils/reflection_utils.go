package reflection_utils

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
)

func GetFieldByName(structPtr interface{}, fieldName string) (interface{}, error) {

	p := reflect.ValueOf(structPtr)
	if p.Kind() != reflect.Ptr {
		return nil, errors.New("not a struct ptr")
	}

	v := reflect.Indirect(p)
	if v.Kind() != reflect.Struct {
		return nil, errors.New("not a struct")
	}

	value := v.FieldByName(fieldName)
	if !value.CanInterface() {
		return nil, errors.New("inaccessible field: " + fieldName)
	}

	return value.Interface(), nil
}

func GetFieldValuesFromStructPtr(structPtr interface{}, fieldNames []string) ([]interface{}, error) {
    var fieldValues []interface{}

    p := reflect.ValueOf(structPtr)
    if p.Kind() != reflect.Ptr {
        return nil, errors.New("not a struct ptr")
    }

    v := reflect.Indirect(p)
    if v.Kind() != reflect.Struct {
        return nil, errors.New("not a struct")
    }

    for _, fieldName := range fieldNames {
        value := v.FieldByName(fieldName)
        if value.CanInterface() {
            fieldValues = append(fieldValues, value.Interface())
        } else {
            return nil, errors.New("inaccessible field: " + fieldName)
        }
    }

    return fieldValues, nil
}

func MarshalStructPtrToBson(s interface{}) (bson.M, error) {
	bsonMRepresentation := make(map[string]interface{})

	p := reflect.ValueOf(s)
	if p.Kind() != reflect.Ptr {
		return nil, errors.New("not a struct ptr")
	}

	v := reflect.Indirect(p)
	if v.Kind() != reflect.Struct {
		return nil, errors.New("not a struct")
	}

	numFields := v.NumField()
	for i := 0; i < numFields; i++ {
		value := v.Field(i)
		if value.CanInterface() {
			field := v.Type().Field(i)
			bsonTag, ok := field.Tag.Lookup("bson")
			if !ok || !strings.Contains(bsonTag, "ignore") {
				bsonMRepresentation[field.Name] = value.Interface()
			}
		}
	}

	return bsonMRepresentation, nil
}
