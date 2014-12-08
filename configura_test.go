package configura

import (
	"os"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBasicLoading(t *testing.T) {
	type Config struct {
		SomeString string
		SomeInt    int
	}

	Convey("If the variables are set", t, func() {
		expectedString := "this is just a test"
		expectedInt := 1

		err := os.Setenv("TEST_SOMESTRING", expectedString)
		So(err, ShouldBeNil)
		err = os.Setenv("TEST_SOMEINT", strconv.Itoa(expectedInt))
		So(err, ShouldBeNil)

		c := Config{}
		err = Load("TEST_", &c)
		So(err, ShouldBeNil)

		So(c.SomeString, ShouldEqual, expectedString)
		So(c.SomeInt, ShouldEqual, expectedInt)
	})

	Convey("If at least one variable is not set", t, func() {
		c := Config{}
		err := Load("SOMERANDOMSTUFF", &c)
		So(err, ShouldNotBeNil)
	})

	Convey("If the type doesn't match", t, func() {
		err := os.Setenv("TEST_SOMEINT", "this can not be an int")
		So(err, ShouldBeNil)

		c := Config{}
		err = Load("TEST_", &c)
		So(err, ShouldNotBeNil)
	})
}

func TestStructTagsLoading(t *testing.T) {
}

func TestWithDefault(t *testing.T) {
}
