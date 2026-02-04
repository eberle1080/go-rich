package progress

import (
	"bytes"
	"testing"
)

func TestProgressWriter(t *testing.T) {
	var buf bytes.Buffer

	totalWritten := 0
	wrapped := NewWriter(&buf, func(n int) {
		totalWritten += n
	})

	data := []byte("Hello, World!")
	n, err := wrapped.Write(data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if n != len(data) {
		t.Errorf("Expected to write %d bytes, got %d", len(data), n)
	}

	if totalWritten != len(data) {
		t.Errorf("Expected callback totalWritten=%d, got %d", len(data), totalWritten)
	}

	if buf.String() != string(data) {
		t.Errorf("Data mismatch: expected '%s', got '%s'", data, buf.String())
	}
}

func TestProgressWriterMultiple(t *testing.T) {
	var buf bytes.Buffer

	callCount := 0
	totalWritten := 0

	wrapped := NewWriter(&buf, func(n int) {
		callCount++
		totalWritten += n
	})

	// Write multiple times
	writes := [][]byte{
		[]byte("Hello"),
		[]byte(", "),
		[]byte("World"),
		[]byte("!"),
	}

	for _, data := range writes {
		n, err := wrapped.Write(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if n != len(data) {
			t.Errorf("Expected to write %d bytes, got %d", len(data), n)
		}
	}

	if callCount != len(writes) {
		t.Errorf("Expected %d callback calls, got %d", len(writes), callCount)
	}

	expectedTotal := 0
	for _, data := range writes {
		expectedTotal += len(data)
	}

	if totalWritten != expectedTotal {
		t.Errorf("Expected totalWritten=%d, got %d", expectedTotal, totalWritten)
	}

	if buf.String() != "Hello, World!" {
		t.Errorf("Final data mismatch: expected 'Hello, World!', got '%s'", buf.String())
	}
}

func TestProgressWriterNoCallback(t *testing.T) {
	var buf bytes.Buffer

	// Create without callback (nil callback)
	wrapped := NewWriter(&buf, nil)

	data := []byte("Test")
	n, err := wrapped.Write(data)
	// Should not panic and should write data normally
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if n != len(data) {
		t.Errorf("Expected to write %d bytes, got %d", len(data), n)
	}

	if buf.String() != "Test" {
		t.Errorf("Data mismatch: expected 'Test', got '%s'", buf.String())
	}
}
