package helpers

import (
	"fmt"

	"github.com/pborman/uuid"
)

func RandomString() string {
	return uuid.New()
}

func RandomName(format string) string {
	return fmt.Sprintf(format, RandomString())
}
