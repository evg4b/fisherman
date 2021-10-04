package testutils

import (
	"strings"

	"gopkg.in/yaml.v3"
)

func DecodeYaml(source string, object interface{}) error {
	reader := strings.NewReader(source)
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	return decoder.Decode(object)
}
