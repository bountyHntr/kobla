package types

type Copier[T any] interface {
	Copy() *T
}

type Serializable interface {
	Serialize() ([]byte, error)
}

type Void struct{}

func (v Void) Copy() *Void {
	return &v
}

var _ Copier[Void] = Void{}
