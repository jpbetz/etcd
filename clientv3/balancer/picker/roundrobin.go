// Copyright 2018 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package picker

import (
	"context"
	"sync"

	"github.com/coreos/etcd/pkg/logutil"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

// NewRoundrobin returns a new roundrobin picker.
func NewRoundrobin(readySCs map[resolver.Address]balancer.SubConn, lg logutil.Logger) Picker {
	scs := make([]balancer.SubConn, 0, len(readySCs))
	for _, sc := range readySCs {
		scs = append(scs, sc)
	}

	pk := &rrPicker{lg: lg, subConns: scs}

	// TODO: update endpoints
	// TODO: add logging
	// pk.updateSubConns = func(scs []balancer.SubConn) {
	// 	pk.mu.Lock()
	// 	pk.subConns = scs
	// 	pk.mu.Unlock()
	// }

	// TODO: sc.UpdateAddresses?
	pk.updateAddrs = func(addrs []resolver.Address) {
		pk.mu.Lock()
		// close all resolved sub-connections first
		for _, sc := range pk.subConns {
			sc.UpdateAddresses([]resolver.Address{})
		}
		pk.mu.Unlock()
	}

	return pk
}

type rrPicker struct {
	lg logutil.Logger

	mu sync.RWMutex
	// index of currently pinned address(or connection)
	pin      int
	subConns []balancer.SubConn

	updateAddrs func(addrs []resolver.Address)
}

func (p *rrPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	p.mu.RLock()
	n := len(p.subConns)
	p.mu.RUnlock()
	if n == 0 {
		return nil, nil, balancer.ErrNoSubConnAvailable
	}

	p.mu.Lock()
	sc := p.subConns[p.pin]
	p.pin = (p.pin + 1) % len(p.subConns)
	p.lg.Infof("(p *rrPicker) Pick on %d", p.pin)
	p.mu.Unlock()

	return sc, func(info balancer.DoneInfo) { p.lg.Infof("balancer.DoneInfo: %+v", info) }, nil
}

func (p *rrPicker) UpdateAddrs(addrs []resolver.Address) { p.updateAddrs(addrs) }
