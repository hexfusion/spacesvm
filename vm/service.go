// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vm

import (
	"errors"
	"net/http"

	"github.com/ava-labs/avalanchego/ids"
	log "github.com/inconshreveable/log15"

	"github.com/ava-labs/quarkvm/chain"
)

var (
	ErrPoWFailed      = errors.New("PoW failed")
	ErrInvalidEmptyTx = errors.New("invalid empty transaction")
)

type Service struct {
	vm *VM
}

type PingArgs struct {
}

type PingReply struct {
	Success bool `serialize:"true" json:"success"`
}

func (svc *Service) Ping(_ *http.Request, args *PingArgs, reply *PingReply) (err error) {
	log.Info("ping")
	reply.Success = true
	return nil
}

type IssueTxArgs struct {
	Tx []byte `serialize:"true" json:"tx"`
}

type IssueTxReply struct {
	TxID    ids.ID `serialize:"true" json:"txID"`
	Success bool   `serialize:"true" json:"success"`
}

func (svc *Service) IssueTx(_ *http.Request, args *IssueTxArgs, reply *IssueTxReply) error {
	tx := new(chain.Transaction)
	if _, err := chain.Unmarshal(args.Tx, tx); err != nil {
		return err
	}
	svc.vm.Submit(tx)
	reply.TxID = tx.ID()
	reply.Success = true
	return nil
}

type CheckTxArgs struct {
	TxID ids.ID `serialize:"true" json:"txID"`
}

type CheckTxReply struct {
	Confirmed bool `serialize:"true" json:"confirmed"`
}

func (svc *Service) CheckTx(_ *http.Request, args *CheckTxArgs, reply *CheckTxReply) error {
	has, err := chain.HasTransaction(svc.vm.db, args.TxID)
	if err != nil {
		return err
	}
	reply.Confirmed = has
	return nil
}

type CurrBlockArgs struct {
}

type CurrBlockReply struct {
	BlockID ids.ID `serialize:"true" json:"blockID"`
}

func (svc *Service) CurrBlock(_ *http.Request, args *CurrBlockArgs, reply *CurrBlockReply) error {
	reply.BlockID = svc.vm.preferred
	return nil
}

type ValidBlockIDArgs struct {
	BlockID ids.ID `serialize:"true" json:"blockID"`
}

type ValidBlockIDReply struct {
	Valid bool `serialize:"true" json:"valid"`
}

func (svc *Service) ValidBlockID(_ *http.Request, args *ValidBlockIDArgs, reply *ValidBlockIDReply) error {
	valid, err := svc.vm.ValidBlockID(args.BlockID)
	if err != nil {
		return err
	}
	reply.Valid = valid
	return nil
}

type DifficultyEstimateArgs struct {
}

type DifficultyEstimateReply struct {
	Difficulty uint `serialize:"true" json:"valid"`
}

func (svc *Service) DifficultyEstimate(_ *http.Request, args *DifficultyEstimateArgs, reply *DifficultyEstimateReply) error {
	diff, err := svc.vm.DifficultyEstimate()
	if err != nil {
		return err
	}
	reply.Difficulty = diff
	return nil
}

type PrefixInfoArgs struct {
	Prefix []byte `serialize:"true" json:"prefix"`
}

type PrefixInfoReply struct {
	Info *chain.PrefixInfo `serialize:"true" json:"info"`
}

func (svc *Service) PrefixInfo(_ *http.Request, args *PrefixInfoArgs, reply *PrefixInfoReply) error {
	i, _, err := chain.GetPrefixInfo(svc.vm.db, args.Prefix)
	if err != nil {
		return err
	}
	reply.Info = i
	return nil
}
