package user

import "errors"

type PhoneNumber struct {
	AreaCode    string `json:"area_code"`
	Number      string `json:"number"`
	CountryCode string `json:"country_code"`
}

// Validate validates a phone number
func (p PhoneNumber) Validate() error {
	//TODO: More detailed validation can be done here
	if len(p.AreaCode) == 0 {
		return errors.New("area code is required")
	}

	if len(p.Number) == 0 {
		return errors.New("number is required")
	}

	if len(p.CountryCode) == 0 {
		return errors.New("country code is required")
	}

	return nil
}
