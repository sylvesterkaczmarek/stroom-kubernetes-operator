package v1

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNodeSetHeadlessServiceHostName(t *testing.T) {
	cluster := &StroomCluster{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "example"}}
	nodeSet := &NodeSet{Name: "data"}

	want := "stroom-test-node-data.example.svc.cluster.local"
	if got := cluster.GetNodeSetHeadlessServiceHostName(nodeSet); got != want {
		t.Fatalf("GetNodeSetHeadlessServiceHostName() = %q, want %q", got, want)
	}
}
