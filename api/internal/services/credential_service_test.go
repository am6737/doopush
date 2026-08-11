package services

import (
	"reflect"
	"testing"

	"github.com/doopush/doopush/api/internal/models"
)

func TestNormalizeScopes(t *testing.T) {
	got, err := normalizeScopes([]string{models.ScopePushSend, models.ScopePushBroadcast, models.ScopePushSend})
	if err != nil {
		t.Fatalf("normalizeScopes returned error: %v", err)
	}
	want := models.StringList{models.ScopePushBroadcast, models.ScopePushSend}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeScopes() = %#v, want %#v", got, want)
	}
}

func TestNormalizeScopesRejectsUnknownAndEmpty(t *testing.T) {
	if _, err := normalizeScopes(nil); err == nil {
		t.Fatal("expected empty scopes to be rejected")
	}
	if _, err := normalizeScopes([]string{"app:delete"}); err == nil {
		t.Fatal("expected unknown scope to be rejected")
	}
}

func TestCredentialPrefixes(t *testing.T) {
	if models.AppKeyPrefix != "dp_ak_" {
		t.Fatalf("unexpected App Key prefix: %q", models.AppKeyPrefix)
	}
	if models.AppSecretPrefix != "dp_as_" {
		t.Fatalf("unexpected App Secret prefix: %q", models.AppSecretPrefix)
	}
	if models.InstallationTokenPrefix != "dp_ins_" {
		t.Fatalf("unexpected Installation Token prefix: %q", models.InstallationTokenPrefix)
	}
}
