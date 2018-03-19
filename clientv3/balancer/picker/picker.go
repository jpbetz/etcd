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
	"fmt"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

// Policy defines balancer picker policy.
type Policy uint8

const (
	// TODO: custom picker is not supported yet.
	// custom defines custom balancer picker.
	custom Policy = iota

	// Roundrobin switches endpoints in order.
	Roundrobin Policy = iota

	// TODO: priotize leader, health-check, weighted roundrobin
)

func (p Policy) String() string {
	switch p {
	case custom:
		panic("'custom' picker policy is not supported yet")
	case Roundrobin:
		// TODO: include gRPC version?
		return "clientv3-balancer-roundrobin"
	default:
		panic(fmt.Errorf("invalid balancer picker policy (%d)", p))
	}
}

// Picker defines balancer Picker methods.
type Picker interface {
	balancer.Picker

	// UpdateAddrs updates current endpoints in picker.
	// Used when endpoints are updated manually.
	UpdateAddrs(addrs []resolver.Address)
}

/*
	// UpdateSubConns updates sub-connection snapshot in picker.
	// Used when endpoints are updated manually.
	UpdateSubConns(scs []balancer.SubConn)
*/
