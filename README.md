# RRFToH2

## Reasonably Rapid File Transfer Over HTTP/2
Simple Golang HTTP/2 service which enforces the use of `h2` and allows listing and downloading files.
There's nothing particularly performant streaming data over single HTTP/2 stream, but since all major 
browsers support using `gzip` encoding, this service makes use of it and compresses files on the fly.
Over slower connection speeds it can significantly decrease download times depending on type of data.

## Usage
This implementation enforces clients use `h2` and thus always requires use of TLS. There's example 
certificate included for quick testing, but of course recommendation is to provide your own. That 
and every other setting passed to server must be configured in `config.json` which contains following:

**Server** object containing networking related settings both for TCP and TLS layers. Additionally window 
size sets initial flow control window of Golang's HTTP/2 server.
```
"server": {
    "hostName": "Test",
    "bindAddr": "0.0.0.0",
    "port": 443,
    "cert": "./server.crt",
    "key": "./server.key",
    "windowSize": 104857600
}
```

**File** object contains root path (`docRoot`) from which files are shared as well as size of buffered reads 
in bytes. Adjusting `bufferedRead` could affect read performance depending on storage medium.
```
"file": {
    "docRoot": "./",
    "bufferedRead": 262144
}
```

**Compression** object controls compression usage. You can enable or disable gzip compression for files 
and set treshold (in bytes) below which files are not eligible for compression. Treshold is set in
`compressionTreshold` field. You can also list file extensions which will not be compressed during 
download. Files with specific extensions can be assumed to already contain compressed data. Gzip 
stream is buffered before sending. You can control size of buffer (in bytes) by adjusting 
`gzipBuffer` property.
```
"compression": {
    "enableGzip": true,
    "compressionTreshold": 10485760,
    "gzipBuffer": 65536,
    "omitExtensions": [
        "zip", "jpg", "mkv", "png", "tar", "gz", 
        "rar", "mp4", "avi", "gif", "mp3", "webp"
    ]
}
```