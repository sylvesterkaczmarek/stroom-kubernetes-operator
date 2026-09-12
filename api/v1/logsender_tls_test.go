package v1

import (
	"encoding/json"
	"testing"
)

func TestLogSenderTlsJSON(t *testing.T) {
	for _, configured := range []bool{false, true} {
		settings := LogSenderSettings{}
		if configured {
			settings.Tls = &LogSenderTlsSettings{SecretName: "client-tls"}
		}
		data, err := json.Marshal(settings)
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(data, &fields); err != nil {
			t.Fatal(err)
		}
		if _, present := fields["tls"]; present != configured {
			t.Fatalf("TLS configured=%v, unexpected JSON: %s", configured, data)
		}
		clone := settings.DeepCopy()
		if configured {
			clone.Tls.SecretName = "other-secret"
			if settings.Tls.SecretName != "client-tls" {
				t.Fatal("DeepCopy aliases TLS settings")
			}
		} else if !settings.Tls.IsZero() {
			t.Fatal("nil TLS settings must be unset")
		}
	}
}
