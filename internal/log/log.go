package log

import (
	"context"
	"log"
)

// Errorf sends an error message to the error management system.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	log.Printf("ERROR: "+format, args...)
}
