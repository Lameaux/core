package models

import "regexp"

var shortcodeRegex = regexp.MustCompile("^[0-9]{3,14}$")

func ValidateShortcode(shortcode string) bool {
	return shortcodeRegex.MatchString(shortcode)
}
