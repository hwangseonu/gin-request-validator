package gin_validator

import (
	"errors"
	"regexp"
)

func EmailValidator(name string, data interface{}, args ...string) error {
	var emailRegex = regexp.MustCompile(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	str, err := mustString(data)
	if err != nil {
		return err
	}
	if !emailRegex.MatchString(str) {
		return errors.New(str + " is not email")
	}
	return nil
}
