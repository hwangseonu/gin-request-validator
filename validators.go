package gin_validator

import (
	"errors"
	"regexp"
	"strconv"
)

var emailRegex = regexp.MustCompile(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)

func EmailValidator(name string, i ...interface{}) error {
	str, err := mustString(i[0])
	if err != nil {
		return err
	}
	if !emailRegex.MatchString(str) {
		return errors.New(str + " is not email")
	}
	return nil
}

func MinValidator(name string, i ...interface{}) error {
	data, err := mustInt(i[0])
	if err != nil {
		return err
	}
	min, err := mustInt(i[1])
	if err != nil {
		return err
	}
	if min > data {
		return errors.New(name + " must greater than " + strconv.Itoa(min) + ", but actual value is " + strconv.Itoa(data))
	}
	return nil
}

func MaxValidator(name string, i ...interface{}) error {
	data, err := mustInt(i[0])
	if err != nil {
		return err
	}
	max, err := mustInt(i[1])
	if err != nil {
		return err
	}
	if max < data {
		return errors.New(name + " must less than " + strconv.Itoa(max) + ", but actual value is " + strconv.Itoa(data))
	}
	return nil
}

func PatternValidator(name string, i ...interface{}) error {
	var data, regexStr string
	var err error
	var r *regexp.Regexp

	if data, err = mustString(i[0]); err != nil {
		return err
	}
	if regexStr, err = mustString(i[1]); err != nil {
		return err
	}
	if r, err = regexp.Compile(regexStr); err != nil {
		return nil
	}
	if r.MatchString(data) {
		return nil
	} else {
		return errors.New(name + "is not matched pattern")
	}
}

func NotBlackValidator(name string, i ...interface{}) error {
	if str, err := mustString(i); err != nil {
		return err
	} else if str == "" {
		return errors.New(name + " should not blank")
	} else {
		return nil
	}
}

func mustString(i interface{}) (string, error) {
	if s, ok := i.(string); ok {
		return s, nil
	} else {
		return "", errors.New("interface is not string")
	}
}

func mustFloat(i interface{}) (float64, error) {
	if s, ok := i.(float64); ok {
		return s, nil
	} else {
		return 0, errors.New("interface is not float")
	}
}

func mustInt(i interface{}) (int, error) {
	if s, ok := i.(int); ok {
		return s, nil
	} else {
		return 0, errors.New("interface is not integer")
	}
}