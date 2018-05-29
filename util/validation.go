package util

import "regexp"

func IsEmail(email string) bool {
	format := `^([\w-_]+(?:\.[\w-_]+)*)@((?:[a-z0-9]+(?:-[a-zA-Z0-9]+)*)+\.[a-z]{2,6})$`
	if ok, _ := regexp.MatchString(format, email); !ok {
		return false
	}
	return true
}

func IsMobile(mobile string) bool {
	format := `^[0|1][0-9]{10}$`
	if ok, _ := regexp.MatchString(format, mobile); !ok {
		return false
	}

	return false
}
