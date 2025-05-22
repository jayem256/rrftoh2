package file

// GzipBufferedWriter performs buffered writes of compression to output channel
type GzipBufferedWriter struct {
	out        chan []byte
	buf        []byte
	bufferSize int
}

// NewWriter sets up new writer with backing buffer and output channel
func (g *GzipBufferedWriter) NewWriter(out chan []byte) *GzipBufferedWriter {
	g.out = out
	g.buf = make([]byte, 0, g.bufferSize)

	return g
}

// Write takes in compressed chunks and either buffers or writes to channel
func (g *GzipBufferedWriter) Write(p []byte) (int, error) {
	// Check if buffer can take another slice.
	if len(g.buf)+len(p) > g.bufferSize {
		// Stream current buffer contents.
		g.out <- g.buf
		// New empty buffer.
		g.buf = make([]byte, 0, g.bufferSize)
		// Append the slice which didn't fit.
		g.buf = append(g.buf, p...)
	} else {
		// Append slice to buffer.
		g.buf = append(g.buf, p...)
	}
	return len(p), nil
}

// Close flushes remaining data and closes output channel
func (g *GzipBufferedWriter) Close() {
	// Flush any remaining bytes.
	if len(g.buf) > 0 {
		g.out <- g.buf
	}
	// Close channel.
	close(g.out)
}
