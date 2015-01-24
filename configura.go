package configura

import (
	"errors"
	"fmt"
	"log"
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

type parser struct {
	prefix string
}

func newParser(prefix string) *parser {
	return &parser{strings.ToUpper(prefix)}
}

// getValue will try to read the value from the env var. If it was not found
// and a default was set, the default value is going to be returned.
func (p *parser) getValue(v reflect.StructField) (value string, err error) {
	tags := strings.Split(v.Tag.Get("configura"), ",")
	env := tags[0]

	if env == "" {
		if p.prefix != "" {
			env = p.prefix + "_"
		}

		env += strings.ToUpper(v.Name)
	}
	value = os.Getenv(env)

	if value == "" {
		if len(tags) > 1 {
			value = tags[1]
			return
		}
		err = noDefaultsError(v.Name, env)
	}

	return
}

// Load will go through all the fields defined in your struct and trying to
// load their config values from environemnt variables.
//
// Bear in mind that a underscore will be appended to the prefix.
//
// - The var name to be looked up on the system can be override using struct
// tags: `configura:"OVERRIDE"`
//
// - The user will also be able to set some defaults int case that the variable
// was not found on the system: `configura:",defaultvalue"`
//
// - Or both: `configura:"OVERRIDE,defaultvalue"`
func Load(prefix string, c interface{}) (err error) {
	log.Println("---")

	p := newParser(prefix)
	log.Println("p")
	log.Println(p)

	t := reflect.TypeOf(c)
	log.Println("t")
	log.Println(t)
	te := t.Elem()
	log.Println(te)

	if te.Kind() != reflect.Struct {
		return errors.New("the config must be a struct")
	}

	v := reflect.ValueOf(c)
	ve := v.Elem()
	log.Println("v")
	log.Println(v)
	log.Println(ve)

	for i := 0; i < te.NumField(); i++ {
		sf := te.Field(i)

		name := sf.Name
		log.Println("name")
		log.Println(name)
		field := ve.FieldByName(name)
		if name == "typ" {
			panic(te.NumField())
		}

		kind := field.Kind()

		var value string
		if kind != reflect.Struct {
			value, err = p.getValue(sf) // no default for the struct kind
			if err != nil {
				return err
			}
		}

		switch kind {
		case reflect.String:
			field.SetString(value)
		case reflect.Int:
			n, err := strconv.Atoi(value)
			if err != nil {
				return mismatchError(name, n, kind)
			}
			field.SetInt(int64(n))
		case reflect.Float32, reflect.Float64:
			bitSize := 32
			if kind == reflect.Float64 {
				bitSize = 64
			}
			n, err := strconv.ParseFloat(value, bitSize)
			if err != nil {
				return mismatchError(name, n, kind)
			}
			field.SetFloat(n)
		case reflect.Bool:
			b, err := strconv.ParseBool(value)
			if err != nil {
				return mismatchError(name, b, kind)
			}
			field.SetBool(b)
		case reflect.Int64: // time.Duration
			t, err := time.ParseDuration(value)
			if err != nil {
				return mismatchError(name, t, kind)
			}
			field.Set(reflect.ValueOf(t))
		case reflect.Struct:
			subPrefix := p.prefix
			if subPrefix != "" {
				subPrefix += "_"
			}
			subPrefix += name
			log.Println("subPrefix")
			log.Println(subPrefix)

			if err := Load(subPrefix, field); err != nil {
				panic(err)
				return err
			}
		default:
			return fmt.Errorf("%s is not parsable", kind)
		}
	}

	return nil
}
