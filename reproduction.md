Reproduction for https://github.com/kubernetes/kubernetes/issues/83028

### 1. Checkout this branch and cd into it since we're going to use these test certs:

- https://github.com/jpbetz/etcd/blob/etcd-lb-dnsname-failover/integration/fixtures/server-ca-csr-dnsname1.json
- https://github.com/jpbetz/etcd/blob/etcd-lb-dnsname-failover/integration/fixtures/server-ca-csr-dnsname2.json
- https://github.com/jpbetz/etcd/blob/etcd-lb-dnsname-failover/integration/fixtures/server-ca-csr-dnsname3.json


### 2. Start an etcd cluster

> etcd --name infra1 --listen-client-urls https://127.0.0.1:2379 --advertise-client-urls https://member1.etcd.local:2379 --listen-peer-urls http://127.0.0.1:12380 --initial-advertise-peer-urls http://127.0.0.1:12380 --initial-cluster-token etcd-cluster-1 --initial-cluster 'infra1=http://127.0.0.1:12380,infra2=http://127.0.0.1:22380,infra3=http://127.0.0.1:32380' --initial-cluster-state new --enable-pprof --cert-file=integration/fixtures/server-dnsname1.crt --key-file=integration/fixtures/server-dnsname1.key.insecure

> etcd --name infra2 --listen-client-urls https://127.0.0.1:22379 --advertise-client-urls https://member2.etcd.local:22379 --listen-peer-urls http://127.0.0.1:22380 --initial-advertise-peer-urls http://127.0.0.1:22380 --initial-cluster-token etcd-cluster-1 --initial-cluster 'infra1=http://127.0.0.1:12380,infra2=http://127.0.0.1:22380,infra3=http://127.0.0.1:32380' --initial-cluster-state new --enable-pprof --cert-file=integration/fixtures/server-dnsname2.crt --key-file=integration/fixtures/server-dnsname2.key.insecure

> etcd --name infra3 --listen-client-urls https://127.0.0.1:32379 --advertise-client-urls https://member3.etcd.local:32379 --listen-peer-urls http://127.0.0.1:32380 --initial-advertise-peer-urls http://127.0.0.1:32380 --initial-cluster-token etcd-cluster-1 --initial-cluster 'infra1=http://127.0.0.1:12380,infra2=http://127.0.0.1:22380,infra3=http://127.0.0.1:32380' --initial-cluster-state new --enable-pprof --cert-file=integration/fixtures/server-dnsname3.crt --key-file=integration/fixtures/server-dnsname3.key.insecure

### 3. Modify /etc/hosts to include

```

```g when kube-apiserver is started*

kube-apiserver terminates with log:

> ...
> member2.etcd.local, not member1.etcd.local". Reconnecting...
> W0924 14:42:29.368784  248333 clientconn.go:1120] grpc: addrConn.createTransport failed to connect to {https://member3.etcd.local:32379 0  <nil>}. Err :connection error: desc = "transport: authentication handshake failed: x509: certificate is valid for member3.etcd.local, not member1.etcd.local". Reconnecting...
> panic: context deadline exceeded
> 
> goroutine 1 [running]:
> k8s.io/kubernetes/vendor/k8s.io/apiextensions-apiserver/pkg/registry/customresourcedefinition.NewREST(0xc000992690, 0x7aaae80, 0xc000278240, 0xc000278468)
>         /go/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apiextensions-apiserver/pkg/registry/customresourcedefinition/etcd.go:56 +0x41c
> k8s.io/kubernetes/vendor/k8s.io/apiextensions-apiserver/pkg/apiserver.completedConfig.New(0xc000b18ca0, 0xc0004fe6c8, 0x7b63340, 0xaaa58d8, 0x10, 0x0, 0x0)
>         /go/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/k8s.io/apiextensions-apiserver/pkg/apiserver/apiserver.go:147 +0x1586
> k8s.io/kubernetes/cmd/kube-apiserver/app.createAPIExtensionsServer(0xc0004fe6c0, 0x7b63340, 0xaaa58d8, 0x0, 0x7aaaae0, 0xc000440f30)
>         /go/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/cmd/kube-apiserver/app/apiextensions.go:95 +0x59
> k8s.io/kubernetes/cmd/kube-apiserver/app.CreateServerChain(0xc0008822c0, 0xc0002fc720, 0x44babba, 0xc, 0xc000aa3ca8)
>         /go/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/cmd/kube-apiserver/app/server.go:182 +0x2bb
> k8s.io/kubernetes/cmd/kube-apiserver/app.Run(0xc0008822c0, 0xc0002fc720, 0x0, 0x0)
>         /go/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/cmd/kube-apiserver/app/server.go:151 +0x102
> k8s.io/kubernetes/cmd/kube-apiserver/app.NewAPIServerCommand.func1(0xc0002af680, 0xc0001c6fc0, 0x0, 0x9, 0x0, 0x0)
>         /go/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/cmd/kube-apiserver/app/server.go:118 +0x104
> k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).execute(0xc0002af680, 0xc00004c0b0, 0x9, 0x9, 0xc0002af680, 0xc00004c0b0)
>         /go/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:826 +0x465
> k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).ExecuteC(0xc0002af680, 0x464cee0, 0xaa87560, 0xc000aa3f88)
>         /go/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:914 +0x2fc
> k8s.io/kubernetes/vendor/github.com/spf13/cobra.(*Command).Execute(...)
>         /go/src/k8s.io/kubernetes/_output/dockerized/go/src/k8s.io/kubernetes/vendor/github.com/spf13/cobra/command.go:864
> main.main()
>         _output/dockerized/go/src/k8s.io/kubernetes/cmd/kube-apiserver/apiserver.go:43 +0xc9

### Results with 1.15

*With all etcd members running when kube-apiserver is started and then "infra1" etcd is stopped*

kube-apiserver log:

> ...
> W0924 14:40:33.985557  245246 clientconn.go:1251] grpc: addrConn.createTransport failed to connect to {member3.etcd.local:32379 0  <nil>}. Err :connection error: desc = "transport: authentication handshake failed: x509: certificate is valid for member3.etcd.local, not member1.etcd.local". Reconnecting...
> W0924 14:40:34.003799  245246 clientconn.go:1251] grpc: addrConn.createTransport failed to connect to {member2.etcd.local:22379 0  <nil>}. Err :connection error: desc = "transport: authentication handshake failed: x509: certificate is valid for member2.etcd.local, not member1.etcd.local". Reconnecting...
> ...


*With "infra1" etcd not running when kube-apiserver is started*

kube-apiserver terminates with log:

> ...
> W0924 14:36:55.367773  242178 clientconn.go:1251] grpc: addrConn.createTransport failed to connect to {member1.etcd.local:2379 0  <nil>}. Err :connection error: desc = "transport: Error while dialing dial tcp 127.0.0.1:2379: connect: connection refused".
> Reconnecting...
> W0924 14:36:55.705694  242178 clientconn.go:1251] grpc: addrConn.createTransport failed to connect to {member2.etcd.local:22379 0  <nil>}. Err :connection error: desc = "transport: authentication handshake failed: x509: certificate is valid for member2.etcd.local, not member1.etcd.local". Reconnecting...
> W0924 14:36:58.312776  242178 asm_amd64.s:1337] Failed to dial member2.etcd.local:22379: context canceled; please retry.
> W0924 14:36:58.312792  242178 asm_amd64.s:1337] Failed to dial member3.etcd.local:32379: context canceled; please retry.
> F0924 14:36:58.312719  242178 storage_decorator.go:57] Unable to create storage backend: config (&{ /registry {[https://member1.etcd.local:2379 https://member2.etcd.local:22379 https://member3.etcd.local:32379]   /usr/local/google/home/jpbetz/projects/etcd-io/src/go.etcd.io/etcd/integration/fixtures/ca.crt} true 0xc000832480 apiextensions.k8s.io/v1beta1 <nil> 5m0s 1m0s}), err (dial tcp 127.0.0.1:2379: connect: connection refused)
