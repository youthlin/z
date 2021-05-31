package z_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/youthlin/z"
)

func TestErr(t *testing.T) {
	Convey("Err", t, func() {
		err := z.Err("%+v", err)
		So(err, ShouldNotBeNil)
		e, ok := err.(*z.ErrJSON)
		So(ok, ShouldBeTrue)
		So(e.Error(), ShouldEqual, "error")
		bytes, err := e.MarshalJSON()
		So(err, ShouldBeNil)
		t.Logf("%s\n", bytes)
		So(string(bytes), ShouldContainSubstring, "error")
	})
}
