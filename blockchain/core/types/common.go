package types

type Serializable interface {
	Serialize() ([]byte, error)
}
