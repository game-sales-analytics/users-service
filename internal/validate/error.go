package validate

import (
	"fmt"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("'%s' for field: '%s'", e.Message, e.Field)
}
