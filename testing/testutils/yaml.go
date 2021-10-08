package testutils

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// DecodeYaml decodes object from yaml string (KnownFields = true).
func DecodeYaml(source string, object interface{}) error {
	reader := strings.NewReader(source)
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	return decoder.Decode(object)
}
