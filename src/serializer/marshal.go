package serializer

import (
	"io"

	"github.com/BurntSushi/toml"
	"gitlab.com/yarbelk/adr/src/adr"
)

// BufferMarshaller is a simple wrapper arround a reader
type BufferMarshaller struct {
	io.Writer
}

// NewMarshal is what it says
func NewMarshal(r io.Writer) Marshaller {
	return BufferMarshaller{r}
}

// Marshal into the supplied ADR reference
func (u BufferMarshaller) Marshal(adr *adr.ADR) error {
	return toml.NewEncoder(u.Writer).Encode(adr)
}
