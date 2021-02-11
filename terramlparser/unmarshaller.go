package terramlparser

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

func CustomizedYAMLUnmarshal(in []byte, out interface{}) error {
	var res interface{}

	if err := yaml.Unmarshal(in, &res); err != nil {
		return err
	}
	*out.(*interface{}) = cleanupMapValue(res)

	return nil
}

func cleanupInterfaceArray(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = cleanupMapValue(v)
	}
	return res
}

func cleanupInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[fmt.Sprintf("%v", k)] = cleanupMapValue(v)
	}
	return res
}

func cleanupMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanupInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanupInterfaceMap(v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}