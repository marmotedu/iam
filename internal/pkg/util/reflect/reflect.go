// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package reflect

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func ToGormDBMap(obj interface{}, fields []string) (map[string]interface{}, error) {
	reflectType := reflect.ValueOf(obj).Type()
	reflectValue := reflect.ValueOf(obj)
	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
		reflectValue = reflect.ValueOf(obj).Elem()
	}

	ret := make(map[string]interface{}, 0)
	for _, f := range fields {
		fs, exist := reflectType.FieldByName(f)
		if !exist {
			return nil, fmt.Errorf("unknow field " + f)
		}

		tagMap := parseTagSetting(fs.Tag)
		gormfiled, exist := tagMap["COLUMN"]
		if !exist {
			return nil, fmt.Errorf("undef gorm field " + f)
		}

		ret[gormfiled] = reflectValue.FieldByName(f)
	}
	return ret, nil
}

func parseTagSetting(tags reflect.StructTag) map[string]string {
	setting := map[string]string{}
	for _, str := range []string{tags.Get("sql"), tags.Get("gorm")} {
		if str == "" {
			continue
		}
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				setting[k] = strings.Join(v[1:], ":")
			} else {
				setting[k] = k
			}
		}
	}
	return setting
}

func GetObjFieldsMap(obj interface{}, fields []string) map[string]interface{} {
	ret := make(map[string]interface{})

	modelReflect := reflect.ValueOf(obj)
	if modelReflect.Kind() == reflect.Ptr {
		modelReflect = modelReflect.Elem()
	}

	modelRefType := modelReflect.Type()
	fieldsCount := modelReflect.NumField()
	var fieldData interface{}
	for i := 0; i < fieldsCount; i++ {
		field := modelReflect.Field(i)
		if len(fields) != 0 && !findString(fields, modelRefType.Field(i).Name) {
			continue
		}

		switch field.Kind() {
		case reflect.Struct:
			fallthrough
		case reflect.Ptr:
			fieldData = GetObjFieldsMap(field.Interface(), []string{})
		default:
			fieldData = field.Interface()
		}

		ret[modelRefType.Field(i).Name] = fieldData
	}

	return ret
}

func CopyObj(from interface{}, to interface{}, fields []string) (changed bool, err error) {
	fromMap := GetObjFieldsMap(from, fields)
	toMap := GetObjFieldsMap(to, fields)
	if reflect.DeepEqual(fromMap, toMap) {
		return false, nil
	}

	t := reflect.ValueOf(to).Elem()
	for k, v := range fromMap {
		val := t.FieldByName(k)
		val.Set(reflect.ValueOf(v))
	}
	return true, nil
}

// CopyObjViaYaml marshal "from" to yaml data, then unMarshal data to "to".
func CopyObjViaYaml(to interface{}, from interface{}) error {
	if from == nil || to == nil {
		return nil
	}

	data, err := yaml.Marshal(from)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, to)
}

// findString return true if target in slice, return false if not.
func findString(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}
