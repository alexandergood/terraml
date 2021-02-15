package utils

import (
	"gopkg.in/yaml.v2"
)

func CustomizedYAMLUnmarshal(in []byte, out interface{}) error {
	var res interface{}

	if err := yaml.Unmarshal(in, &res); err != nil {
		return err
	}
	*out.(*interface{}) = CleanupMapValue(res)

	return nil
}