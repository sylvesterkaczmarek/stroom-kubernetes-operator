package controllers

import (
	"testing"

	stroomv1 "github.com/gradata-systems/stroom-k8s-operator/api/v1"
)

func TestLogSenderSecurityArgs(t *testing.T) {
	tests := []struct {
		name     string
		settings stroomv1.LogSenderSettings
		want     string
	}{
		{
			name: "insecure by default",
			want: "--no-secure",
		},
		{
			name: "client certificate",
			settings: stroomv1.LogSenderSettings{Tls: &stroomv1.LogSenderTlsSettings{
				SecretName: "client-tls",
			}},
			want: `--secure --cert "/stroom-log-sender/certs/client.crt" --key "/stroom-log-sender/certs/client.key"`,
		},
		{
			name: "client certificate and CA",
			settings: stroomv1.LogSenderSettings{Tls: &stroomv1.LogSenderTlsSettings{
				SecretName:       "client-tls",
				CaCertificateKey: "ca.pem",
			}},
			want: `--secure --cert "/stroom-log-sender/certs/client.crt" --key "/stroom-log-sender/certs/client.key" --cacert "/stroom-log-sender/certs/ca.crt"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := getLogSenderSecurityArgs(test.settings); got != test.want {
				t.Fatalf("getLogSenderSecurityArgs() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestCreateLogSenderTlsVolume(t *testing.T) {
	if got := createLogSenderTlsVolume(&stroomv1.LogSenderTlsSettings{}); got != nil {
		t.Fatalf("expected no TLS volume, got %#v", got)
	}

	volume := createLogSenderTlsVolume(&stroomv1.LogSenderTlsSettings{SecretName: "client-tls"})
	if volume == nil || volume.Name != LogSenderTlsVolumeName || volume.Secret == nil {
		t.Fatalf("unexpected TLS volume: %#v", volume)
	}
	if volume.Secret.SecretName != "client-tls" {
		t.Fatalf("secret name = %q, want client-tls", volume.Secret.SecretName)
	}
	if len(volume.Secret.Items) != 2 ||
		volume.Secret.Items[0].Key != "tls.crt" || volume.Secret.Items[0].Path != "client.crt" ||
		volume.Secret.Items[1].Key != "tls.key" || volume.Secret.Items[1].Path != "client.key" {
		t.Fatalf("unexpected default secret projection: %#v", volume.Secret.Items)
	}
}

func TestCreateLogSenderTlsVolumeWithCustomKeys(t *testing.T) {
	volume := createLogSenderTlsVolume(&stroomv1.LogSenderTlsSettings{
		SecretName:       "client-tls",
		CertificateKey:   "cert.pem",
		PrivateKeyKey:    "key.pem",
		CaCertificateKey: "root.pem",
	})
	if volume == nil || volume.Secret == nil || len(volume.Secret.Items) != 3 {
		t.Fatalf("unexpected TLS volume: %#v", volume)
	}

	want := [][2]string{{"cert.pem", "client.crt"}, {"key.pem", "client.key"}, {"root.pem", "ca.crt"}}
	for i, item := range volume.Secret.Items {
		if item.Key != want[i][0] || item.Path != want[i][1] {
			t.Fatalf("secret item %d = %q -> %q, want %q -> %q", i, item.Key, item.Path, want[i][0], want[i][1])
		}
	}
}
