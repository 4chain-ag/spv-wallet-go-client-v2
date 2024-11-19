package totp

import "errors"

var (
	// ErrMissingXpriv is returned when the xpriv is missing.
	ErrMissingXpriv = errors.New("xpriv is missing")
	// ErrContactPubKeyInvalid is returned when the contact's PubKey is invalid.
	ErrContactPubKeyInvalid = errors.New("contact's PubKey is invalid")
)
