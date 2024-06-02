package file_test

import (
	"bytes"
	"sigs.k8s.io/kwok/pkg/utils/file"
	"testing"
)

func TestCompress(t *testing.T) {
	testCases := []struct {
		name         string
		filename     string
		decompressed []byte
		compressed   []byte
	}{
		{
			name:         "Compress gzip data",
			filename:     "test.gz",
			decompressed: []byte("Hello, World!"),
			compressed:   []byte{31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 242, 72, 205, 201, 201, 215, 81, 8, 207, 47, 202, 73, 81, 4, 4, 0, 0, 255, 255, 208, 195, 74, 236, 13, 0, 0, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			w := file.Compress(tc.filename, buf)
			_, err := w.Write(tc.decompressed)
			if err != nil {
				t.Errorf("Compress failed: %v", err)
			}
			err = w.Close()
			if err != nil {
				t.Errorf("Compress failed: %v", err)
			}

			if !bytes.Equal(buf.Bytes(), tc.compressed) {
				t.Errorf("Compress(%q) = %v, want %v", tc.filename, buf.Bytes(), tc.compressed)
			}
		})
	}
}

func TestDecompress(t *testing.T) {
	testCases := []struct {
		name         string
		filename     string
		compressed   []byte
		decompressed []byte
	}{
		{
			name:         "Decompress gzip data",
			filename:     "test.gz",
			compressed:   []byte{31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 242, 72, 205, 201, 201, 215, 81, 8, 207, 47, 202, 73, 81, 4, 4, 0, 0, 255, 255, 208, 195, 74, 236, 13, 0, 0, 0},
			decompressed: []byte("Hello, World!"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := make([]byte, len(tc.decompressed))
			r, err := file.Decompress(tc.filename, bytes.NewReader(tc.decompressed))
			if err != nil {
				t.Errorf("Decompress failed: %v", err)
			}
			_, err = r.Read(buf)
			if err != nil {
				t.Errorf("Decompress failed: %v", err)
			}
			err = r.Close()
			if err != nil {
				return
			}
			if !bytes.Equal(buf, tc.decompressed) {
				t.Errorf("Decompress(%q) = %v, want %v", tc.filename, buf, tc.compressed)
			}
		})
	}
}
