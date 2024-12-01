package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// loadConfigFromFile loads the configuration from the given file into the given configuration struct.
func loadConfigFromFile(configFile string, config interface{}) error {
	if configFile == "" {
		return errors.New("the configuration file is not set")
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	jsonData, err := convertToJson(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, config)
	if err != nil {
		return err
	}

	return nil

}

// loadConfigFromEnv loads the configuration from the environment variables into the given configuration struct.
func loadConfigFromEnv(config interface{}, prefix string) {
	v := reflect.ValueOf(config).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		envKey := prefix + "_" + fieldType.Name

		if field.Kind() == reflect.Struct {
			loadConfigFromEnv(field.Addr().Interface(), envKey)
		} else {
			envValue := getEnv(strings.ToUpper(envKey), "")
			if envValue != "" {
				setField(field, envValue)
			}
		}
	}
}

// getEnv returns the value of the environment variable if it exists, otherwise it returns the default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// setField sets the value of the field based on the field type.
func setField(field reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		if intValue, err := strconv.Atoi(value); err == nil {
			field.SetInt(int64(intValue))
		}
	case reflect.Map:
		// Convert the string to a map
		mapValue := map[string]string{}
		if err := yaml.Unmarshal([]byte(value), &mapValue); err != nil {
			panic(err)
		}
		field.Set(reflect.ValueOf(mapValue))
	default:
		panic("unhandled default case")
	}
}

func convertToJson(yamlData []byte) ([]byte, error) {
	var data interface{}

	err := yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, err
	}

	data = convertMapKeysToString(data)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func convertMapKeysToString(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v := range x {
			m[fmt.Sprintf("%v", k)] = convertMapKeysToString(v)
		}
		return m
	case []interface{}:
		for i, v := range x {
			x[i] = convertMapKeysToString(v)
		}
	}
	return i
}
