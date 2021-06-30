package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ReadYamlExpand(filename string, v interface{}) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	expand := os.ExpandEnv(string(buf))
	return yaml.Unmarshal([]byte(expand), v)
}

func ReadYaml(filename string, v interface{}) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buf, v)
}

func MustReadYaml(filename string, v interface{}) {
	if err := ReadYaml(filename, v); err != nil {
		panic(err)
	}
}

func MustReadYamlExpand(filename string, v interface{}) {
	if err := ReadYamlExpand(filename, v); err != nil {
		panic(err)
	}
}
