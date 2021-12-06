package drip

// Code is an Code returned with errors from the Drip API.
// https://www.getdrip.com/docs/rest-api#errors
type Code string

const (
	// PresenceError means the attribute is required.
	PresenceError Code = "presence_error"
	// LengthError means the length of the attribute is out of bounds.
	LengthError Code = "length_error"
	// UniquenessError means the attribute must be unique.
	UniquenessError Code = "uniqueness_error"
	// EmailError means the attribute must be a valid email address.
	EmailError Code = "email_error"
	// URLError means the attribute must be a valid URL.
	URLError Code = "url_error"
	// DomainError means the attribute must be a valid domain name.
	DomainError Code = "domain_error"
	// TimeError means the attribute must be a valid time in ISO-8601 format.
	TimeError Code = "time_error"
	// EmailAddressListError means the attribute must be a valid comma-separated list of email addresses.
	EmailAddressListError Code = "email_address_list_error"
	// DaysOfTheWeekError means the attribute must be a valid days of the week mask of the format /\A(0|1){7}\z/ (excluding 0000000).
	DaysOfTheWeekError Code = "days_of_the_week_error"
	// UnavailableError means a resource has been disabled or deleted.
	UnavailableError Code = "unavailable_error"
	// FormatError means a resource identifier or object is not formatted correctly.
	FormatError Code = "format_error"
	// RangeError means a numeric value is out of range.
	RangeError Code = "range_error"
)

// CodeError is a error with a code.
// https://www.getdrip.com/docs/rest-api#errors
type CodeError struct {
	Code      string `json:"code"`
	Attribute string `json:"attribute"`
	Message   string `json:"message"`
}

// Error returns the error message.
func (e CodeError) Error() string {
	return e.Message
}
