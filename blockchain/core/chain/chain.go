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
	ErrInvalidParentBlock  = errors.New("invalid parent block")
	ErrInvalidTxSignature  = errors.New("invalid tx signature")
	ErrZeroAddressScam     = errors.New("zero address scam")
	ErrDuplicateCoinbaseTx = errors.New("duplicate coinbase tx")
	ErrNoCoinbaseTx        = errors.New("no coinbase tx")
)

type Config struct {
	DBPath          string
	Consensus       types.ConsesusProtocol
	URL             string
	SyncNode        string
	Miner           bool
	CoinbaseAddress string
	Genesis         bool
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
	log.Info("init database")

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

	bc.comm = newCommunicationManager(cfg.URL, cfg.SyncNode, &bc)

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

func (bc *Blockchain) MineBlock(txs []*types.Transaction, coinbase types.Address) error {

	txs = append(txs, newCoinbaseTx(coinbase))
	if err := validateTxs(txs, coinbase); err != nil {
		return err
	}

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

func (bc *Blockchain) addBlock(block *types.Block) error {

	if err := validateTxs(block.Transactions, block.Coinbase); err != nil {
		return err
	}

	if err := bc.saveNewBlock(block); err != nil {
		return fmt.Errorf("save new block %d: %w", block.Number, err)
	}

	bc.blockSubs.notify(block)
	return nil
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
