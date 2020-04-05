package resterror

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	gopkgerrors "github.com/pkg/errors"
)

// RestError represents a Rest HTTP Error that can be returned from a controller
type RestError struct {
	Code          int      `json:"code"`
	Message       string   `json:"message"`
	Messages      []string `json:"messages"`
	OriginalError error    `json:"-"`
}

func (re *RestError) Error() string {
	return re.Message
}

// ErrorWrapper is an error trap that attempts to make sense out of common errors
// it handles errors and cleans up the messages before going out
func ErrorWrapper(err error, msg string) error {
	if err == nil {
		return nil
	}

	// Trap Common SQL Errors
	if err == sql.ErrNoRows {
		err = &RestError{
			Code:          404,
			Message:       fmt.Sprintf("Resource not found: %v", msg),
			OriginalError: err,
		}
	}

	// Trap Postgres SQL Errors
	if pgerr, ok := err.(*pq.Error); ok {
		switch pgerr.Code.Name() {
		case "unique_violation":
			err = &RestError{
				Code:          409,
				Message:       fmt.Sprintf("Record already exists: %s", msg),
				OriginalError: pgerr,
			}
		}
	}

	return gopkgerrors.Wrap(err, msg)
}
