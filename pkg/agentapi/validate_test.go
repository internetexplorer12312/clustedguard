package agentapi

import "testing"

func TestMetrics_StructFields(t *testing.T) {
	m := Metrics{CPUPercent: 10, MemoryUsedPercent: 20, DiskUsedPercent: 30}
	if m.CPUPercent != 10 {
		t.Fatalf("unexpected cpu %v", m.CPUPercent)
	}
}
