package middleware

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/api/resterror"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Sugar to de-doup useage below
// be sure an return whenever you write
// or else multiple writes can go out
func writeResponse(w http.ResponseWriter, r *http.Request, status int, response []byte) {
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Content-Type", "text")
	w.WriteHeader(status)
	_, _ = w.Write(response)
	return
}

// ServeHTTP Serves up the HTTP response
func (fnH Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, payload, err := fnH.RouteHandler(fnH.AppCtx, w, r)

	if err == nil {
		fnH.AppCtx.Metrics.StatHTTPResponseCount.WithLabelValues(strconv.Itoa(status), r.RequestURI, r.Method, r.Proto).Inc()
		if payload != nil {
			response, _ := json.Marshal(payload)
			writeResponse(w, r, status, response)
			return
		}
		w.WriteHeader(status)
		return
	}

	// Get Cause of Original Error
	originalErr := errors.Cause(err)

	// Attempt to catch "known" errors.
	// these should typically be 4xx errors that were gift wrapped somewhere in the call stack
	// as a developer you can return &resterror.RestError{} from methods in the Controller/Service/DB
	// which will be caught here.
	// Typically, you will want to return those errors from the Service tier
	if rerr, ok := originalErr.(*resterror.RestError); ok {
		fnH.AppCtx.Logger.For(r.Context()).Error(
			"internal http rest error",
			zap.Error(rerr),
			zap.Int("code", rerr.Code),
			zap.String("user_agent", r.UserAgent()),
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.String("protocol", r.Proto))
		response, _ := json.Marshal(rerr)
		fnH.AppCtx.Metrics.StatHTTPResponseCount.WithLabelValues(strconv.Itoa(rerr.Code), r.RequestURI, r.Method, r.Proto).Inc()
		writeResponse(w, r, rerr.Code, response)
		return
	}

	// Oh No! Something unexpected happened here
	// these errors should sound an alarm (Prometheus, AlertManager, Logging Levels, etc.)
	// these are "unknown" server errors that need to be investigated
	// http.StatusInternalServerError
	fnH.AppCtx.Logger.For(r.Context()).Error(
		"unknown internal server error",
		zap.Error(err),
		zap.String("error_type", reflect.TypeOf(err).String()),
		zap.Int("code", http.StatusInternalServerError),
		zap.String("user_agent", r.UserAgent()),
		zap.String("uri", r.RequestURI),
		zap.String("method", r.Method),
		zap.String("protocol", r.Proto))

	response, _ := json.Marshal(resterror.RestError{
		Code:    http.StatusInternalServerError,
		Message: "Unknown internal server error has occurred",
	})
	fnH.AppCtx.Metrics.StatHTTPResponseCount.WithLabelValues("500", r.RequestURI, r.Method, r.Proto).Inc()
	writeResponse(w, r, http.StatusInternalServerError, response)
	return
}

func extractAuthBearerToken(r *http.Request) string {
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(auth) == 2 && auth[0] == "Bearer" {
		return auth[1]
	}
	return ""
}
