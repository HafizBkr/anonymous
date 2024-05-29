package validator

import (
	"regexp"
	"time"

	"github.com/google/uuid"
)

func IsUUID(v string) bool {
	_, err := uuid.Parse(v)
	return err == nil
}

func IsEmptyString(v string) bool {
	return v == ""
}

func IsEqual(x any, y any) bool {
	return x == y
}

func IsNotEqual(x any, y any) bool {
	return x != y
}

func IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

func IsOneOf(x any, ys ...string) bool {
	ok := false
	for _, y := range ys {
		if x == y {
			ok = true
			break
		}
	}
	return ok
}

func IsValidDate(v string) bool {
	_, err := time.Parse("2006-01-02", v)
	return err == nil
}

func IsOfType[T any](v any) bool {
	_, ok := v.(T)
	return ok
}

func IsCorrectPhoneNumber(v string) bool {
  return true
}
