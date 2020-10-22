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
	assert.Equal(t, "{\"Maintainability\":\"A\","+
		"\"Remediation\":\"0.0 minute\","+
		"\"Technical Debt Ratio\":\"0.0 percent\"}",
		w.Body.String())

}
