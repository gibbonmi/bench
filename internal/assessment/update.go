package assessment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func compatible(old, next Run) error {
	a, _ := json.Marshal(old)
	b, _ := json.Marshal(next)
	left, right := jsonObject(a), jsonObject(b)
	return preserve(left, right, "")
}

func preserve(old, next any, path string) error {
	if old == nil || old == "" {
		return nil
	}
	if reflect.DeepEqual(old, next) {
		return nil
	}
	if path == "/state" {
		return nil
	}
	if strings.HasSuffix(path, "/state") && old == "running" {
		return nil
	}
	switch a := old.(type) {
	case map[string]any:
		b, ok := next.(map[string]any)
		if !ok {
			break
		}
		for k, v := range a {
			if err := preserve(v, b[k], path+"/"+k); err != nil {
				return err
			}
		}
		return nil
	case []any:
		b, ok := next.([]any)
		if !ok || len(b) < len(a) {
			break
		}
		for i, v := range a {
			if err := preserve(v, b[i], fmt.Sprintf("%s/%d", path, i)); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("conflicting known evidence at %s", path)
}

func jsonObject(data []byte) any {
	var value any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	dec.Decode(&value)
	return value
}
