package erofs

import (
	"bytes"
	"testing"

	"github.com/W-Floyd/erofs/internal/ondisk"
)

// A non-interlaced plain extent that starts mid-lcluster (clusterofs > 0, as
// mkfs.erofs writes a file's incompressible tail) stores its bytes from the
// start of the pcluster, not shifted by clusterofs.
func TestDecodePlainExtentStartsAtPclusterStart(t *testing.T) {
	const bs = 4096
	want := make([]byte, 3518)
	for i := range want {
		want[i] = byte(i*7 + 1)
	}
	block := make([]byte, bs)
	copy(block, want)

	cf := &compressedFile{
		fsys:     &FS{r: bytes.NewReader(block), blockSize: bs, blockSzBits: 12},
		lclustSz: bs,
	}
	const clusterOfs = 928
	if err := cf.decodeExtent(ondisk.LClusterTypePlain, 0xFF, 0, bs, 5*bs+clusterOfs, 5*bs+clusterOfs+int64(len(want)), clusterOfs); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(cf.cacheData, want) {
		t.Errorf("plain extent bytes differ (first %x, want %x)", cf.cacheData[:8], want[:8])
	}
}
