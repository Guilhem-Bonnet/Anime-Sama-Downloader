package app

import (
	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/ports"
)

var ErrNotFound = ports.ErrNotFound

// CodedError permet aux executors de renvoyer un code d'erreur stable,
// persisté dans Job.errorCode.
//
// Exemples de codes: invalid_params, http_status, network_error, io_error.
type CodedError struct {
	Code    string
	Message string
	Err     error
}

func (e *CodedError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err == nil {
		return e.Message
	}
	if e.Message == "" {
		return e.Err.Error()
	}
	return e.Message + ": " + e.Err.Error()
}

func (e *CodedError) Unwrap() error { return e.Err }

