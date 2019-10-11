package serializer

import (
	"io"

	"github.com/BurntSushi/toml"
	"gitlab.com/yarbelk/adr/src/adr"
)

// BufferUnmarshaller is a simple wrapper arround a reader
type BufferUnmarshaller struct {
	io.Reader
}

// NewUnmarshal is what it says
func NewUnmarshal(r io.Reader) BufferUnmarshaller {
	return BufferUnmarshaller{r}
}

// Unmarshal into the supplied ADR reference
func (u BufferUnmarshaller) Unmarshal(adr *adr.ADR) error {
	_, err := toml.DecodeReader(u.Reader, adr)
	return err
}
