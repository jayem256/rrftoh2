package file

import (
	"compress/gzip"
	"rrftoh2/constants"
)

// GzipStreamify takes channel feed of raw chunks and returns channel for compressed chunks
func GzipStreamify(raw <-chan []byte, gzipBufferSize int) chan []byte {
	compressed := make(chan []byte, constants.MAX_QUEUE)
	gw := (&GzipBufferedWriter{bufferSize: gzipBufferSize}).NewWriter(compressed)
	gz, _ := gzip.NewWriterLevel(gw, gzip.BestSpeed)

	// Goroutine for compression stream.
	go func(chunkStream <-chan []byte, gs *GzipBufferedWriter, gz *gzip.Writer) {
		// Get raw chunks.
		for chunk := range chunkStream {
			// Compress it.
			gz.Write(chunk)
		}
		gz.Close()
		gs.Close()
	}(raw, gw, gz)

	return compressed
}
