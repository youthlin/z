package z_test

import (
	"encoding/json"
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
	type logConfig struct {
		Logs z.Configs `json:"logs" yaml:"logs"`
	}
	Convey("new config", t, func() {
		Convey("yaml", func() {
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

			var config logConfig
			err = yaml.Unmarshal(bytes, &config)
			if err != nil {
				t.Fatalf("Unmarshal fail|%+v", err)
				return
			}
			t.Logf("print: %v\n", config.Logs)
			b, _ := json.Marshal(config.Logs)
			t.Logf(" json: %s\n", b)
			Convey("Console", func() {
				logger := z.NewLogger(config.Logs[0:1])
				logger.Debug("Hello, Console") // Level Info
				logger.Info("Hello, Console")  // Level with color
			})
		})
		Convey("json", func() {
			f, _ := os.Open("config.json")
			b, _ := ioutil.ReadAll(f)
			var config logConfig
			err := json.Unmarshal(b, &config)
			So(err, ShouldBeNil)
			Convey("Console", func() {
				logger := z.NewLogger(config.Logs[0:1])
				logger.Debug("Hello, Console") // Level Info
				logger.Info("Hello, Console")  // Level with color
			})
		})
	})
}
