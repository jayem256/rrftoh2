package file

import (
	"bufio"
	"context"
	"os"
	"rrftoh2/constants"
)

// StartReading creates buffered reader goroutine and returns chunk stream
func StartReading(path string, readsize int, h2context context.Context) chan []byte {
	file, err := os.Open(path)

	stream := make(chan []byte, constants.MAX_QUEUE)

	if err != nil {
		close(stream)
		return stream
	}

	// Goroutine for reading the file in chunks.
	go func(out chan []byte, ctx context.Context, file *os.File) {
		// Buffered reader.
		br := bufio.NewReaderSize(file, readsize)
		// Read file contents in chunks.
		for {
			// Attempt to read up to readsize number of bytes at once.
			buf := make([]byte, readsize)
			read, err := br.Read(buf)
			select {
			// Connection terminated.
			case <-ctx.Done():
				// Stop streaming file contents.
				close(stream)
				file.Close()
				return
			// Send raw chunk.
			default:
				if read > 0 {
					stream <- buf[:read]
				}
				// EOF.
				if read <= 0 || err != nil {
					close(stream)
					file.Close()
					return
				}
			}
		}
	}(stream, h2context, file)

	return stream
}
