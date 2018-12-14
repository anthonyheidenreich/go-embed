package embed

import "io"

// Embedder defines basic functionality of an Embedder module
type Embedder interface {
	// EmbedFile in the Go code
	EmbedFile(fileName string, contents []byte) (err error)
	// Finalize writes the structure to the file
	Finalize(fileDescriptor io.Writer) (err error)
}
