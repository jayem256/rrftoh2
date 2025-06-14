package handler

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"rrftoh2/constants"
	"rrftoh2/file"
	"slices"
	"strconv"
	"strings"
)

// RequestHandler handles incoming HTTP2 requests
type RequestHandler struct {
	Conn        *tls.Conn
	Hostname    string
	DocRoot     string
	Compression bool
	Treshold    int
	ReadSize    int
	Omitted     []string
	GzBuffer    int
}

const head = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <link rel="shortcut icon" href="data:," />
	<style>
	  a {
	    color: blue;	
	  }
	
	  a:hover {
	    color: purple;  
	  }
	</style>
`
const body = `
  </head>
  <body>
    <table style="border-spacing: 5px;">
	  <tr>
	    <td style="border: 1px solid; min-width: 100px;">Name</td>
		<td style="border: 1px solid;">Size</td>
	  </tr>
`

const tail = `    </table>
  </body>
</html>`

// ServeHTTP handle request and response logic
func (s *RequestHandler) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	// Check protocol version.
	if rq.Proto != constants.HTTP2_VERSION {
		// Probably not necessary but nuke anyway if we somehow do end up here.
		s.Conn.Close()
		return
	}

	// Only accept GET method.
	if rq.Method == http.MethodGet {
		if strings.Contains(rq.RequestURI, "..") || strings.Count(rq.RequestURI, "/") > 1 {
			rw.WriteHeader(401)
			return
		}
		// Root path provides file listing.
		if rq.RequestURI == "/" {
			// Compose and send HTML.
			rw.Write([]byte(head +
				"    <title>" +
				s.Hostname +
				"</title>" +
				body +
				s.listFiles(s.DocRoot) +
				tail))
		} else {
			fmt.Println(s.Conn.RemoteAddr().String() + " requesting " + rq.RequestURI)

			filepath := s.DocRoot + rq.RequestURI
			info, err := os.Stat(filepath)
			// Check if requested file exists.
			if errors.Is(err, os.ErrNotExist) || info.IsDir() {
				// File does not exist.
				rw.WriteHeader(404)
			} else {
				flusha := rw.(http.Flusher)
				// Files may contain arbitrary data. Use same content type for all.
				rw.Header().Set("Content-Type", "application/octet-stream")

				// Start reading file in chunks.
				chunkStream := file.StartReading(filepath, s.ReadSize, rq.Context())

				// Enable gzip compression if conditions are met.
				if s.doCompress(filepath, int(info.Size())) {
					// Append content encoding header.
					rw.Header().Set("Content-Encoding", "gzip")
					flusha.Flush()
					// Get gzip stream.
					chunkStream = file.GzipStreamify(chunkStream, s.GzBuffer)
				} else {
					// No compression. We know the file size ahead of time.
					rw.Header().Set("Content-Length", strconv.Itoa(int(info.Size())))
				}

				// Stream chunks.
				for chunk := range chunkStream {
					rw.Write(chunk)
					flusha.Flush()
				}
			}
		}
	} else {
		// Method not allowed.
		rw.WriteHeader(405)
	}
}

// doCompress returns true if given file should be compressed
func (s *RequestHandler) doCompress(file string, size int) bool {
	if !s.Compression || size < s.Treshold {
		return false
	}
	// Get file extension.
	splat := strings.Split(file, ".")
	ext := splat[len(splat)-1]
	// Match extension against list of omitted extensions.
	return !slices.Contains(s.Omitted, ext)
}

// listFiles lists all files within given path
func (s *RequestHandler) listFiles(docroot string) string {
	list := ""
	files, err := os.ReadDir(docroot)
	if err == nil {
		for _, entry := range files {
			if !entry.IsDir() {
				list += "      <tr>\n"

				list += "        <td>"
				list += "<a href=\"/" + entry.Name() + "\" download>"
				list += entry.Name()
				list += "</a></td>\n"

				list += "        <td>"
				info, _ := entry.Info()
				list += sizeUnit(int(info.Size()))
				list += "</td>\n"

				list += "      </tr>\n"
			}
		}
	} else {
		fmt.Println(err.Error())
	}
	return list
}

// sizeUnit returns size in variable unit depending on size
func sizeUnit(byteSize int) string {
	if byteSize >= 1048576 {
		return strconv.Itoa(byteSize/1048576) + " MB"
	} else if byteSize >= 1024 {
		return strconv.Itoa(byteSize/1024) + " KB"
	}
	return strconv.Itoa(byteSize) + " Bytes"
}
