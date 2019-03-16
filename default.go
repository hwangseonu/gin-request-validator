package gin_validator

import (
	"errors"
	"regexp"
	"strconv"
)

//구조체 필드가 이메일 형식을 지키는지 확인합니다.
//문자열이 아니거나 이메일 형식이 아니면 error 를 반환합니다.
//구조체 필드에 `validate:"email"` 과 같이 사용합니다.
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

// 구조체 필드가 비어있지 않는지 확인합니다.
// 문자열이 아니거나 비어있으면 error 를 반환합니다.
//구조체 필드에 `validate:"notblank"` 과 같이 사용합니다.
func NotBlankValidator(name string, data interface{}, args ...string) error {
	if str, err := mustString(data); err != nil {
		return errors.New(name + err.Error())
	} else if str == "" {
		return errors.New(name + " should not blank")
	} else {
		return nil
	}
}

// 구조체 필드가 첫번째 인자로 받은 값보다 크거나 같은지 확인합니다.
// Data 와 첫번때 인자가 정수형이 아니거나 Data 가 더 작으면 error 를 반환합니다.
//구조체 필드에 `validate:"min=1"` 과 같이 사용합니다.
func MinValidator(name string, data interface{}, args ...string) error {
	if i, err := mustInt(data); err != nil {
		return errors.New(name + err.Error())
	} else {
		if min, err := strconv.Atoi(args[0]); err != nil {
			return errors.New("first argument of min validate must int")
		} else if min > i {
			return errors.New(name + " must greater than " + strconv.Itoa(min))
		} else {
			return nil
		}
	}
}

// 구조체 필드가 첫번째 인자로 받은 값보다 작거나 같은지 확인합니다.
// Data 와 첫번때 인자가 정수형이 아니거나 Data 가 더 크면 error 를 반환합니다.
//구조체 필드에 `validate:"max=100"` 과 같이 사용합니다.
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

// 구조체 필드가 첫번째 인자로 정규식과 일치한지 확인합니다.
// Data 와 첫번때 인자가 문자열이 아니거나 Data 정규식을 지키지 않으면 error 를 반환합니다.
//구조체 필드에 `validate:"pattern=^[A-Z]$"` 과 같이 사용합니다.
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
