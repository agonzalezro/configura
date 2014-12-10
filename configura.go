package configura

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func noDefaultsError(n, v string) error {
	return fmt.Errorf("%s doesn't have defaults and %s is not set", n, v)
}

func mismatchError(n string, i interface{}, t reflect.Kind) error {
	return fmt.Errorf("%s=%v must be %s", n, i, t)
}

func getStructInfo(v reflect.StructField) (fieldName, envVar, defVal string) {
	fieldName = v.Name
	tags := strings.Split(v.Tag.Get("configura"), ",")
	envVar = tags[0]
	if len(tags) > 1 {
		defVal = tags[1]
	}
	return
}

// Load will go through all the fields defined in your struct and trying to
// load their config values from environemnt variables.
//
// - The var name to be looked up on the system can be override using struct
// tags: `configura:"OVERRIDE"`
//
// - The user will also be able to set some defaults int case that the variable
// was not found on the system: `configura:",defaultvalue"`
//
// - Or both: `configura:"OVERRIDE,defaultvalue"`
func Load(prefix string, c interface{}) error {
	t := reflect.TypeOf(c)
	te := t.Elem()
	v := reflect.ValueOf(c)
	ve := v.Elem()

	if te.Kind() != reflect.Struct {
		return errors.New("the config must be a struct")
	}

	for i := 0; i < te.NumField(); i++ {
		sf := te.Field(i)
		fieldName, envVar, defVal := getStructInfo(sf)

		field := ve.FieldByName(fieldName)

		if envVar == "" {
			envVar = prefix + strings.ToUpper(fieldName)
		}
		env := os.Getenv(envVar)

		if env == "" && defVal != "" {
			env = defVal
		} else if env == "" {
			return noDefaultsError(fieldName, envVar)
		}

		kind := field.Kind()

		switch kind {
		case reflect.String:
			field.SetString(env)
		case reflect.Int:
			n, err := strconv.Atoi(env)
			if err != nil {
				return mismatchError(fieldName, n, kind)
			}
			field.SetInt(int64(n))
		case reflect.Float32, reflect.Float64:
			bitSize := 32
			if kind == reflect.Float64 {
				bitSize = 64
			}
			n, err := strconv.ParseFloat(env, bitSize)
			if err != nil {
				return mismatchError(fieldName, n, kind)
			}
			field.SetFloat(n)
		case reflect.Bool:
			b, err := strconv.ParseBool(env)
			if err != nil {
				return mismatchError(fieldName, b, kind)
			}
			field.SetBool(b)
		case reflect.Int64: // time.Duration
			t, err := time.ParseDuration(env)
			if err != nil {
				return mismatchError(fieldName, t, kind)
			}
			field.Set(reflect.ValueOf(t))
		default:
			return fmt.Errorf("%s is not parsable", kind)
		}
	}

	return nil
}
