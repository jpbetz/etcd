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

package mockserver

import (
	"context"
	"fmt"
	"net"

	"github.com/coreos/etcd/etcdserver/etcdserverpb"
	"google.golang.org/grpc"
)

// MockServer provides a mocked out grpc server of the etcdserver interface.
type MockServer struct {
	GrpcServer *grpc.Server
	Address    string
}

// MockServers provides a cluster of mocket out gprc servers of the etcdserver interface.
type MockServers []*MockServer

// StartMockServers creates the desired count of mock servers
// and starts them.
func StartMockServers(count int) (svrs MockServers, err error) {
	svrs = make(MockServers, count)
	defer func() {
		if err != nil {
			svrs.Stop()
		}
	}()

	for i := 0; i < count; i++ {
		listener, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			return nil, fmt.Errorf("Failed to listen %v", err)
		}

		svr := grpc.NewServer()
		etcdserverpb.RegisterKVServer(svr, &mockKVServer{})
		svrs[i] = &MockServer{GrpcServer: svr, Address: listener.Addr().String()}
		go func(svr *grpc.Server, l net.Listener) {
			svr.Serve(l)
		}(svr, listener)
	}

	return svrs, nil
}

// Stop stops the mock server, immediately closing all open connections and listeners.
func (svrs MockServers) Stop() {
	for _, svr := range svrs {
		svr.GrpcServer.Stop()
	}
}

type mockKVServer struct{}

func (m *mockKVServer) Range(context.Context, *etcdserverpb.RangeRequest) (*etcdserverpb.RangeResponse, error) {
	return &etcdserverpb.RangeResponse{}, nil
}
func (m *mockKVServer) Put(context.Context, *etcdserverpb.PutRequest) (*etcdserverpb.PutResponse, error) {
	return &etcdserverpb.PutResponse{}, nil
}
func (m *mockKVServer) DeleteRange(context.Context, *etcdserverpb.DeleteRangeRequest) (*etcdserverpb.DeleteRangeResponse, error) {
	return &etcdserverpb.DeleteRangeResponse{}, nil
}
func (m *mockKVServer) Txn(context.Context, *etcdserverpb.TxnRequest) (*etcdserverpb.TxnResponse, error) {
	return &etcdserverpb.TxnResponse{}, nil
}
func (m *mockKVServer) Compact(context.Context, *etcdserverpb.CompactionRequest) (*etcdserverpb.CompactionResponse, error) {
	return &etcdserverpb.CompactionResponse{}, nil
}
