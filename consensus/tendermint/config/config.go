// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package config

import (
	"sync"
)

type ProposerPolicy uint64

const (
	RoundRobin ProposerPolicy = iota
	Sticky
)

type Config struct {
	RequestTimeout uint64         `toml:",omitempty"` // The timeout for each Istanbul round in milliseconds.
	BlockPeriod    uint64         `toml:",omitempty"` // Default minimum difference between two consecutive block's timestamps in second
	ProposerPolicy ProposerPolicy `toml:",omitempty"` // The policy for proposer selection
	Epoch          uint64         `toml:",omitempty"` // The number of blocks after which to checkpoint and reset the pending votes

	sync.RWMutex
}

func DefaultConfig() *Config {
	return &Config{
		RequestTimeout: 10000,
		BlockPeriod:    1,
		ProposerPolicy: RoundRobin,
		Epoch:          30000,
	}
}

func (cfg *Config) SetProposerPolicy(p ProposerPolicy) {
	cfg.Lock()
	cfg.ProposerPolicy = p
	cfg.Unlock()
}

func (cfg *Config) GetProposerPolicy() ProposerPolicy {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.ProposerPolicy
}
