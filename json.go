package z

import (
	"encoding/json"
	"fmt"
)

// toJSON 将入参转为 json
func toJSON(args ...interface{}) []interface{} {
	var result = make([]interface{}, 0, len(args))
	for i, arg := range args {
		bytes, err := json.Marshal(arg)
		if err != nil {
			var msg = map[string]string{
				"index": fmt.Sprintf("%d", i),
				"value": fmt.Sprintf("%#v", arg),
				"err":   fmt.Sprintf("%v", err),
			}
			bytes, _ = json.Marshal(msg)
			result = append(result, fmt.Sprintf("[< toJSON error|%s >]", bytes))
		} else {
			result = append(result, string(bytes))
		}
	}
	return result
}

// errJSON 序列化为 JSON 时使用指定的格式转为 string
type errJSON struct {
	error
	verb string
}

// MarshalJSON 序列化为 JSON 时调用
func (e *errJSON) MarshalJSON() ([]byte, error) {
	if e == nil || e.error == nil {
		return []byte(`null`), nil
	}
	str := fmt.Sprintf(e.verb, e.error)
	str = fmt.Sprintf("%q", str)
	return []byte(str), nil
}

// Err 返回一个 error, 其序列化为 JSON 时，会使用指定的 verb 格式序列化为 JSON string
func Err(verb string, err error) error {
	return &errJSON{err, verb}
}
