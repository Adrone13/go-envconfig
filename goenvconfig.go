package goenvconfig

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const tag = "env"

// Load populates provided config struct with environment variables
// using `env:"ENV_VAR"` tag and converting them to specified types
func Load[T any](config *T) error {
	val := reflect.ValueOf(*config)

	for i := 0; i < val.NumField(); i++ {
		propName := val.Type().Field(i).Name
		propType := val.Type().Field(i).Type
		envVar := val.Type().Field(i).Tag.Get(tag)

		v := os.Getenv(envVar)
		if v == "" {
			return errors.New(fmt.Sprintf("Config error: var %s not found", envVar))
		}

		switch propType.String() {
		case "string":
			setProperty(config, propName, v)

		case "int":
			parsed, err := strconv.Atoi(v)
			if err != nil {
				return errors.New(fmt.Sprintf("Config error: var %s failed to convert to int", envVar))
			}

			setProperty(config, propName, parsed)

		default:
			panic(fmt.Sprintf("Config error: type %s not supported", propType))
		}
	}

	return nil
}

func setProperty[T any](i *T, propName string, propValue any) {
	reflect.ValueOf(i).Elem().FieldByName(propName).Set(reflect.ValueOf(propValue))
}
