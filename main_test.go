package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"Code Smells\":0,\"Duplication\":0,"+
		"\"Maintainability(mins)\":0,\"Other Issues\":0,"+
		"\"Test Coverage(%)\":92}", w.Body.String())

}
