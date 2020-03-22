package field

import "github.com/sirupsen/logrus"

type contextKey = string

const (
	LogFieldErrorType = contextKey("error_kind")
)

// ErrorFields create a set of field for an error event
func ErrorFields(event string, errorType string) logrus.Fields {
	return logrus.Fields{
		"event":           event,
		LogFieldErrorType: errorType,
	}
}
