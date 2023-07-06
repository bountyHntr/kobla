package types

type ConsesusProtocol interface {
	Run(block *Block) error
	Validate(block *Block) bool
}
