// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vm

import (
	"time"

	"github.com/ava-labs/avalanchego/database/versiondb"
	log "github.com/inconshreveable/log15"

	"github.com/ava-labs/quarkvm/chain"
)

func (vm *VM) pruneCall() bool {
	// Lock to prevent concurrent modification of state
	vm.ctx.Lock.Lock()
	defer vm.ctx.Lock.Unlock()

	vdb := versiondb.New(vm.db)
	defer vdb.Abort()
	removals, err := chain.PruneNext(vdb, vm.pruneLimit)
	if err != nil {
		log.Warn("unable to prune next range", "error", err)
		return false
	}
	if err := vdb.Commit(); err != nil {
		log.Warn("unable to commit pruning work", "error", err)
		return false
	}
	if err := vm.lastAccepted.SetChildrenDB(vm.db); err != nil {
		log.Error("unable to update child databases of last accepted block", "error", err)
	}
	return removals == vm.pruneLimit
}

func (vm *VM) prune() {
	log.Debug("starting prune loops")
	defer close(vm.donecPrune)

	// should retry less aggressively
	t := time.NewTimer(vm.pruneInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
		case <-vm.stopc:
			return
		}
		if vm.pruneCall() {
			t.Reset(vm.fullPruneInterval)
		} else {
			t.Reset(vm.pruneInterval)
		}
	}
}
