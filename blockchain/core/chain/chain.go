package chain

import (
	"errors"
	"fmt"
	"kobla/blockchain/core/db"
	"kobla/blockchain/core/types"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	ErrInvalidParentBlock = errors.New("invalid parent block")
	ErrInvalidTxSignature = errors.New("invalid tx signature")
)

type Config struct {
	DBPath        string
	Consensus     types.ConsesusProtocol
	SyncNode      string
	Miner         bool
	MinderAddress string
}

type Blockchain struct {
	mu sync.RWMutex

	cons types.ConsesusProtocol

	tail    *types.Block // last block, use getter lastBlock()
	db      *db.Database
	mempool *memoryPool
	comm    *communicationManager

	blockSubs *subscriptionManager[types.Block]
}

func New(cfg *Config) (*Blockchain, error) {
	logCtx := log.WithField("service", "blockchain")

	logCtx.Info("sync blockchain")
	logCtx.Info("init database")

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

	bc.comm = newCommunicationManager(cfg.SyncNode, &bc)

	hash, err := database.LastBlockHash()
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, fmt.Errorf("get last block hash: %w", err)
		}

		if err := bc.addGenesisBlock(); err != nil {
			return nil, fmt.Errorf("genesis block: %w", err)
		}
	} else {
		lastBlock, err := bc.BlockByHash(hash)
		if err != nil {
			return nil, fmt.Errorf("get last block: %w", err)
		}
		bc.tail = lastBlock
	}

	bc.comm.listen()

	return &bc, nil
}

func (bc *Blockchain) MineBlock(txs []*types.Transaction, coinbase types.Address) error {

	if err := validateTxs(txs); err != nil {
		return err
	}

	txs = append(txs, newCoinbaseTx(coinbase))
	newBlock, err := types.NewBlock(bc.cons, txs, bc.lastBlock(), coinbase)
	if err != nil {
		return fmt.Errorf("create new block: %w", err)
	}

	if err = bc.saveNewBlock(newBlock); err != nil {
		return fmt.Errorf("save new block %d: %w", newBlock.Number, err)
	}

	bc.blockSubs.notify(newBlock)
	return nil
}

func validateTxs(txs []*types.Transaction) error {

	for _, tx := range txs {
		ok, err := tx.Sender.Verify(tx.Hash, tx.Signature)
		if err != nil {
			return fmt.Errorf("verify tx signature: %w", err)
		}

		if !ok {
			return fmt.Errorf("verify tx signature: sender: %s, hash: %s: %w",
				tx.Sender.String(), tx.Hash.String(), ErrInvalidTxSignature)
		}
	}

	return nil
}

func newCoinbaseTx(coinbase types.Address) *types.Transaction {
	return types.NewTransaction(types.ZeroAddress, coinbase, types.BlockReward, []byte("coinbase"))
}

func (bc *Blockchain) addGenesisBlock() error {
	bc.tail = &types.Block{Number: -1}
	return bc.MineBlock(nil, types.ZeroAddress)
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

func (bc *Blockchain) lastBlock() *types.Block {
	bc.mu.RLock()
	block := bc.tail
	bc.mu.RUnlock()

	return block
}
