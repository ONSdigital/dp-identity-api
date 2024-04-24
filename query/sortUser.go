package query

import (
	"fmt"
	"github.com/ONSdigital/dp-identity-api/models"
	"reflect"
	"sort"
)

func sortBy(jsonField string, arr []models.UsersList) {
	if len(arr) < 1 {
		return
	}

	// first we find the field based on the json tag
	valueType := reflect.TypeOf(arr[0])

	var field reflect.StructField

	for i := 0; i < valueType.NumField(); i++ {
		field = valueType.Field(i)

		if field.Tag.Get("json") == jsonField {
			break
		}
	}

	// then we sort based on the type of the field
	sort.Slice(arr, func(i, j int) bool {
		v1 := reflect.ValueOf(arr[i]).FieldByName(field.Name)
		v2 := reflect.ValueOf(arr[j]).FieldByName(field.Name)

		switch field.Type.Name() {
		case "int":
			return int(v1.Int()) < int(v2.Int())
		case "string":
			return v1.String() < v2.String()
		case "bool":
			return !v1.Bool() // return small numbers first
		default:
			return false // return unmodified
		}
	})

	fmt.Printf("\nsort by %s:\n", jsonField)
	prettyPrint(arr)
}

func prettyPrint(arr []models.UsersList) {
	for _, v := range arr {
		fmt.Printf("%+v\n", v)
	}
}
