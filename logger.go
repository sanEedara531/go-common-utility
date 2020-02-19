/*
 * Delhivery VMS
 *
 * Delhivery Vehicle Management Service worflow.
 *
 * API version: v1
 * Contact: platform@delhivery.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package common

import (
	"log"
	"net/http"
	"time"
)

import (
	"log"
	"net/http"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}