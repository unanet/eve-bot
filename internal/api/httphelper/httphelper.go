package httphelper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"gitlab.unanet.io/devops/eve-bot/internal/api/resterror"
)

// AppErr is syntax sugar to return an Err
func AppErr(err error, message string) (int, interface{}, error) {
	return 0, nil, errors.Wrap(err, message)
}

// AppResponse is syntax sugar to return a Response
func AppResponse(code int, response interface{}) (int, interface{}, error) {
	return code, response, nil
}

// Vars gets the url variables
func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// Query gets the query parameter values
func Query(r *http.Request) map[string][]string {
	return r.URL.Query()
}

// ParseBody parses the incoming json request body
func ParseBody(r *http.Request, model interface{}) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		return &resterror.RestError{
			Code:          400,
			Message:       fmt.Sprintf("Invalid Post Body: %s", err),
			OriginalError: err,
		}
	}
	return nil
}
