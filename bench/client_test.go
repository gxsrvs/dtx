package bench

import (
	"encoding/json"
	"fmt"
	"testing"
)

var clientProfiles = []Profile{ProfileAllValid, ProfileAllNull, ProfileMixed}

// TestClientRoundTrip is a sanity check that every (variant, profile)
// combination Marshal/Unmarshal-s without error before the benchmark
// numbers are read.
func TestClientRoundTrip(t *testing.T) {
	for _, p := range clientProfiles {
		t.Run("Null/"+p.String(), func(t *testing.T) {
			src := NewClientNull(p)
			data, err := json.Marshal(src)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var dst ClientNull
			if err := json.Unmarshal(data, &dst); err != nil {
				t.Fatalf("unmarshal: %v\nbytes: %s", err, data)
			}
		})
		t.Run("Ptr/"+p.String(), func(t *testing.T) {
			src := NewClientPtr(p)
			data, err := json.Marshal(src)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var dst ClientPtr
			if err := json.Unmarshal(data, &dst); err != nil {
				t.Fatalf("unmarshal: %v\nbytes: %s", err, data)
			}
		})
	}
}

// TestClientJSONSizes prints the marshalled byte sizes for each
// (variant, profile) combination. Run with -v to see the table:
//
//	go test ./bench/ -run TestClientJSONSizes -v
func TestClientJSONSizes(t *testing.T) {
	for _, p := range clientProfiles {
		cnBytes, err := json.Marshal(NewClientNull(p))
		if err != nil {
			t.Fatalf("null marshal: %v", err)
		}
		cpBytes, err := json.Marshal(NewClientPtr(p))
		if err != nil {
			t.Fatalf("ptr marshal: %v", err)
		}
		t.Logf("%-9s  Null=%4d B   Ptr=%4d B", p, len(cnBytes), len(cpBytes))
	}
}

// BenchmarkClient runs Marshal and Unmarshal for both variants across
// the three input profiles. The result table has 12 rows (2 variants ×
// 2 operations × 3 profiles); compare a row against its sibling in the
// other variant to read off the dtx-vs-pointer trade-off for that
// shape.
func BenchmarkClient(b *testing.B) {
	for _, p := range clientProfiles {
		cn := NewClientNull(p)
		cp := NewClientPtr(p)

		cnBytes, err := json.Marshal(cn)
		if err != nil {
			b.Fatalf("null marshal setup: %v", err)
		}
		cpBytes, err := json.Marshal(cp)
		if err != nil {
			b.Fatalf("ptr marshal setup: %v", err)
		}

		b.Run(fmt.Sprintf("Null/Marshal/%s", p), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := json.Marshal(cn); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(fmt.Sprintf("Ptr/Marshal/%s", p), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := json.Marshal(cp); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(fmt.Sprintf("Null/Unmarshal/%s", p), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var dst ClientNull
				if err := json.Unmarshal(cnBytes, &dst); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(fmt.Sprintf("Ptr/Unmarshal/%s", p), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var dst ClientPtr
				if err := json.Unmarshal(cpBytes, &dst); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
