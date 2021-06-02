package z_test

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/youthlin/z"
)

func TestErr(t *testing.T) {
	Convey("Err", t, func() {
		e := z.Err("%+v", err)
		So(e, ShouldNotBeNil)

		So(e.Error(), ShouldEqual, "error")
		bytes, err := json.Marshal(e)
		So(err, ShouldBeNil)
		t.Logf("%s\n", bytes)
		So(string(bytes), ShouldContainSubstring, "error")
	})
}
