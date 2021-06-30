package utils

import "errors"

func Assert(is bool, message string) {
	if is {
		panic(errors.New(message))
	}
}

func Asserts(message string, checks ...bool) {
	flg := false
	for _, is := range checks {
		if is {
			flg = true
			break
		}
	}
	if !flg {
		panic(errors.New(message))
	}
}
