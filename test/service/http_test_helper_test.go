package service_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func newJSONContext(method string, path string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	return c, rec
}