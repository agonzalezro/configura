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
		}

		Convey("When all the variables are set", func() {
			expectedString := "this is just a test"
			expectedInt := 1
			expectedBool := true
			expectedDuration, err := time.ParseDuration("1s")
			So(err, ShouldBeNil)

			err = os.Setenv("TEST_SOMESTRING", expectedString)
			So(err, ShouldBeNil)
			err = os.Setenv("TEST_SOMEINT", strconv.Itoa(expectedInt))
			So(err, ShouldBeNil)
			err = os.Setenv("TEST_SOMEBOOL", strconv.FormatBool(expectedBool))
			So(err, ShouldBeNil)
			err = os.Setenv("TEST_SOMEDURATION", expectedDuration.String())
			So(err, ShouldBeNil)

			c := Config{}
			err = Load("TEST_", &c)
			So(err, ShouldBeNil)

			So(c.SomeString, ShouldEqual, expectedString)
			So(c.SomeInt, ShouldEqual, expectedInt)
			So(c.SomeBool, ShouldEqual, expectedBool)
			So(c.SomeDuration, ShouldEqual, expectedDuration)
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
}

func TestWithDefault(t *testing.T) {
}
