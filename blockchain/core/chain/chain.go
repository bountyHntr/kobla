package chain

import (
	"errors"
	"fmt"
	"kobla/blockchain/core/db"
	"kobla/blockchain/core/types"
	"sync"

	log "github.com/sirupsen/logrus"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

func init() {
	log.SetLevel(log.DebugLevel)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	ErrInvalidParentBlock  = errors.New("invalid parent block")
	ErrOldBlock            = errors.New("old block")
	ErrInvalidTxSignature  = errors.New("invalid tx signature")
	ErrZeroAddressScam     = errors.New("zero address scam")
	ErrDuplicateCoinbaseTx = errors.New("duplicate coinbase tx")
	ErrNoCoinbaseTx        = errors.New("no coinbase tx")
	ErrConsensusBroken     = errors.New("consensus broken")
)

////////////////////////////////////////////////////////////////////////////////////////////////////

type Config struct {
	DBPath    string
	Consensus types.ConsesusProtocol
	Url       string
	Nodes     []string
	Genesis   bool
}

////////////////////////////////////////////////////////////////////////////////////////////////////

type Blockchain struct {
	mu sync.RWMutex

	cons types.ConsesusProtocol

	tail    *types.Block // last block, use getter lastBlock()
	db      *db.Database
	mempool *memoryPool
	comm    *communicationManager

	blockSubs *subscriptionManager[types.Block]
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func New(cfg *Config) (*Blockchain, error) {
	log.WithField("path", cfg.DBPath).Info("init database")

	database, err := db.New(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("create db: %w", err)
	}

	bc := Blockchain{
		cons: cfg.Consensus,

		db:      database,
		mempool: newMempool(),

		blockSubs: newSubscription[types.Block](),
	}

	bc.comm, err = newCommunicationManager(cfg.Url, cfg.Nodes, &bc)
	if err != nil {
		return nil, fmt.Errorf("new communication manager: %w", err)
	}

	hash, err := database.LastBlockHash()
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, fmt.Errorf("get last block hash: %w", err)
		}

		bc.tail = &types.Block{Number: -1}

		if cfg.Genesis {
			if err := bc.addGenesisBlock(); err != nil {
				return nil, fmt.Errorf("genesis block: %w", err)
			}
		}
	} else {
		lastBlock, err := bc.BlockByHash(hash)
		if err != nil {
			return nil, fmt.Errorf("get last block: %w", err)
		}
		bc.tail = lastBlock
	}

	if err := bc.comm.listen(); err != nil {
		return nil, fmt.Errorf("run listener: %w", err)
	}

	return &bc, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func (bc *Blockchain) mineBlock(txs []*types.Transaction, coinbase types.Address) error {

	txs = append(txs, newCoinbaseTx(coinbase))
	if err := validateTxs(txs, coinbase); err != nil {
		return err
	}

	newBlock, err := types.NewBlock(bc.cons, txs, bc.lastBlock(), coinbase)
	if err != nil {
		return fmt.Errorf("create new block: %w", err)
	}

	if err = bc.newBlock(newBlock); err != nil {
		return fmt.Errorf("new block %d: %w", newBlock.Number, err)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func (bc *Blockchain) addBlock(block *types.Block) (bool, error) {
	if block.Number <= bc.lastBlock().Number {
		return true, ErrOldBlock
	}

	if err := validateTxs(block.Transactions, block.Coinbase); err != nil {
		return false, err
	}

	if err := bc.newBlock(block); err != nil {
		return false, fmt.Errorf("new block %d: %w", block.Number, err)
	}

	return true, nil
}

func validateTxs(txs []*types.Transaction, coinbase types.Address) error {
	hasCoinbase := false
	for _, tx := range txs {

		if tx.Sender == types.ZeroAddress {
			if tx.Receiver != coinbase {
				return fmt.Errorf("verify coinbase tx receiver: hash: %s: %w",
					tx.Hash.String(), ErrZeroAddressScam)
			}

			if hasCoinbase {
				return fmt.Errorf("verify number of coinbase txs: hash %s: %w",
					tx.Hash.String(), ErrDuplicateCoinbaseTx)
			}

			hasCoinbase = true
			continue
		}

		ok, err := tx.Sender.Verify(tx.Hash, tx.Signature)
		if err != nil {
			return fmt.Errorf("verify tx signature: %w", err)
		}

		if !ok {
			return fmt.Errorf("verify tx signature: sender: %s, hash: %s: %w",
				tx.Sender.String(), tx.Hash.String(), ErrInvalidTxSignature)
		}
	}

	if !hasCoinbase {
		return ErrNoCoinbaseTx
	}

	return nil
}

func newCoinbaseTx(coinbase types.Address) *types.Transaction {
	return types.NewTransaction(types.ZeroAddress, coinbase, types.BlockReward, []byte("coinbase"))
}

func (bc *Blockchain) addGenesisBlock() error {
	return bc.mineBlock(nil, types.ZeroAddress)
}

func (bc *Blockchain) saveNewBlock(block *types.Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if bc.tail.Hash != block.PrevBlockHash {
		return fmt.Errorf("parent: %d, block: %d: %w",
			bc.tail.Number, block.Number, ErrInvalidParentBlock)
	}

	if err := bc.db.SaveBlock(block); err != nil {
		return err
	}

	bc.tail = block
	return nil
}

func (bc *Blockchain) newBlock(block *types.Block) error {
	if !bc.cons.Validate(block) {
		return ErrConsensusBroken
	}

	if err := bc.saveNewBlock(block); err != nil {
		return fmt.Errorf("save new block %d: %w", block.Number, err)
	}

	bc.releaseMempool(block)
	bc.blockSubs.notify(block)
	bc.comm.broadcast(block)

	return nil
}

func (bc *Blockchain) releaseMempool(b *types.Block) {
	for _, tx := range b.Transactions {
		bc.mempool.remove(tx.Hash)
	}
}

func (bc *Blockchain) lastBlock() *types.Block {
	bc.mu.RLock()
	block := bc.tail
	bc.mu.RUnlock()

	return block
}

////////////////////////////////////////////////////////////////////////////////////////////////////
