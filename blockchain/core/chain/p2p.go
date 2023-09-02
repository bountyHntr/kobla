package chain

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"kobla/blockchain/core/pb"
	"kobla/blockchain/core/types"
	"net"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	ErrReqPartiallySent = errors.New("request partially sent")
	ErrInvalidResponse  = errors.New("invalid response")
	ErrInvalidRequest   = errors.New("invalid request")
	ErrInternalError    = errors.New("internal error")
)

const (
	defaultAddress = "localhost:8888"
	netProtocol    = "tcp"
)

//go:generate stringer -type=Command
type Command int

const (
	commandResponse Command = iota
	commandSync
	commandGetBlock
	commandNewBlock
	commandNewTx
)

////////////////////////////////////////////////////////////////////////////////////////////////////

type communicationManager struct {
	mu         sync.RWMutex
	url        string
	dialer     net.Dialer
	knownNodes map[string]struct{}
	bc         *Blockchain
}

func newCommunicationManager(url string, nodes []string, bc *Blockchain) (*communicationManager, error) {
	knownNodes := make(map[string]struct{})
	for _, node := range nodes {
		knownNodes[node] = struct{}{}
	}

	if url == "" {
		url = defaultAddress
	}

	address, err := net.ResolveTCPAddr(netProtocol, url)
	if err != nil {
		return nil, fmt.Errorf("resolve TCP address: %w", err)
	}

	return &communicationManager{
		url:        url,
		dialer:     net.Dialer{LocalAddr: address},
		knownNodes: knownNodes,
		bc:         bc,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func (cm *communicationManager) listen() error {
	ln, err := net.Listen(netProtocol, cm.url)
	if err != nil {
		return err
	}

	go func() {
		const syncInterval = 5 * time.Minute

		ticker := time.NewTicker(syncInterval)
		defer ticker.Stop()

		for ; true; <-ticker.C {
			if err := cm.sendSync(); err != nil {
				log.WithError(err).Error("sync")
			}
		}
	}()

	go func() {
		defer ln.Close()

		log.WithField("url", cm.url).Info("listen")
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.WithError(err).Panic("accept connection")
			}

			go cm.handleConnection(conn)
		}
	}()

	return nil
}

func (cm *communicationManager) handleConnection(conn net.Conn) {
	defer conn.Close()

	remote := conn.RemoteAddr().String()

	logCtx := log.WithField("url", remote)
	logCtx.Debug("accept connection")

	request, err := io.ReadAll(conn)
	if err != nil {
		logCtx.WithError(err).
			WithField("node", remote).
			Error("read request")
		return
	}

	command, request := parseCommand(request)
	logCtx.WithField("command", command.String()).Debug("got request")

	switch command {
	case commandSync:
		err = cm.handleSync(conn, request)
	case commandGetBlock:
		err = cm.handleGetBlock(conn, request)
	case commandNewBlock:
		err = cm.handleNewBlock(request, remote)
	case commandNewTx:
		err = cm.handleNewTx(request)
	}

	if err != nil {
		logCtx.WithError(err).
			WithField("command", command.String()).
			Error("execute command")
	}
}

func (cm *communicationManager) handleSync(conn net.Conn, requestData []byte) error {
	var request pb.ChainStatus
	if err := proto.Unmarshal(requestData, &request); err != nil {
		writeError(conn, ErrInvalidRequest)
		return nil
	}

	cm.addNewNode(request.AddressFrom)
	for _, address := range request.KnownAddresses {
		cm.addNewNode(address)
	}

	responseData, err := proto.Marshal(cm.chainStatus())
	if err != nil {
		writeError(conn, ErrInternalError)
		return fmt.Errorf("marshal: %w", err)
	}

	if err := writeResponse(conn, responseData); err != nil {
		return fmt.Errorf("send response: %w", err)
	}

	return nil
}

func (cm *communicationManager) handleGetBlock(conn net.Conn, request []byte) (err error) {
	if len(request) != 8 {
		return writeError(conn, ErrInvalidRequest)
	}
	blockNumber := int64(binary.BigEndian.Uint64(request))

	block, err := cm.bc.BlockByNumber(blockNumber)
	if err != nil {
		return writeError(conn, fmt.Errorf("get block by number: %w", err))
	}

	data, err := block.Serialize()
	if err != nil {
		writeError(conn, ErrInternalError)
		return fmt.Errorf("serialize block: %w", err)
	}

	return writeResponse(conn, data)
}

func (cm *communicationManager) handleNewBlock(data []byte, meta any) error {
	block, err := types.DeserializeBlock(data)
	if err != nil {
		return fmt.Errorf("deserialize block: %w", err)
	}

	if err := cm.bc.addBlock(block, meta); err != nil {
		return fmt.Errorf("add block: %w", err)
	}

	return nil
}

func (cm *communicationManager) handleNewTx(data []byte) error {
	tx, err := types.DeserializeTx(data)
	if err != nil {
		return fmt.Errorf("deserialize tx: %w", err)
	}

	if ok := cm.bc.mempool.add(tx); ok {
		cm.broadcast(tx)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func (cm *communicationManager) sendSync() error {
	syncNode := cm.randomNode()
	if syncNode == "" {
		log.Debug("skip sync")
		return nil
	}

	log.WithField("sync_node", syncNode).Debug("send sync")

	conn, err := cm.newConnection(syncNode)
	if err != nil {
		cm.removeNode(syncNode)
		return err
	}
	defer conn.Close()

	requestData, err := proto.Marshal(cm.chainStatus())
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := write(conn, commandSync, requestData); err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	responseData, err := read(conn)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	var response pb.ChainStatus
	if err := proto.Unmarshal(responseData, &response); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	if err := cm.sync(syncNode, response.Height); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	return nil
}

func (cm *communicationManager) sendGetBlock(syncNode string, blockNumber int64) (*types.Block, error) {
	log.WithField("sync_node", syncNode).
		WithField("block_number", blockNumber).
		Debug("send get block")

	conn, err := cm.newConnection(syncNode)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(blockNumber))

	if err := write(conn, commandGetBlock, data); err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	data, err = read(conn)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	block, err := types.DeserializeBlock(data)
	if err != nil {
		return nil, fmt.Errorf("deserialize block: %w", err)
	}

	return block, nil
}

func (cm *communicationManager) sync(syncNode string, lastBlockNumber int64) error {
	localLastBlock := cm.bc.lastBlock().Number

	for blockNumber := localLastBlock + 1; blockNumber <= lastBlockNumber; blockNumber++ {

		block, err := cm.sendGetBlock(syncNode, blockNumber)
		if err != nil {
			return fmt.Errorf("get block: %w", err)
		}

		if err := cm.bc.addBlock(block, syncNode); err != nil {
			return fmt.Errorf("add block: %w", err)
		}
	}

	return nil
}

func (cm *communicationManager) chainStatus() *pb.ChainStatus {
	return &pb.ChainStatus{
		Height:         cm.bc.lastBlock().Number,
		AddressFrom:    cm.url,
		KnownAddresses: cm.copyNodes(),
	}
}

func (cm *communicationManager) broadcast(msg types.Serializable) {
	data, _ := msg.Serialize()

	var command Command
	switch msg.(type) {
	case *types.Transaction:
		command = commandNewTx
	case *types.Block:
		command = commandNewBlock
	}

	for node := range cm.knownNodes {

		conn, err := cm.newConnection(node)
		if err != nil {
			log.WithField("node", node).WithError(err).Error("new connection")
			continue
		}

		if err := write(conn, command, data); err != nil {
			log.WithField("node", node).WithError(err).Error("send data")
		}

		conn.Close()
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func (cm *communicationManager) addNewNode(node string) {
	if cm.bc.cons.NodesAreFixed() {
		return
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, ok := cm.knownNodes[node]; ok {
		return
	}

	log.WithField("node", node).Debug("add new node")
	cm.knownNodes[node] = struct{}{}
}

func (cm *communicationManager) removeNode(node string) {
	if cm.bc.cons.NodesAreFixed() {
		return
	}

	cm.mu.Lock()
	log.WithField("node", node).Debug("remove node")
	delete(cm.knownNodes, node)
	cm.mu.Unlock()
}

func (cm *communicationManager) randomNode() (node string) {
	cm.mu.RLock()
	for n := range cm.knownNodes {
		node = n
		break
	}
	cm.mu.RUnlock()

	return
}

func (cm *communicationManager) copyNodes() (knownNodes []string) {
	cm.mu.RLock()
	for node := range cm.knownNodes {
		knownNodes = append(knownNodes, node)
	}
	cm.mu.RUnlock()

	return
}

func (cm *communicationManager) newConnection(node string) (net.Conn, error) {
	conn, err := cm.dialer.Dial(netProtocol, node)
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", node, err)
	}
	return conn, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////

func write(conn net.Conn, command Command, data []byte) (err error) {
	defer func() {
		if err == nil {
			err = conn.(*net.TCPConn).CloseWrite()
		}
	}()

	payload := addCommand(command, data)

	n, err := conn.Write(payload)
	if err != nil {
		return err
	}
	if n != len(payload) {
		return ErrReqPartiallySent
	}

	return nil
}

func writeResponse(conn net.Conn, data []byte) error {
	return write(conn, commandResponse, data)
}

func writeError(conn net.Conn, err error) error {
	payload := fmt.Sprintf("error: %s", err)
	return writeResponse(conn, []byte(payload))
}

func read(conn net.Conn) ([]byte, error) {
	data, err := io.ReadAll(conn)
	if err != nil {
		return nil, err
	}

	command, data := parseCommand(data)
	if command != commandResponse {
		return nil, ErrInvalidResponse
	}

	if strings.Contains(string(data), "error") {
		return nil, errors.New(string(data))
	}

	return data, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func parseCommand(request []byte) (Command, []byte) {
	return Command(request[0]), request[1:]
}

func addCommand(command Command, request []byte) []byte {
	return append([]byte{byte(command)}, request...)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
