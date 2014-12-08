package configura

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func noDefaultsError(n, v string) error {
	return fmt.Errorf("%s doesn't have defaults and %s is not set", n, v)
}

func mismatchError(n string, i interface{}, t reflect.Kind) error {
	return fmt.Errorf("%s=%v must be %s", n, i, t)
}

func Load(prefix string, c interface{}) error {
	t := reflect.TypeOf(c)
	te := t.Elem()
	v := reflect.ValueOf(c)
	ve := v.Elem()

	if te.Kind() != reflect.Struct {
		return errors.New("the config must be a struct")
	}

	for i := 0; i < te.NumField(); i++ {
		name := te.Field(i).Name
		field := ve.FieldByName(name)
		// env vars will be uppercase
		varName := prefix + strings.ToUpper(name)
		env := os.Getenv(varName)
		if env == "" { // TODO: we will have default here
			return noDefaultsError(name, varName)
		}

		kind := field.Kind()

		switch kind {
		case reflect.String:
			field.SetString(env)
		case reflect.Int:
			n, err := strconv.Atoi(env)
			if err != nil {
				return mismatchError(name, n, kind)
			}
			field.SetInt(int64(n))
		}
	}

	return nil
}
