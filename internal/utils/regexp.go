package utils

import "regexp"

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var PhoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
var UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,16}$`)
