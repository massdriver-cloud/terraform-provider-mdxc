package azure

import (
	"net"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
)

func responseWasNotFound(resp autorest.Response) bool {
	return responseWasStatusCode(resp, http.StatusNotFound)
}

func responseWasBadRequest(resp autorest.Response) bool {
	return responseWasStatusCode(resp, http.StatusBadRequest)
}

func responseWasForbidden(resp autorest.Response) bool {
	return responseWasStatusCode(resp, http.StatusForbidden)
}

func responseWasConflict(resp autorest.Response) bool {
	return responseWasStatusCode(resp, http.StatusConflict)
}

func responseErrorIsRetryable(err error) bool {
	if arerr, ok := err.(autorest.DetailedError); ok {
		err = arerr.Original
	}

	// nolint gocritic
	switch e := err.(type) {
	case net.Error:
		if e.Temporary() || e.Timeout() {
			return true
		}
	}

	return false
}

func responseWasStatusCode(resp autorest.Response, statusCode int) bool { // nolint: unparam
	if r := resp.Response; r != nil {
		if r.StatusCode == statusCode {
			return true
		}
	}

	return false
}
