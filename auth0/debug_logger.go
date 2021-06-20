package auth0

import (
	"encoding/json"
	"log"
	"os"
)

func shouldLog() bool {
	tflog := os.Getenv("TF_LOG_PROVIDER")
	if len(tflog) == 0 {
		tflog = os.Getenv("TF_LOG")
	}

	return tflog == "TRACE"
}

func TfLogString(context string, message string) {
	if shouldLog() {
		log.Printf("\n\n[%s] %s\n\n", context, message)
	}
}

func TfLogJson(context string, boxed interface{}) {
	if shouldLog() {
		jsonBytes, err := json.Marshal(boxed)

		if err != nil {
			log.Printf("[debug_logger] failed to marshal log: %v", err)
		} else {
			log.Printf("\n\n[%s] %s\n\n", context, string(jsonBytes))
		}
	}
}
