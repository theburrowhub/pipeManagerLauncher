package envvars

import (
	"os"
	"strings"
)

const prefix = "PIPELINE_"

var Variables map[string]string

// GetEnvVars reads all environment variables with the prefix "PIPELINE_TRAINER_"
// and returns them as a map[string]string
func GetEnvVars() {
	envVars := make(map[string]string)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 && strings.HasPrefix(parts[0], prefix) {
			key := strings.TrimPrefix(parts[0], prefix)
			envVars[key] = parts[1]
		}
	}

	Variables = envVars
}
