package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	w := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/health", nil)
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	Health(ginContext)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"status": "ok"}`, w.Body.String())
}
