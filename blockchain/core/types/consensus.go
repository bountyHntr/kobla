package types

type ConsesusProtocol interface {
	Run(block *Block, meta any) error
	Validate(block *Block, meta any) bool
	NodesAreFixed() bool
}
