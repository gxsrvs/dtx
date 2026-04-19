package bench

import (
	"runtime"
	"testing"
	"unsafe"
)

// TestClientHeapSize measures the steady-state heap footprint of
// holding N client instances in memory, for each (variant, profile)
// combination. Run with:
//
//	go test ./bench/ -run TestClientHeapSize -v
//
// The reported "per-instance" number is total heap delta divided by
// N; it captures both the slice element itself and every heap
// allocation performed while constructing the instance (inner pointers,
// pointed-to strings, etc.). The "sizeof" column shows what a single
// element of the slice occupies inline — anything beyond that lives
// on the heap behind a pointer.
func TestClientHeapSize(t *testing.T) {
	const N = 200_000

	t.Logf("sizeof ClientNull = %d B   sizeof ClientPtr = %d B",
		unsafe.Sizeof(ClientNull{}), unsafe.Sizeof(ClientPtr{}))

	for _, p := range clientProfiles {
		measureNull := func() {
			runtime.GC()
			var before runtime.MemStats
			runtime.ReadMemStats(&before)

			arr := make([]ClientNull, N)
			for i := range N {
				arr[i] = NewClientNull(p)
			}

			runtime.GC()
			var after runtime.MemStats
			runtime.ReadMemStats(&after)
			runtime.KeepAlive(arr)

			delta := int64(after.HeapAlloc) - int64(before.HeapAlloc)
			t.Logf("Null/%-9s  total=%11d B  per-instance=%7.1f B",
				p, delta, float64(delta)/float64(N))
		}
		measurePtr := func() {
			runtime.GC()
			var before runtime.MemStats
			runtime.ReadMemStats(&before)

			arr := make([]ClientPtr, N)
			for i := range N {
				arr[i] = NewClientPtr(p)
			}

			runtime.GC()
			var after runtime.MemStats
			runtime.ReadMemStats(&after)
			runtime.KeepAlive(arr)

			delta := int64(after.HeapAlloc) - int64(before.HeapAlloc)
			t.Logf("Ptr/%-9s   total=%11d B  per-instance=%7.1f B",
				p, delta, float64(delta)/float64(N))
		}

		measureNull()
		measurePtr()
	}
}
