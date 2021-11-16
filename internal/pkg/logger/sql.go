// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package logger

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"gorm.io/gorm/utils"
)

func isPrintable(s []byte) bool {
	for _, r := range s {
		if !unicode.IsPrint(rune(r)) {
			return false
		}
	}

	return true
}

// NULL defines a NULL string.
const NULL = "NULL"

var convertableTypes = []reflect.Type{reflect.TypeOf(time.Time{}), reflect.TypeOf(false), reflect.TypeOf([]byte{})}

// ExplainSQL explain a SQL.
// nolint: gocognit,gocyclo
func ExplainSQL(sql string, numericPlaceholder *regexp.Regexp, escaper string, avars ...interface{}) string {
	var convertParams func(interface{}, int)
	vars := make([]string, len(avars))

	convertParams = func(v interface{}, idx int) {
		switch v := v.(type) {
		case bool:
			vars[idx] = strconv.FormatBool(v)
		case time.Time:
			if v.IsZero() {
				vars[idx] = escaper + "0000-00-00 00:00:00" + escaper
			} else {
				vars[idx] = escaper + v.Format("2006-01-02 15:04:05.999") + escaper
			}
		case *time.Time:
			if v != nil {
				if v.IsZero() {
					vars[idx] = escaper + "0000-00-00 00:00:00" + escaper
				} else {
					vars[idx] = escaper + v.Format("2006-01-02 15:04:05.999") + escaper
				}
			} else {
				vars[idx] = NULL
			}
		case fmt.Stringer:
			vars[idx] = escaper + strings.ReplaceAll(fmt.Sprintf("%v", v), escaper, "\\"+escaper) + escaper
		case driver.Valuer:
			reflectValue := reflect.ValueOf(v)
			if v != nil && reflectValue.IsValid() && ((reflectValue.Kind() == reflect.Ptr && !reflectValue.IsNil()) || reflectValue.Kind() != reflect.Ptr) {
				r, _ := v.Value()
				convertParams(r, idx)
			} else {
				vars[idx] = NULL
			}
		case []byte:
			if isPrintable(v) {
				vars[idx] = escaper + strings.ReplaceAll(string(v), escaper, "\\"+escaper) + escaper
			} else {
				vars[idx] = escaper + "<binary>" + escaper
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			vars[idx] = utils.ToString(v)
		case float64, float32:
			vars[idx] = fmt.Sprintf("%.6f", v)
		case string:
			vars[idx] = escaper + strings.ReplaceAll(v, escaper, "\\"+escaper) + escaper
		default:
			rv := reflect.ValueOf(v)
			//nolint: nestif
			if v == nil || !rv.IsValid() || rv.Kind() == reflect.Ptr && rv.IsNil() {
				vars[idx] = NULL
			} else if valuer, ok := v.(driver.Valuer); ok {
				v, _ = valuer.Value()
				convertParams(v, idx)
			} else if rv.Kind() == reflect.Ptr && !rv.IsZero() {
				convertParams(reflect.Indirect(rv).Interface(), idx)
			} else {
				for _, t := range convertableTypes {
					if rv.Type().ConvertibleTo(t) {
						convertParams(rv.Convert(t).Interface(), idx)

						return
					}
				}
				vars[idx] = escaper + strings.ReplaceAll(fmt.Sprint(v), escaper, "\\"+escaper) + escaper
			}
		}
	}

	for idx, v := range avars {
		convertParams(v, idx)
	}

	//nolint: nestif
	if numericPlaceholder == nil {
		var idx int
		var newSQL strings.Builder

		for _, v := range []byte(sql) {
			if v == '?' {
				if len(vars) > idx {
					newSQL.WriteString(vars[idx])
					idx++

					continue
				}
			}
			newSQL.WriteByte(v)
		}

		sql = newSQL.String()
	} else {
		sql = numericPlaceholder.ReplaceAllString(sql, "$$$1$$")
		for idx, v := range vars {
			sql = strings.Replace(sql, "$"+strconv.Itoa(idx+1)+"$", v, 1)
		}
	}

	return sql
}
