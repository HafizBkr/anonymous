package helpers

import (
	"fmt"
	"time"
)

func GenerateRandomCode() string {
  return fmt.Sprint(time.Now().Nanosecond())[:6]
}
