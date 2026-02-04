package progress

import "io"

// ProgressWriter wraps an io.Writer and calls a callback function
// with the number of bytes written on each Write() call.
//
// This allows automatic progress tracking when writing to files,
// network connections, or any other io.Writer.
//
// Example:
//
//	prog := progress.New(console)
//	task := prog.AddBar("Writing", fileSize)
//	prog.Start()
//
//	writer := progress.NewWriter(file, func(n int) {
//		prog.Advance(task, int64(n))
//	})
//
//	io.Copy(writer, source)
//	prog.Stop()
type ProgressWriter struct {
	writer   io.Writer
	callback func(int)
}

// NewWriter creates a new ProgressWriter that wraps the given writer.
// The callback is called with the number of bytes written after each successful Write().
//
// The callback is only invoked when n > 0 (bytes were actually written).
// Errors are passed through unchanged.
//
// Example:
//
//	writer := progress.NewWriter(file, func(n int) {
//		fmt.Printf("Wrote %d bytes\n", n)
//	})
func NewWriter(w io.Writer, callback func(int)) io.Writer {
	return &ProgressWriter{
		writer:   w,
		callback: callback,
	}
}

// Write implements io.Writer.
// Writes to the underlying writer and invokes the callback with the byte count.
func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.writer.Write(p)
	if n > 0 && pw.callback != nil {
		pw.callback(n)
	}
	return n, err
}

// Close implements io.Closer if the underlying writer implements it.
// This allows ProgressWriter to be used with defer file.Close() patterns.
func (pw *ProgressWriter) Close() error {
	if closer, ok := pw.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Seek implements io.Seeker if the underlying writer implements it.
// This allows ProgressWriter to be used with seekable destinations.
func (pw *ProgressWriter) Seek(offset int64, whence int) (int64, error) {
	if seeker, ok := pw.writer.(io.Seeker); ok {
		return seeker.Seek(offset, whence)
	}
	// If underlying writer doesn't support seeking, return error
	return 0, io.ErrUnexpectedEOF
}
