package z_test

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/youthlin/z"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func TestDefaultConfig(t *testing.T) {
	Convey("DefaultConfig", t, func() {
		c := z.DefaultConfig()
		So(c, ShouldNotBeNil)
		So(len(c), ShouldEqual, 1)
		config := c[0]
		So(config, ShouldNotBeNil)
		So(config.Enable, ShouldBeTrue)
		So(config.Level, ShouldEqual, zap.DebugLevel)
		So(config.Output.Type, ShouldEqual, z.Console)
	})
}
func TestConfig(t *testing.T) {
	Convey("new config", t, func() {
		file, err := os.Open("config.yaml")
		if err != nil {
			t.Fatalf("Open fail|%+v", err)
			return
		}
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			t.Fatalf("ReadAll fail|%+v", err)
			return
		}
		type logConfig struct {
			Logs z.Configs `yaml:"logs"`
		}
		var config logConfig
		err = yaml.Unmarshal(bytes, &config)
		if err != nil {
			t.Fatalf("Unmarshal fail|%+v", err)
			return
		}
		t.Logf("%v\n", config.Logs)
		Convey("Console", func() {
			logger := z.NewLogger(config.Logs[0:1])
			logger.Debug("Hello, Console") // Level Info
			logger.Info("Hello, Console")  // Level with color
		})
	})
}
