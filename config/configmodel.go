package config

// Config root type
type Config struct {
	Server      ServerType      `json:"server"`
	Client      ClientType      `json:"client"`
	File        FileType        `json:"file"`
	Compression CompressionType `json:"compression"`
}

// ServerType contains server specific settings
type ServerType struct {
	HostName   string `json:"hostName"`
	BindAddr   string `json:"bindAddr"`
	Port       int    `json:"port"`
	Cert       string `json:"cert"`
	Key        string `json:"key"`
	WindowSize int    `json:"windowSize"`
}

// ClientType contains client certificate pool
type ClientType struct {
	CertPool []string `json:"certPool"`
}

// FileType contains file I/O related settings
type FileType struct {
	DocRoot      string `json:"docRoot"`
	BufferedRead int    `json:"bufferedRead"`
}

// CompressionType contains file compression related settings
type CompressionType struct {
	EnableGzip          bool     `json:"enableGzip"`
	CompressionTreshold int      `json:"compressionTreshold"`
	GzipBuffer          int      `json:"gzipBuffer"`
	OmitExtensions      []string `json:"omitExtensions"`
}
