package valueobject

import "regexp"

type Email string

func (email Email) IsValid() bool {
	emailRegexpValidation := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegexpValidation.MatchString(string(email))
}
