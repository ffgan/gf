//go:build linux

package cli

import "testing"

func TestUtsStringBothElementTypes(t *testing.T) {
	var signed [65]int8
	var unsigned [65]uint8
	for i, c := range []byte("Linux") {
		signed[i] = int8(c)
		unsigned[i] = c
	}
	// A byte above 127 must keep its bit pattern on both element types.
	var hi byte = 0xE9
	signed[5] = int8(hi)
	unsigned[5] = hi

	want := string(append([]byte("Linux"), 0xE9))
	if got := utsString(signed); got != want {
		t.Fatalf("int8: got %q", got)
	}
	if got := utsString(unsigned); got != want {
		t.Fatalf("uint8: got %q", got)
	}
}
