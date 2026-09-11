package v1

import (
	"testing"

	netv1 "k8s.io/api/networking/v1"
)

func TestIngressSettingsGetPathType(t *testing.T) {
	t.Run("uses route default when override is omitted", func(t *testing.T) {
		settings := IngressSettings{}
		if got := settings.GetPathType(netv1.PathTypeExact); got != netv1.PathTypeExact {
			t.Fatalf("expected Exact, got %q", got)
		}
		if got := settings.GetPathType(netv1.PathTypePrefix); got != netv1.PathTypePrefix {
			t.Fatalf("expected Prefix, got %q", got)
		}
	})

	t.Run("uses configured override", func(t *testing.T) {
		settings := IngressSettings{PathType: netv1.PathTypeImplementationSpecific}
		if got := settings.GetPathType(netv1.PathTypeExact); got != netv1.PathTypeImplementationSpecific {
			t.Fatalf("expected ImplementationSpecific, got %q", got)
		}
		if got := settings.GetPathType(netv1.PathTypePrefix); got != netv1.PathTypeImplementationSpecific {
			t.Fatalf("expected ImplementationSpecific, got %q", got)
		}
	})
}
