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
	"github.com/coreos/etcd/clientv3/balancer/picker"
	"github.com/coreos/etcd/pkg/logutil"
)

// Config defines balancer configurations.
type Config struct {
	Policy picker.Policy

	// ExperimentalLogger configures balancer logging.
	// TODO: use structured logging
	ExperimentalLogger logutil.Logger
}
