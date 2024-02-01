package goenvconfig

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const tag = "env"

func Load[T any](config *T) {
	val := reflect.ValueOf(*config)

	for i := 0; i < val.NumField(); i++ {
		propName := val.Type().Field(i).Name
		propType := val.Type().Field(i).Type
		envVar := val.Type().Field(i).Tag.Get(tag)

		v := os.Getenv(envVar)
		if v == "" {
			panic(fmt.Sprintf("Config error: var %s not found", envVar))
		}

		switch propType.String() {
		case "string":
			setProperty(config, propName, v)

		case "int":
			parsed, err := strconv.Atoi(v)
			if err != nil {
				panic(fmt.Sprintf(`Config error: var "%s" failed to convert to int`, envVar))
			}

			setProperty(config, propName, parsed)

		default:
			panic(fmt.Sprintf("Config error: type %s not supported", propType))
		}
	}
}

func setProperty[T any](i *T, propName string, propValue any) *T {
	reflect.ValueOf(i).Elem().FieldByName(propName).Set(reflect.ValueOf(propValue))

	return i
}
