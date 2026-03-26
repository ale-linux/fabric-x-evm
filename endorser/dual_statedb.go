/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: LGPL-3.0-or-later
*/

package endorser

import (
	"github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/stateless"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie/utils"
	"github.com/holiman/uint256"
	"github.com/hyperledger/fabric-x-sdk/blocks"
	"github.com/hyperledger/fabric-x-sdk/state"
)

// DualStateDB implements the vm.StateDB interface by delegating all calls
// to both a go-ethereum StateDB and an endorser SnapshotDB.
// This allows both state implementations to be kept in sync during execution.
type DualStateDB struct {
	ethStateDB *ethstate.StateDB
	snapshotDB *SnapshotDB
}

// NewDualStateDB creates a new DualStateDB that wraps both state implementations.
// The constructor takes concrete types (not interfaces) so that callers can
// access non-interface methods on both implementations.
func NewDualStateDB(ethStateDB *ethstate.StateDB, snapshotDB *SnapshotDB) *DualStateDB {
	return &DualStateDB{
		ethStateDB: ethStateDB,
		snapshotDB: snapshotDB,
	}
}

// EthStateDB returns the underlying go-ethereum StateDB for accessing
// non-interface methods.
func (d *DualStateDB) EthStateDB() *ethstate.StateDB {
	return d.ethStateDB
}

// SnapshotDB returns the underlying endorser SnapshotDB for accessing
// non-interface methods.
func (d *DualStateDB) SnapshotDB() *SnapshotDB {
	return d.snapshotDB
}

// CreateAccount creates an account in both state implementations.
func (d *DualStateDB) CreateAccount(addr common.Address) {
	d.ethStateDB.CreateAccount(addr)
	d.snapshotDB.CreateAccount(addr)
}

// CreateContract creates a contract account in both state implementations.
func (d *DualStateDB) CreateContract(addr common.Address) {
	d.ethStateDB.CreateContract(addr)
	d.snapshotDB.CreateContract(addr)
}

// SubBalance subtracts balance from an account in both state implementations.
// Returns the previous balance from the eth StateDB.
func (d *DualStateDB) SubBalance(addr common.Address, amount *uint256.Int, reason tracing.BalanceChangeReason) uint256.Int {
	prev := d.ethStateDB.SubBalance(addr, amount, reason)
	d.snapshotDB.SubBalance(addr, amount, reason)
	return prev
}

// AddBalance adds balance to an account in both state implementations.
// Returns the previous balance from the eth StateDB.
func (d *DualStateDB) AddBalance(addr common.Address, amount *uint256.Int, reason tracing.BalanceChangeReason) uint256.Int {
	prev := d.ethStateDB.AddBalance(addr, amount, reason)
	d.snapshotDB.AddBalance(addr, amount, reason)
	return prev
}

// GetBalance returns the balance from the SnapshotDB.
func (d *DualStateDB) GetBalance(addr common.Address) *uint256.Int {
	return d.snapshotDB.GetBalance(addr)
}

// GetNonce returns the nonce from the SnapshotDB.
func (d *DualStateDB) GetNonce(addr common.Address) uint64 {
	return d.snapshotDB.GetNonce(addr)
}

// SetNonce sets the nonce in both state implementations.
func (d *DualStateDB) SetNonce(addr common.Address, nonce uint64, reason tracing.NonceChangeReason) {
	d.ethStateDB.SetNonce(addr, nonce, reason)
	d.snapshotDB.SetNonce(addr, nonce, reason)
}

// GetCodeHash returns the code hash from the SnapshotDB.
func (d *DualStateDB) GetCodeHash(addr common.Address) common.Hash {
	return d.snapshotDB.GetCodeHash(addr)
}

// GetCode returns the code from the SnapshotDB.
func (d *DualStateDB) GetCode(addr common.Address) []byte {
	return d.snapshotDB.GetCode(addr)
}

// SetCode sets the code in both state implementations.
// Returns the previous code from the eth StateDB.
func (d *DualStateDB) SetCode(addr common.Address, code []byte, reason tracing.CodeChangeReason) []byte {
	prev := d.ethStateDB.SetCode(addr, code, reason)
	d.snapshotDB.SetCode(addr, code, reason)
	return prev
}

// GetCodeSize returns the code size from the SnapshotDB.
func (d *DualStateDB) GetCodeSize(addr common.Address) int {
	return d.snapshotDB.GetCodeSize(addr)
}

// AddRefund adds a gas refund in both state implementations.
func (d *DualStateDB) AddRefund(gas uint64) {
	d.ethStateDB.AddRefund(gas)
	d.snapshotDB.AddRefund(gas)
}

// SubRefund subtracts a gas refund in both state implementations.
func (d *DualStateDB) SubRefund(gas uint64) {
	d.ethStateDB.SubRefund(gas)
	d.snapshotDB.SubRefund(gas)
}

// GetRefund returns the refund counter from the SnapshotDB.
func (d *DualStateDB) GetRefund() uint64 {
	return d.snapshotDB.GetRefund()
}

// GetStateAndCommittedState returns both current and committed state from the SnapshotDB.
func (d *DualStateDB) GetStateAndCommittedState(addr common.Address, hash common.Hash) (common.Hash, common.Hash) {
	return d.snapshotDB.GetStateAndCommittedState(addr, hash)
}

// GetState returns the state from the SnapshotDB.
func (d *DualStateDB) GetState(addr common.Address, hash common.Hash) common.Hash {
	return d.snapshotDB.GetState(addr, hash)
}

// SetState sets the state in both state implementations.
// Returns the previous state from the eth StateDB.
func (d *DualStateDB) SetState(addr common.Address, key common.Hash, value common.Hash) common.Hash {
	prev := d.ethStateDB.SetState(addr, key, value)
	d.snapshotDB.SetState(addr, key, value)
	return prev
}

// GetStorageRoot returns the storage root from the SnapshotDB.
func (d *DualStateDB) GetStorageRoot(addr common.Address) common.Hash {
	return d.snapshotDB.GetStorageRoot(addr)
}

// GetTransientState returns the transient state from the SnapshotDB.
func (d *DualStateDB) GetTransientState(addr common.Address, key common.Hash) common.Hash {
	return d.snapshotDB.GetTransientState(addr, key)
}

// SetTransientState sets the transient state in both state implementations.
func (d *DualStateDB) SetTransientState(addr common.Address, key, value common.Hash) {
	d.ethStateDB.SetTransientState(addr, key, value)
	d.snapshotDB.SetTransientState(addr, key, value)
}

// SelfDestruct performs self-destruct in both state implementations.
// Returns the balance from the eth StateDB.
func (d *DualStateDB) SelfDestruct(addr common.Address) uint256.Int {
	balance := d.ethStateDB.SelfDestruct(addr)
	d.snapshotDB.SelfDestruct(addr)
	return balance
}

// HasSelfDestructed checks if an account has self-destructed in the SnapshotDB.
func (d *DualStateDB) HasSelfDestructed(addr common.Address) bool {
	return d.snapshotDB.HasSelfDestructed(addr)
}

// SelfDestruct6780 performs EIP-6780 self-destruct in both state implementations.
// Returns the balance and destruction status from the eth StateDB.
func (d *DualStateDB) SelfDestruct6780(addr common.Address) (uint256.Int, bool) {
	balance, destructed := d.ethStateDB.SelfDestruct6780(addr)
	d.snapshotDB.SelfDestruct6780(addr)
	return balance, destructed
}

// Exist checks if an account exists in the SnapshotDB.
func (d *DualStateDB) Exist(addr common.Address) bool {
	return d.snapshotDB.Exist(addr)
}

// Empty checks if an account is empty in the SnapshotDB.
func (d *DualStateDB) Empty(addr common.Address) bool {
	return d.snapshotDB.Empty(addr)
}

// AddressInAccessList checks if an address is in the access list in the SnapshotDB.
func (d *DualStateDB) AddressInAccessList(addr common.Address) bool {
	return d.snapshotDB.AddressInAccessList(addr)
}

// SlotInAccessList checks if a slot is in the access list in the SnapshotDB.
func (d *DualStateDB) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	return d.snapshotDB.SlotInAccessList(addr, slot)
}

// AddAddressToAccessList adds an address to the access list in both state implementations.
func (d *DualStateDB) AddAddressToAccessList(addr common.Address) {
	d.ethStateDB.AddAddressToAccessList(addr)
	d.snapshotDB.AddAddressToAccessList(addr)
}

// AddSlotToAccessList adds a slot to the access list in both state implementations.
func (d *DualStateDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	d.ethStateDB.AddSlotToAccessList(addr, slot)
	d.snapshotDB.AddSlotToAccessList(addr, slot)
}

// PointCache returns the point cache from the SnapshotDB.
func (d *DualStateDB) PointCache() *utils.PointCache {
	return d.snapshotDB.PointCache()
}

// Prepare prepares both state implementations for transaction execution.
func (d *DualStateDB) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	d.ethStateDB.Prepare(rules, sender, coinbase, dest, precompiles, txAccesses)
	d.snapshotDB.Prepare(rules, sender, coinbase, dest, precompiles, txAccesses)
}

// RevertToSnapshot reverts to a snapshot in both state implementations.
func (d *DualStateDB) RevertToSnapshot(snapshot int) {
	d.ethStateDB.RevertToSnapshot(snapshot)
	d.snapshotDB.RevertToSnapshot(snapshot)
}

// Snapshot creates a snapshot in both state implementations.
// Returns the snapshot ID from the SnapshotDB.
func (d *DualStateDB) Snapshot() int {
	d.ethStateDB.Snapshot()
	snapSnapshot := d.snapshotDB.Snapshot()
	return snapSnapshot
}

// AddLog adds a log to both state implementations.
func (d *DualStateDB) AddLog(log *types.Log) {
	d.ethStateDB.AddLog(log)
	d.snapshotDB.AddLog(log)
}

// AddPreimage adds a preimage to both state implementations.
func (d *DualStateDB) AddPreimage(hash common.Hash, preimage []byte) {
	d.ethStateDB.AddPreimage(hash, preimage)
	d.snapshotDB.AddPreimage(hash, preimage)
}

// Witness returns the witness from the SnapshotDB.
func (d *DualStateDB) Witness() *stateless.Witness {
	return d.snapshotDB.Witness()
}

// AccessEvents returns the access events from the SnapshotDB.
func (d *DualStateDB) AccessEvents() *ethstate.AccessEvents {
	return d.snapshotDB.AccessEvents()
}

// Finalise finalizes both state implementations.
func (d *DualStateDB) Finalise(deleteEmptyObjects bool) {
	d.ethStateDB.Finalise(deleteEmptyObjects)
	d.snapshotDB.Finalise(deleteEmptyObjects)
}

// Result returns the read-write set from the SnapshotDB.
// This is a SnapshotDB-specific method not part of vm.StateDB interface.
func (d *DualStateDB) Result() blocks.ReadWriteSet {
	return d.snapshotDB.Result()
}

// Logs returns the logs from the SnapshotDB.
// This is a SnapshotDB-specific method not part of vm.StateDB interface.
func (d *DualStateDB) Logs() []state.Log {
	return d.snapshotDB.Logs()
}

// Ops returns the recorded state operations from the SnapshotDB.
// This is a SnapshotDB-specific method not part of vm.StateDB interface.
func (d *DualStateDB) Ops() []StateOp {
	return d.snapshotDB.Ops()
}

// Made with Bob
