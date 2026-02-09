package validator

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorMessage struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

type ValidationError struct {
	StatusCode   int            `json:"statusCode"`
	Message      string         `json:"message"`
	ErrorMessage []ErrorMessage `json:"errors"`
}

func (err *ValidationError) Error() string {
	return err.Message
}

// messageFunc generates a human-readable message from a FieldError.
type messageFunc func(fe validator.FieldError) string

// simple returns a messageFunc that ignores the FieldError and returns a static string.
func simple(msg string) messageFunc {
	return func(fe validator.FieldError) string { return fmt.Sprintf("%s: %s", fe.Field(), msg) }
}

// withParam returns a messageFunc that interpolates the field name and param.
func withParam(format string) messageFunc {
	return func(fe validator.FieldError) string {
		return fmt.Sprintf(format, fe.Field(), fe.Param())
	}
}

// tagMessages maps every common validator tag to a human-friendly message generator.
var tagMessages = map[string]messageFunc{
	// ── Required / presence ──
	"required":             simple("is required"),
	"required_if":          simple("is required based on a condition"),
	"required_unless":      simple("is required unless a condition is met"),
	"required_with":        simple("is required when related fields are present"),
	"required_with_all":    simple("is required when all related fields are present"),
	"required_without":     simple("is required when related fields are absent"),
	"required_without_all": simple("is required when all related fields are absent"),

	// ── Strings ──
	"alpha":           simple("must contain only alphabetic characters"),
	"alphanum":        simple("must contain only alphanumeric characters"),
	"alphaunicode":    simple("must contain only unicode alphabetic characters"),
	"alphanumunicode": simple("must contain only unicode alphanumeric characters"),
	"ascii":           simple("must contain only ASCII characters"),
	"lowercase":       simple("must be lowercase"),
	"uppercase":       simple("must be uppercase"),
	"contains":        withParam("%s must contain '%s'"),
	"containsany":     withParam("%s must contain at least one of '%s'"),
	"containsrune":    withParam("%s must contain the character '%s'"),
	"excludes":        withParam("%s must not contain '%s'"),
	"excludesall":     withParam("%s must not contain any of '%s'"),
	"excludesrune":    withParam("%s must not contain the character '%s'"),
	"startswith":      withParam("%s must start with '%s'"),
	"endswith":        withParam("%s must end with '%s'"),
	"startsnotwith":   withParam("%s must not start with '%s'"),
	"endsnotwith":     withParam("%s must not end with '%s'"),

	// ── Comparisons ──
	"eq":             withParam("%s must be equal to %s"),
	"eq_ignore_case": withParam("%s must be equal to %s (case-insensitive)"),
	"ne":             withParam("%s must not be equal to %s"),
	"ne_ignore_case": withParam("%s must not be equal to %s (case-insensitive)"),
	"gt":             withParam("%s must be greater than %s"),
	"gte":            withParam("%s must be greater than or equal to %s"),
	"lt":             withParam("%s must be less than %s"),
	"lte":            withParam("%s must be less than or equal to %s"),
	"min":            withParam("%s must be at least %s"),
	"max":            withParam("%s must be at most %s"),
	"len":            withParam("%s must have a length of %s"),

	// ── Field comparisons ──
	"eqfield":    withParam("%s must be equal to %s"),
	"nefield":    withParam("%s must not be equal to %s"),
	"gtfield":    withParam("%s must be greater than %s"),
	"gtefield":   withParam("%s must be greater than or equal to %s"),
	"ltfield":    withParam("%s must be less than %s"),
	"ltefield":   withParam("%s must be less than or equal to %s"),
	"eqcsfield":  withParam("%s must be equal to %s"),
	"necsfield":  withParam("%s must not be equal to %s"),
	"gtcsfield":  withParam("%s must be greater than %s"),
	"gtecsfield": withParam("%s must be greater than or equal to %s"),
	"ltcsfield":  withParam("%s must be less than %s"),
	"ltecsfield": withParam("%s must be less than or equal to %s"),

	// ── Enums / choices ──
	"oneof":            withParam("%s must be one of [%s]"),
	"excluded_with":    simple("must be empty when related fields are present"),
	"excluded_without": simple("must be empty when related fields are absent"),

	// ── Network ──
	"ip":       simple("must be a valid IP address"),
	"ipv4":     simple("must be a valid IPv4 address"),
	"ipv6":     simple("must be a valid IPv6 address"),
	"cidr":     simple("must be a valid CIDR notation"),
	"cidrv4":   simple("must be a valid CIDRv4 notation"),
	"cidrv6":   simple("must be a valid CIDRv6 notation"),
	"tcp_addr": simple("must be a valid TCP address"),
	"udp_addr": simple("must be a valid UDP address"),
	"mac":      simple("must be a valid MAC address"),
	"fqdn":     simple("must be a valid FQDN"),

	// ── Formats ──
	"email":        simple("must be a valid email address"),
	"url":          simple("must be a valid URL"),
	"http_url":     simple("must be a valid HTTP URL"),
	"uri":          simple("must be a valid URI"),
	"base64":       simple("must be a valid Base64 string"),
	"base64url":    simple("must be a valid Base64URL string"),
	"base64rawurl": simple("must be a valid raw Base64URL string"),
	"hexadecimal":  simple("must be a valid hexadecimal string"),
	"hex_color":    simple("must be a valid hex color code"),
	"rgb":          simple("must be a valid RGB color"),
	"rgba":         simple("must be a valid RGBA color"),
	"hsl":          simple("must be a valid HSL color"),
	"hsla":         simple("must be a valid HSLA color"),
	"json":         simple("must be valid JSON"),
	"jwt":          simple("must be a valid JWT token"),
	"html":         simple("must be valid HTML"),
	"html_encoded": simple("must be HTML-encoded"),
	"url_encoded":  simple("must be URL-encoded"),

	// ── Identifiers ──
	"uuid":   simple("must be a valid UUID"),
	"uuid3":  simple("must be a valid UUID v3"),
	"uuid4":  simple("must be a valid UUID v4"),
	"uuid5":  simple("must be a valid UUID v5"),
	"ulid":   simple("must be a valid ULID"),
	"isbn":   simple("must be a valid ISBN"),
	"isbn10": simple("must be a valid ISBN-10"),
	"isbn13": simple("must be a valid ISBN-13"),

	// ── Numbers ──
	"numeric": simple("must be a numeric value"),
	"number":  simple("must be a number"),
	"boolean": simple("must be a boolean value"),

	// ── Date / time ──
	"datetime": withParam("%s must match the format %s"),
	"timezone": simple("must be a valid IANA timezone"),

	// ── Phone ──
	"e164": simple("must be a valid E.164 phone number"),

	// ── Collections ──
	"unique": simple("must contain unique values"),
	"dive":   simple("has invalid nested elements"),

	// ── File ──
	"file":     simple("must be a valid file path"),
	"filepath": simple("must be a valid file path"),
	"image":    simple("must be a valid image"),
	"dir":      simple("must be a valid directory"),
	"dirpath":  simple("must be a valid directory path"),

	// ── Credit cards / finance ──
	"credit_card": simple("must be a valid credit card number"),

	// ── Country / locale ──
	"iso3166_1_alpha2":   simple("must be a valid ISO 3166-1 alpha-2 country code"),
	"iso3166_1_alpha3":   simple("must be a valid ISO 3166-1 alpha-3 country code"),
	"iso3166_1_numeric":  simple("must be a valid ISO 3166-1 numeric country code"),
	"bcp47_language_tag": simple("must be a valid BCP 47 language tag"),
	"country_code":       simple("must be a valid country code"),
	"locale":             simple("must be a valid locale"),

	// ── Custom tags (from original code) ──
	"timeformat":  simple("must match the format hh:mm"),
	"questionpos": simple("question position must increment by one"),
}

// getErrorMessage resolves a human-readable message for a single field error.
func getErrorMessage(fe validator.FieldError) string {
	if fn, ok := tagMessages[fe.Tag()]; ok {
		return fn(fe)
	}
	// Fallback for any unregistered / custom tags.
	if fe.Param() != "" {
		return fmt.Sprintf("%s failed on '%s' validation with param '%s'", fe.Field(), fe.Tag(), fe.Param())
	}
	return fmt.Sprintf("%s failed on '%s' validation", fe.Field(), fe.Tag())
}

// RegisterTagMessage lets callers add or override tag messages at runtime.
//
//	validation.RegisterTagMessage("my_custom_tag", func(fe validator.FieldError) string {
//	    return fmt.Sprintf("%s is not valid", fe.Field())
//	})
func RegisterTagMessage(tag string, fn messageFunc) {
	tagMessages[tag] = fn
}

// ValidateRequest converts an error into a structured *ValidationError
// if the underlying error is a validator.ValidationErrors.
// Otherwise it returns the original error unchanged.
func ValidateRequest(err error) error {
	if err == nil {
		return nil
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		msgs := make([]ErrorMessage, 0, len(ve))
		for _, fe := range ve {
			msgs = append(msgs, ErrorMessage{
				Field:   fieldName(fe),
				Tag:     fe.Tag(),
				Message: getErrorMessage(fe),
			})
		}
		return &ValidationError{
			StatusCode:   http.StatusUnprocessableEntity,
			Message:      "validation failed",
			ErrorMessage: msgs,
		}
	}

	// Handle JSON unmarshal / type errors gracefully.
	var unmarshalErr *UnmarshalError
	if errors.As(err, &unmarshalErr) {
		return &ValidationError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid request body",
			ErrorMessage: []ErrorMessage{
				{Field: unmarshalErr.Field, Message: unmarshalErr.Message},
			},
		}
	}

	return err
}

// fieldName returns the JSON field name when available via the Namespace,
// falling back to the struct field name.
func fieldName(fe validator.FieldError) string {
	ns := fe.Namespace()
	// Strip the top-level struct name (e.g. "CreateUserReq.Email" → "Email").
	if idx := strings.Index(ns, "."); idx != -1 {
		return ns[idx+1:]
	}
	return fe.Field()
}

// UnmarshalError wraps JSON decoding errors into a structured format.
type UnmarshalError struct {
	Field   string
	Message string
}

func (e *UnmarshalError) Error() string {
	return e.Message
}

func NewUnmarshalError(field, message string) *UnmarshalError {
	return &UnmarshalError{Field: field, Message: message}
}
