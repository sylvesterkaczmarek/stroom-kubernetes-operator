package controllers

import (
	"strings"
	"testing"
)

func TestStroomClusterConfigUsesCurrentOpenIdSettings(t *testing.T) {
	data, err := StaticFiles.ReadFile("static_content/stroomcluster-config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	config := string(data)

	if strings.Contains(config, "useDefaultOpenIdCredentials") {
		t.Fatal("deprecated useDefaultOpenIdCredentials setting is still present")
	}
	for _, expected := range []string{
		`identityProviderType: "${IDENTITY_PROVIDER_TYPE:-INTERNAL_IDP}"`,
		`clientId: "${STROOM_OPERATOR_OPENID_CLIENT_ID}"`,
		`clientSecret: "${STROOM_OPERATOR_OPENID_CLIENT_SECRET}"`,
	} {
		if !strings.Contains(config, expected) {
			t.Fatalf("missing OpenID setting %q", expected)
		}
	}
}
