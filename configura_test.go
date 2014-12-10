package configura

import (
	"os"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBasicLoading(t *testing.T) {
	Convey("Test the basic configuration", t, func() {
		type Config struct {
			SomeString   string
			SomeInt      int
			SomeBool     bool
			SomeDuration time.Duration
			SomeFloat32  float32
			SomeFloat64  float64
		}

		Convey("When all the variables are set", func() {
			expectedString := "this is just a test"
			expectedInt := 1
			expectedBool := true
			expectedDuration, err := time.ParseDuration("1s")
			So(err, ShouldBeNil)
			var (
				expectedFloat32 float32 = 3.2
				expectedFloat64 float64 = 6.4
			)

			envs := map[string]string{
				"TEST_SOMESTRING":   expectedString,
				"TEST_SOMEINT":      strconv.Itoa(expectedInt),
				"TEST_SOMEBOOL":     strconv.FormatBool(expectedBool),
				"TEST_SOMEDURATION": expectedDuration.String(),
				"TEST_SOMEFLOAT32":  strconv.FormatFloat(float64(expectedFloat32), 'f', 2, 32),
				"TEST_SOMEFLOAT64":  strconv.FormatFloat(expectedFloat64, 'f', 2, 64),
			}
			for k, v := range envs {
				err = os.Setenv(k, v)
				So(err, ShouldBeNil)
			}

			c := Config{}
			err = Load("TEST_", &c)
			So(err, ShouldBeNil)

			So(c.SomeString, ShouldEqual, expectedString)
			So(c.SomeInt, ShouldEqual, expectedInt)
			So(c.SomeBool, ShouldEqual, expectedBool)
			So(c.SomeDuration, ShouldEqual, expectedDuration)
			So(c.SomeFloat32, ShouldEqual, expectedFloat32)
			So(c.SomeFloat64, ShouldEqual, expectedFloat64)
		})

		Convey("When at least one variable is not set", func() {
			c := Config{}
			err := Load("SOMERANDOMSTUFF", &c)
			So(err, ShouldNotBeNil)
		})

		Convey("When the type doesn't match", func() {
			err := os.Setenv("TEST_SOMEINT", "this can not be an int")
			So(err, ShouldBeNil)

			c := Config{}
			err = Load("TEST_", &c)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestStructTagsLoading(t *testing.T) {
	Convey("Test the configuration with struct tags", t, func() {
		Convey("When it has different names", func() {
			type Config struct {
				Foo string `configura:"DN"`
				Bar string `configura:"DN,"`
			}
			expectedFoo := "fubar"

			err := os.Setenv("DN", expectedFoo)
			So(err, ShouldBeNil)

			c := Config{}
			err = Load("DOESNTMATTER", &c)
			So(err, ShouldBeNil)

			So(c.Foo, ShouldEqual, expectedFoo)
		})

		Convey("When it has defaults", func() {
			type Config struct {
				Foo string `configura:",sometesthere"`
			}

			c := Config{}
			err := Load("WHATEVER,ITWILLDEFAULT", &c)
			So(err, ShouldBeNil)

			So(c.Foo, ShouldEqual, "sometesthere")
		})

		Convey("When it has defaults and different names", func() {
			type Config struct {
				Foo string `configura:"ACME,corporation"`
			}

			Convey("First test without the env var set", func() {
				c := Config{}
				err := Load("", &c)
				So(err, ShouldBeNil)

				So(c.Foo, ShouldEqual, "corporation")
			})

			Convey("And test is with the env var set now", func() {
				expectedFoo := "more fubar"

				err := os.Setenv("ACME", expectedFoo)
				So(err, ShouldBeNil)

				c := Config{}
				err = Load("", &c)
				So(err, ShouldBeNil)

				So(c.Foo, ShouldEqual, expectedFoo)
			})
		})
	})
}
