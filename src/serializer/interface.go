package serializer

import "gitlab.com/yarbelk/adr/src/adr"

type Unmarshaller interface {
	Unmarshal(*adr.ADR) error
}
type Marshaller interface {
	Marshal(*adr.ADR) error
}

type Interface interface {
	Marshaller
	Unmarshaller
}
