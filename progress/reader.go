package progress

import "io"

// ProgressReader wraps an io.Reader and calls a callback function
// with the number of bytes read on each Read() call.
//
// This allows automatic progress tracking when reading from files,
// network connections, or any other io.Reader.
//
// Example:
//
//	file, _ := os.Open("large-file.dat")
//	fileInfo, _ := file.Stat()
//
//	prog := progress.New(console)
//	task := prog.AddBar("Reading", fileInfo.Size())
//	prog.Start()
//
//	reader := progress.NewReader(file, func(n int) {
//		prog.Advance(task, int64(n))
//	})
//
//	io.Copy(dest, reader)
//	prog.Stop()
type ProgressReader struct {
	reader   io.Reader
	callback func(int)
}

// NewReader creates a new ProgressReader that wraps the given reader.
// The callback is called with the number of bytes read after each successful Read().
//
// The callback is only invoked when n > 0 (bytes were actually read).
// Errors are passed through unchanged.
//
// Example:
//
//	reader := progress.NewReader(file, func(n int) {
//		fmt.Printf("Read %d bytes\n", n)
//	})
func NewReader(r io.Reader, callback func(int)) io.Reader {
	return &ProgressReader{
		reader:   r,
		callback: callback,
	}
}

// Read implements io.Reader.
// Reads from the underlying reader and invokes the callback with the byte count.
func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	if n > 0 && pr.callback != nil {
		pr.callback(n)
	}
	return n, err
}

// Close implements io.Closer if the underlying reader implements it.
// This allows ProgressReader to be used with defer file.Close() patterns.
func (pr *ProgressReader) Close() error {
	if closer, ok := pr.reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Seek implements io.Seeker if the underlying reader implements it.
// This allows ProgressReader to be used with seekable sources.
func (pr *ProgressReader) Seek(offset int64, whence int) (int64, error) {
	if seeker, ok := pr.reader.(io.Seeker); ok {
		return seeker.Seek(offset, whence)
	}
	// If underlying reader doesn't support seeking, return error
	return 0, io.ErrUnexpectedEOF
}
