package progress

import (
	"bytes"
	"io"
	"testing"
)

func TestProgressReader(t *testing.T) {
	data := []byte("Hello, World!")
	reader := bytes.NewReader(data)

	totalRead := 0
	wrapped := NewReader(reader, func(n int) {
		totalRead += n
	})

	// Read all data
	buf := make([]byte, 1024)
	n, err := wrapped.Read(buf)

	if err != nil && err != io.EOF {
		t.Fatalf("Unexpected error: %v", err)
	}

	if n != len(data) {
		t.Errorf("Expected to read %d bytes, got %d", len(data), n)
	}

	if totalRead != len(data) {
		t.Errorf("Expected callback totalRead=%d, got %d", len(data), totalRead)
	}

	if !bytes.Equal(buf[:n], data) {
		t.Errorf("Data mismatch: expected '%s', got '%s'", data, buf[:n])
	}
}

func TestProgressReaderChunked(t *testing.T) {
	data := []byte("0123456789")
	reader := bytes.NewReader(data)

	callCount := 0
	totalRead := 0

	wrapped := NewReader(reader, func(n int) {
		callCount++
		totalRead += n
	})

	// Read in 3-byte chunks
	for {
		buf := make([]byte, 3)
		n, err := wrapped.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if n == 0 {
			break
		}
	}

	// Should have called callback 4 times (3+3+3+1)
	if callCount != 4 {
		t.Errorf("Expected 4 callback calls, got %d", callCount)
	}

	if totalRead != len(data) {
		t.Errorf("Expected totalRead=%d, got %d", len(data), totalRead)
	}
}

func TestProgressReaderNoCallback(t *testing.T) {
	data := []byte("Test")
	reader := bytes.NewReader(data)

	// Create without callback (nil callback)
	wrapped := NewReader(reader, nil)

	buf := make([]byte, 1024)
	n, err := wrapped.Read(buf)

	// Should not panic and should read data normally
	if err != nil && err != io.EOF {
		t.Fatalf("Unexpected error: %v", err)
	}

	if n != len(data) {
		t.Errorf("Expected to read %d bytes, got %d", len(data), n)
	}
}

func TestProgressReaderClose(t *testing.T) {
	// Create a closeable reader
	data := []byte("Test")
	reader := io.NopCloser(bytes.NewReader(data))

	wrapped := NewReader(reader, nil)

	// Should support Close
	if closer, ok := wrapped.(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			t.Errorf("Close failed: %v", err)
		}
	} else {
		t.Error("Expected wrapped reader to support io.Closer")
	}
}
