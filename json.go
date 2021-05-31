package z

import (
	"encoding/json"
	"fmt"
)

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
			bytes, err = json.Marshal(msg)
			result = append(result, fmt.Sprintf("[< toJSON error|%s >]", bytes))
		} else {
			result = append(result, fmt.Sprintf("%s", bytes))
		}
	}
	return result
}

type ErrJSON struct {
	error
	verb string
}

func (e *ErrJSON) MarshalJSON() ([]byte, error) {
	if e == nil || e.error == nil {
		return []byte(`null`), nil
	}
	str := fmt.Sprintf(e.verb, e.error)
	str = fmt.Sprintf("%q", str)
	return []byte(str), nil
}

func Err(verb string, err error) error {
	return &ErrJSON{err, verb}
}
