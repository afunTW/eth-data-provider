package config

import "reflect"

func GetMapStructureTag(obj interface{}) []string {
	t := reflect.TypeOf(obj)
	numFields := t.NumField()
	tags := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		tag := field.Tag.Get("mapstructure")
		tags[i] = tag
	}
	return tags
}
