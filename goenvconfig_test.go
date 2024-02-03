package goenvconfig

import (
	"fmt"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("DATABASE_USER", "postgres")

	config := struct {
		Port         int    `env:"PORT"`
		DatabaseUser string `env:"DATABASE_USER"`
	}{}

	err := Load(&config)

	assert(t, err, nil)

	if config.Port != 8080 {
		t.Errorf("Port should be equal to 8080")
	}
	if config.DatabaseUser != "postgres" {
		t.Errorf("DatabaseUser should be equal to \"postgres\"")
	}

	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_USER")
}

func TestMissingVar(t *testing.T) {
	config := struct {
		Port int `env:"PORT"`
	}{}

	err := Load(&config)
	if err == nil {
		t.Error("Expected to return error")
	} else {
		assert(t, "Config error: var PORT not found", err.Error())
	}
}

func TestInvalidInt(t *testing.T) {
	os.Setenv("PORT", "abc")

	config := struct {
		Port int `env:"PORT"`
	}{}

	err := Load(&config)
	if err == nil {
		t.Error("Expected to return error")
	}
	assert(t, "Config error: var PORT failed to convert to int", err.Error())

	os.Unsetenv("PORT")
}

func TestPanicForUnsupportedType(t *testing.T) {
	os.Setenv("TIMESTAMP", "1706978320428")

	config := struct {
		Timestamp int64 `env:"TIMESTAMP"`
	}{}

	assertPanic(t, func() {
		Load(&config)
	})
}

func assert[T comparable](t *testing.T, expected, received T) {
	if expected != received {
		t.Errorf(fmt.Sprintf("Expected: %+v\nReceived: %+v\n", expected, received))
	}
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}
