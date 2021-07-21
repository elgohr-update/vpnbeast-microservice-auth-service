package options

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
	"reflect"
	"strconv"
)

// getStringEnv gets the specific environment variables with default value, returns default value if variable not set
func getStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// getIntEnv gets the specific environment variables with default value, returns default value if variable not set
func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return convertStringToInt(value)
}

// convertStringToInt converts string environment variables to integer values
func convertStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		logger.Warn("an error occurred while converting from string to int. Setting it as zero",
			zap.String("error", err.Error()))
		i = 0
	}
	return i
}

// unmarshalConfig creates a new *viper.Viper and unmarshalls the config into struct using *viper.Viper
func unmarshalConfig(key string, value interface{}) error {
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	bindEnvs(sub)

	return sub.Unmarshal(value)
}

// bindEnvs takes *viper.Viper as argument and binds structs fields to environments variables to be able to override
// them using environment variables at the runtime
func bindEnvs(sub *viper.Viper) {
	opts := AuthServiceOptions{}
	fieldCount := reflect.TypeOf(opts).NumField()
	for i := 0; i < fieldCount; i++ {
		tag := reflect.TypeOf(opts).Field(i).Tag.Get("env")
		name := reflect.TypeOf(opts).Field(i).Name
		log.Printf("%s - %s\n", name, tag)
		_ = sub.BindEnv(name, tag)
	}
}
