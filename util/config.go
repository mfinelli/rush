package util

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/spf13/viper"
)

func normalizeServerPort() error {
	key := "server.port"

	if err := normalizeViperStringToInt(key); err != nil {
		return err
	}

	if viper.Get(key).(int) < 1024 || viper.Get(key).(int) > 65535 {
		return fmt.Errorf("%s is out of range (%d)", key, viper.Get(key))
	}

	return nil
}

// env variables get parsed as strings, but we expect some variables to be
// integers so convert them if necessary
func normalizeViperStringToInt(v string) error {
	original := viper.Get(v)

	if reflect.TypeOf(original).Kind() == reflect.String {
		intval, err := strconv.Atoi(original.(string))
		if err != nil {
			return err
		}

		viper.Set(v, intval)
	}

	return nil
}
