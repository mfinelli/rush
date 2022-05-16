package util

import (
	"errors"
	"testing"

	"github.com/spf13/viper"
)

func TestNormalizeServerPort(t *testing.T) {
	key := "server.port"

	tests := []struct {
		input  interface{}
		output int
	}{
		{7999, 7999},
		{"7998", 7998},
	}

	errTests := []struct {
		input interface{}
		err   error
	}{
		{
			"not an int",
			errors.New("strconv.Atoi: parsing \"not an int\": invalid syntax"),
		},
		{
			"1000",
			errors.New("server.port is out of range (1000)"),
		},
		{
			90000,
			errors.New("server.port is out of range (90000)"),
		},
	}

	for _, test := range tests {
		viper.Set(key, test.input)
		normalizeServerPort()
		if viper.Get(key) != test.output {
			t.Errorf("expected %v but got %v", test.output, viper.Get(key))
		}
	}

	for _, test := range errTests {
		viper.Set(key, test.input)
		err := normalizeServerPort()
		if err.Error() != test.err.Error() {
			t.Errorf("expected %v but got %v", test.err.Error(), err.Error())
		}
	}

}

func TestNormalizeViperStringToInt(t *testing.T) {
	tests := []struct {
		input  interface{}
		output int
	}{
		{8001, 8001},
		{"8000", 8000},
	}

	errTests := []struct {
		input interface{}
		err   error
	}{
		{
			"not an int",
			errors.New("strconv.Atoi: parsing \"not an int\": invalid syntax"),
		},
	}

	for _, test := range tests {
		viper.Set("test", test.input)
		normalizeViperStringToInt("test")
		if viper.Get("test") != test.output {
			t.Errorf("expected %v but got %v", test.output, viper.Get("test"))
		}
	}

	for _, test := range errTests {
		viper.Set("test", test.input)
		err := normalizeViperStringToInt("test")
		if err.Error() != test.err.Error() {
			t.Errorf("expected %v but got %v", test.err.Error(), err.Error())
		}
	}
}
