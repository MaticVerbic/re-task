package server

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

var RequestIDHeader = http.CanonicalHeaderKey("X-Request-Id")
var trueClientIP = http.CanonicalHeaderKey("True-Client-IP")
var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

// This middleware allows us to log context data about the request, the caller and calculate the round trip
// of the request through our system.
func requestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		s := time.Now()
		// most of the required data is already available from requests context.
		ctx := r.Context()
		requestID := getRequestID(r)
		realIP := getRealIP(r)

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		logger := logrus.WithFields(logrus.Fields{
			"request_id": requestID,
		})

		uri := fmt.Sprintf("%s://%s%s/", scheme, r.Host, r.RequestURI)

		// generate these fields separately, as we will only log them once, to reduce both the visual and memory clutter.
		// all request data can be traced using the request id.
		fields := logrus.Fields{
			"http_scheme": scheme,
			"http_proto":  r.Proto,
			"http_method": r.Method,
			"remote_addr": r.RemoteAddr,
			"request_id":  requestID,
			"user_agent":  r.UserAgent(),
			"real_ip":     realIP,
			"uri":         uri,
		}

		// log the only-once fields
		logger.WithFields(fields).Info("new http request")

		ctx = context.WithValue(r.Context(), "request_id", requestID)

		// defer the execution of this function until after the wrapper has run, this allows us to calculate the round trip
		// and log it.
		defer func(s time.Time, logger *logrus.Entry) {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"elapsed":    time.Since(s),
			}).Info("http request processed")
		}(s, logger)

		// next middleware
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

func getRequestID(r *http.Request) string {
	requestID := r.Header.Get(RequestIDHeader)
	if requestID == "" {
		var err error
		requestID, err = uuid.GenerateUUID()
		if err != nil {
			logrus.Errorf("failed to generate request id with err: %v", err)
			return ""
		}

	}

	return requestID
}

func getRealIP(r *http.Request) string {
	var ip string

	if tcip := r.Header.Get(trueClientIP); tcip != "" {
		ip = tcip
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ",")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	}
	if ip == "" || net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}

// CORSMiddleware returns Cors struct
func CORSMiddleware() func(next http.Handler) http.Handler {
	allowed := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type", "X-Request-Id",
			"X-Forwarded-For", "True-Client-IP", "X-Real-IP"},
		AllowCredentials: true,
	})
	return allowed.Handler
}
