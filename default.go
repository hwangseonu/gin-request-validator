package gin_validator

import (
	"errors"
	"regexp"
	"strconv"
)

func EmailValidator(name string, data interface{}, args ...string) error {
	emailRegex := regexp.MustCompile(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	if str, err := mustString(data); err != nil {
		return errors.New(name + err.Error())
	} else {
		if !emailRegex.MatchString(str) {
			return errors.New(str + " is not email")
		} else {
			return nil
		}
	}
}

func NotBlankValidator(name string, data interface{}, args ...string) error {
	if str, err := mustString(data); err != nil {
		return errors.New(name + err.Error())
	} else if str == "" {
		return errors.New(name + " should not blank")
	} else {
		return nil
	}
}

func MinValidator(name string, data interface{}, args ...string) error {
	if i, err := mustInt(data); err != nil {
		return errors.New(name + err.Error())
	} else {
		if min, err := strconv.Atoi(args[0]); err != nil {
			return errors.New("first argument of min validate must int")
		} else if min < i {
			return errors.New(name + " must greater than " + strconv.Itoa(min))
		} else {
			return nil
		}
	}
}

func MaxValidator(name string, data interface{}, args ...string) error {
	if i, err := mustInt(data); err != nil {
		return errors.New(name + err.Error())
	} else {
		if max, err := strconv.Atoi(args[0]); err != nil {
			return errors.New("first argument of max validate must int")
		} else if max < i {
			return errors.New(name + " must less than " + strconv.Itoa(max))
		} else {
			return nil
		}
	}
}

func PatternValidator(name string, data interface{}, args ...string) error {
	var str string
	var err error
	var r *regexp.Regexp

	if str, err = mustString(data); err != nil {
		return err
	}

	if r, err = regexp.Compile(args[0]); err != nil {
		return err
	}

	if !r.MatchString(str) {
		return errors.New(str + "not matched pattern")
	} else {
		return nil
	}
}
