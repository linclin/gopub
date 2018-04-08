package gokits

import "regexp"

const regex_strict_email_pattern = `(?i)[A-Z0-9!#$%&'*+/=?^_{|}~-]+` +
	`(?:\.[A-Z0-9!#$%&'*+/=?^_{|}~-]+)*` +
	`@(?:[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?\.)+` +
	`[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?`

var (
	digitRegexp             = regexp.MustCompile(`^[0-9]+$`)
	chletterNumUnlineRegexp = regexp.MustCompile(`^[\x{4e00}-\x{9fa5}_a-zA-Z0-9]+$`)
	letterNumUnlineRegexp   = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	cheseRegexp             = regexp.MustCompile(`^[\x{4e00}-\x{9fa5}]+$`)
	emailRegexp             = regexp.MustCompile(`(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`)
	strictemailRegexp       = regexp.MustCompile(regex_strict_email_pattern)
	urlRegexp               = regexp.MustCompile(`(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?`)
)

func IsDigit(str string) bool {
	return digitRegexp.MatchString(str)
}

func IsChineseLetterNumUnline(str string) bool {
	return chletterNumUnlineRegexp.MatchString(str)
}

func IsLetterNumUnline(str string) bool {
	return letterNumUnlineRegexp.MatchString(str)
}

func IsChinese(str string) bool {
	return cheseRegexp.MatchString(str)
}

func IsEmail(str string) bool {
	return emailRegexp.MatchString(str)
}

func IsEmailRFC(str string) bool {
	return strictemailRegexp.MatchString(str)
}

func IsURL(str string) bool {
	return urlRegexp.MatchString(str)
}

func IfEmpty(a, b string) string {
	if a == "" {
		return b
	}
	return a
}
