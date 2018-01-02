package input

import "reflect"

type Variable struct {
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options,omitempty"`
}

func (v Variable) IsEqual(other Variable) bool {
	return reflect.DeepEqual(v, other)
}
