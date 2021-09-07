package z_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/youthlin/logs"
	"github.com/youthlin/logs/pkg/arg"
	"github.com/youthlin/z"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func TestDefaultConfig(t *testing.T) {
	c := z.DefaultConfig()
	logs.Assert(c != nil)
	logs.Assert(len(c.Zap) == 1)
	config := c.Zap[0]
	logs.Assert(config != nil)
	logs.Assert(config.Enable)
	logs.Assert(config.Level == zap.DebugLevel)
	logs.Assert(config.Output.Type == z.Console)
}
func TestConfig(t *testing.T) {
	type logConfig struct {
		Logs z.LogsConfig `json:"logs" yaml:"logs"`
	}
	{
		// json
		var config logConfig
		f, _ := os.Open("testdata/config.json")
		b, _ := ioutil.ReadAll(f)
		err := json.Unmarshal(b, &config)
		logs.Assert(err == nil)
		logs.Assert(config.Logs.Level != nil)
		logs.Assert(config.Logs.Level.Root == logs.LevelWarn)
		logs.Assert(config.Logs.Level.Loggers["github"] == logs.LevelInfo)

		logger := z.NewLogger(config.Logs.Zap[0:1])
		logger.Debug("Hello, Console Debug")
		logger.Info("Hello, Console Info")

		// logs
		z.SetConfig(&config.Logs)
		logs.Debug("Debug")
		logs.With("a", 1).Info("info")

		// 这些 level 应该是彩色的
	}
	{
		// yaml
		var config logConfig
		file, err := os.Open("testdata/config.yaml")
		logs.Assert(err == nil)
		bytes, err := ioutil.ReadAll(file)
		logs.Assert(err == nil)
		err = yaml.Unmarshal(bytes, &config)
		logs.Assert(err == nil)
		logs.Assert(config.Logs.Level != nil)
		logs.Assert(config.Logs.Level.Root == logs.LevelWarn)
		logs.Assert(config.Logs.Level.Loggers["github"] == logs.LevelInfo)

		z.SetConfig(&config.Logs)

		logs.Info(" json: %s", arg.JSON(config))

		logger := z.NewLogger(config.Logs.Zap[0:1])
		logger.Debug("Hello, Console")
		logger.Info("Hello, Console") // Level with color
	}
}
func TestMarshal(t *testing.T) {
	c := z.DefaultConfig()
	b, err := json.Marshal(c)
	logs.Assert(err == nil)

	t.Logf("%s", b)
	err = json.Unmarshal(b, &c)
	logs.AssertThen(err == nil, nil, func() {
		t.Logf("unmarshal err=%+v", err)
	})
	b, err = yaml.Marshal(c)
	logs.Assert(err == nil)

	t.Logf("%s", b)
	logs.Assert(yaml.Unmarshal(b, &c) == nil)
	logs.Assert(c.Zap[0].EncoderConfig.LineEnding == "\n")
}
