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

package balancer

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/coreos/etcd/clientv3/balancer/picker"
	"github.com/coreos/etcd/etcdserver/etcdserverpb"
	"github.com/coreos/etcd/pkg/logutil"
	"github.com/coreos/etcd/pkg/mock/mockserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

func TestBalancer(t *testing.T) {
	testCases := []struct {
		name        string
		serverCount int
	}{
		{
			name:        "rrBalance_1healthy",
			serverCount: 1,
		},
		{
			name:        "rrBalance_3healthy",
			serverCount: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, closeResolver := manual.GenerateAndRegisterManualResolver()
			defer closeResolver()

			logger := logutil.NewLogger(grpclog.NewLoggerV2WithVerbosity(os.Stderr, os.Stderr, os.Stderr, 5))
			balancer := NewRoundRobin(Config{Policy: picker.Roundrobin, ExperimentalLogger: logger})

			svrs, err := mockserver.StartMockServers(tc.serverCount)
			if err != nil {
				t.Fatalf("Failed to start mock servers: %s", err)
			}
			defer svrs.Stop()

			conn, err := grpc.Dial(fmt.Sprintf("%s:///mock.server", r.Scheme()), grpc.WithInsecure(), grpc.WithBalancerName(balancer.Name()))
			if err != nil {
				t.Fatalf("Failed to dial mock server: %s", err)
			}
			defer conn.Close()

			client := etcdserverpb.NewKVClient(conn)

			for _, svr := range svrs {
				r.NewAddress([]resolver.Address{{Addr: svr.Address}})
			}

			if _, err := client.Range(context.Background(), &etcdserverpb.RangeRequest{Key: []byte("/x")}); err != nil {
				t.Fatalf("failed request against load balancer: %s", err)
			}
		})
	}
}
