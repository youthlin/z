package z_test

import (
	"encoding/json"
	"fmt"
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
	assert(c != nil)
	assert(c.Level != nil)
	assert(c.Level.Root == logs.LevelDebug)
	assert(c.Zap != nil)
	assert(len(c.Zap) == 1)
	zapConfig := c.Zap[0]
	assert(zapConfig != nil)
	assert(zapConfig.Enable)
	assert(zapConfig.Level == zap.DebugLevel)
	assert(zapConfig.Output.Type == z.Console)
	assert(zapConfig.Output.File.Filename == z.Stdout)
}

func TestEncoder(t *testing.T) {
	c := z.DefaultConfig()
	e := c.Zap[0].Encoder
	logs.Info("Encoder Config 可以转为 json: %v", arg.JSON(e))
	e = &z.Encoder{}
	var s = []byte(`{
		"levelEncoder": "capitalColor",
		"timeEncoder": "rfc3339nano",
		"durationEncoder": "string"
		}`)
	err := json.Unmarshal(s, &e)
	assert(err == nil)
	logs.Info("能够反序列化，刚反序列化时没有merge默认值 Encoder Config : %v", arg.JSON(e))
	e.Zap()
	logs.Info("调用 Zap 后有默认值 Encoder Config : %v", arg.JSON(e))

	{

		c := z.DefaultConfig()
		e := c.Zap[0].Encoder
		b, err := yaml.Marshal(e)
		assert(err == nil)
		logs.Info("to yaml: %s", b)
		e = &z.Encoder{}
		err = yaml.Unmarshal(b, e)
		assertThen(err == nil, nil, func() {
			logs.Error("unmarshal err: %+v", err)
		})
		assert(e.CallerKey == "C")
	}

	s = []byte(`"nanos"`)
	var te z.TimeEncoder
	err = json.Unmarshal(s, &te)
	assert(err == nil)
	assert(te.Name == "nanos")
	s = []byte(`{"layout":"2006-01-02,15:04:05"}`)
	te = z.TimeEncoder{}
	err = json.Unmarshal(s, &te)
	assert(err == nil)
	assert(te.Name == "")
	assert(te.Layout == "2006-01-02,15:04:05")
}

func TestOutput(t *testing.T) {
	c := z.DefaultConfig()
	o := c.Zap[0].Output
	logs.Info("Output 可以序列化 to json: %v", arg.JSON(o))

	o = &z.Output{}
	json.Unmarshal([]byte(`{
		"type": "file",
		"file": {
			"filename": "test.log"
		}
	}`), o)
	logs.Info("反序列化 json: %v", arg.JSON(o))
	o.WriteSyncer()
	logs.Info("after WriteSyncer to json: %v", arg.JSON(o))
}

func TestUnmarshal(t *testing.T) {
	type AppConfig struct {
		Logs z.LogsConfig `json:"logs" yaml:"logs"`
	}
	{
		// json
		var appConfig AppConfig
		f, err := os.Open("testdata/config.json")
		assert(err == nil)
		b, err := ioutil.ReadAll(f)
		assert(err == nil)
		err = json.Unmarshal(b, &appConfig)
		assert(err == nil)
		assert(appConfig.Logs.Level != nil)
		assert(appConfig.Logs.Level.Root == logs.LevelWarn)
		assert(appConfig.Logs.Level.Loggers["github"] == logs.LevelInfo)

		assert(len(appConfig.Logs.Zap) == 3)
		assert(appConfig.Logs.Zap[0].Output.Type == z.Console)
		assert(appConfig.Logs.Zap[0].Output.File.Filename == z.Stdout)
		assert(appConfig.Logs.Zap[0].Encoder.TimeKey == "")
		appConfig.Logs.Zap[0].Encoder.Zap()
		assert(appConfig.Logs.Zap[0].Encoder.TimeKey == "ts")
		assert(appConfig.Logs.Zap[1].Output.Type == z.File)
		assert(appConfig.Logs.Zap[1].Output.File.Filename == "app.log")

		// 以下日志打印的 level 在控制台中应该是彩色的
		consoleLogger := z.NewLogger(appConfig.Logs.Zap[:1])
		consoleLogger.Debug("Hello, Console Debug")
		consoleLogger.Info("Hello, Console Info")

		// logs
		z.SetConfig(&appConfig.Logs)
		logs.Debug("Debug") // output to console, and file
		logs.With("a", 1).Info("info")
		logs.Error("error = %+v", fmt.Errorf("this is a error"))
	}
	{
		// yaml
		var appConfig AppConfig
		file, err := os.Open("testdata/config.yaml")
		assert(err == nil)
		bytes, err := ioutil.ReadAll(file)
		assert(err == nil)
		err = yaml.Unmarshal(bytes, &appConfig)
		assertThen(err == nil, nil, func() {
			logs.Error("yaml 反序列化失败: %+v", err)
		})
		assert(appConfig.Logs.Level != nil)
		assert(appConfig.Logs.Level.Root == logs.LevelWarn)
		assert(appConfig.Logs.Level.Loggers["github"] == logs.LevelInfo)

		assert(len(appConfig.Logs.Zap) == 3)
		assert(appConfig.Logs.Zap[0].Output.Type == z.Console)
		assert(appConfig.Logs.Zap[0].Output.File.Filename == z.Stdout)
		assert(appConfig.Logs.Zap[0].Encoder.TimeKey == "")
		appConfig.Logs.Zap[0].Encoder.Zap()
		assert(appConfig.Logs.Zap[0].Encoder.TimeKey == "ts")
		assert(appConfig.Logs.Zap[1].Output.Type == z.File)
		assert(appConfig.Logs.Zap[1].Output.File.Filename == "app.log")

		z.SetConfig(&appConfig.Logs)

		logs.Info(" json: %s", arg.JSON(appConfig))

		logger := z.NewLogger(appConfig.Logs.Zap[0:1])
		logger.Debug("Hello, Console")
		logger.Info("Hello, Console") // Level with color
	}
}

func TestMarshal(t *testing.T) {
	c := z.DefaultConfig()
	b, err := json.Marshal(c)
	assert(err == nil)
	logs.Info("default config to json = %s", b)

	c = &z.LogsConfig{}
	err = json.Unmarshal(b, c)
	assertThen(err == nil, nil, func() {
		t.Logf("unmarshal err=%+v", err)
	})
	logs.Info("unmarshal from json: %s", arg.JSON(c))
	assert(c.Zap[0].Encoder.LineEnding == "\n")

	b, err = yaml.Marshal(c)
	assert(err == nil)
	logs.Info("default config to yaml = %s", b)

	c = &z.LogsConfig{}
	assert(yaml.Unmarshal(b, c) == nil)
	logs.Info("unmarshal from yaml =(json) %s", arg.JSON(c))
	e := c.Zap[0].Encoder
	logs.Info("encoder: %s", arg.JSON(e))
	logs.Info("LineEnding = [%q]", e.LineEnding)
	assert(e.LineEnding == "\n")
}
