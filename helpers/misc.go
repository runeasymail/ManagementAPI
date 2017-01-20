package helpers

import "reflect"

func GetDBFields(el interface{}, key string) []string {

	select_elements := []string{}

	t := reflect.TypeOf(el)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(key)

		if tag != ""{
			select_elements = append(select_elements,  tag )
		}
	}

	return select_elements
}
