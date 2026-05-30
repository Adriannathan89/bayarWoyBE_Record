package user_profile_service_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func newJSONContext(method, path string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}
	c.Request = httptest.NewRequest(method, path, bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	return c, rec
}
