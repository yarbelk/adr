package repo

type ADR interface {
	Get(int) (*ADR, error)
	Add(ADR) error
	Update(...ADR) error
	List() []*ADR
}
