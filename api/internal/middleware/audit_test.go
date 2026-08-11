package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParseActionFromRequestAppSecretUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = httptest.NewRequest("PATCH", "/api/v1/apps/12/app-secrets/34", nil)

	action, resource, resourceID := parseActionFromRequest(context)
	if action != "update" || resource != "app_secret" || resourceID != "34" {
		t.Fatalf("parseActionFromRequest() = (%q, %q, %q), want (%q, %q, %q)", action, resource, resourceID, "update", "app_secret", "34")
	}
}
