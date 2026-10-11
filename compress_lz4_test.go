package erofs

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/pierrec/lz4/v4"
)

// An LZ4 stream's final literals can end in zero bytes; they are data, not
// padding, whether or not the pcluster is zero-padded after the stream.
func TestLZ4DecompressKeepsTrailingZeroLiterals(t *testing.T) {
	orig := make([]byte, 600)
	rand.New(rand.NewSource(1)).Read(orig)
	copy(orig[len(orig)-6:], make([]byte, 6))
	// Incompressible, so the stream ends in literals.
	comp := make([]byte, lz4.CompressBlockBound(len(orig)))
	n, err := lz4.CompressBlock(orig, comp, nil)
	if err != nil || n == 0 {
		t.Fatalf("CompressBlock: n=%d err=%v", n, err)
	}
	comp = comp[:n]
	if comp[len(comp)-1] != 0 {
		t.Fatalf("stream does not end in a zero literal")
	}
	for _, pad := range []int{0, 100} {
		src := append(append([]byte{}, comp...), make([]byte, pad)...)
		dst := make([]byte, len(orig))
		got, err := lz4Decompress(src, dst)
		if err != nil {
			t.Fatalf("pad %d: %v", pad, err)
		}
		if !bytes.Equal(dst[:got], orig) {
			t.Fatalf("pad %d: output differs", pad)
		}
	}
}
