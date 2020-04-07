module github.com/DataDog/datadog-agent

go 1.13

// Pinned to kubernetes-1.16.2
replace (
	github.com/kubernetes-incubator/custom-metrics-apiserver => github.com/kubernetes-incubator/custom-metrics-apiserver v0.0.0-20190918110929-3d9be26a50eb
	k8s.io/api => k8s.io/api v0.0.0-20191016110408-35e52d86657a
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20191016113550-5357c4baaf65
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20191016112112-5190913f932d
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20191016114015-74ad18325ed5
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20191016115326-20453efc2458
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20191016115129-c07a134afb42
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20191004115455-8e001e5d1894
	k8s.io/component-base => k8s.io/component-base v0.0.0-20191016111319-039242c015a9
	k8s.io/cri-api => k8s.io/cri-api v0.0.0-20190828162817-608eb1dad4ac
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20191016115521-756ffa5af0bd
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20191016112429-9587704a8ad4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20191016114939-2b2b218dc1df
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20191016114407-2e83b6f20229
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20191016114748-65049c67a58b
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191016120415-2ed914427d51
	k8s.io/kubelet => k8s.io/kubelet v0.0.0-20191016114556-7841ed97f1b2
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20191016115753-cf0698c3a16b
	k8s.io/metrics => k8s.io/metrics v0.0.0-20191016113814-3b1a734dba6e
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20191016112829-06bb3c9d77c9
)

// Fix tooling version
require (
	github.com/fzipp/gocyclo v0.0.0-20150627053110-6acd4345c835 // indirect
	github.com/gordonklaus/ineffassign v0.0.0-20200309095847-7953dde2c7bf // indirect
	github.com/shuLhan/go-bindata v3.4.0+incompatible // indirect
)

replace github.com/spf13/viper v1.6.2 => github.com/DataDog/viper v1.7.1

replace github.com/coreos/etcd v3.2.23+incompatible => github.com/coreos/etcd v3.2.23+incompatible

replace github.com/florianl/go-conntrack v0.1.0 => github.com/florianl/go-conntrack v0.1.1-0.20191002182014-06743d3a59db

// replace github.com/prometheus/common v0.9.1 => github.com/prometheus/common v0.4.1

replace github.com/mholt/archiver v3.1.1+incompatible => github.com/mholt/archiver v2.0.1-0.20171012052341-26cf5bb32d07+incompatible

replace github.com/prometheus/client_golang v1.5.1 => github.com/prometheus/client_golang v0.9.2

replace github.com/ugorji/go v0.0.0 => github.com/ugorji/go v0.0.0-20170107133203-ded73eae5db7

replace github.com/containerd/containerd => github.com/containerd/containerd v1.2.13

replace github.com/containerd/cgroups => github.com/containerd/cgroups v0.0.0-20200327175542-b44481373989

replace github.com/docker/distribution => github.com/docker/distribution v2.7.1-0.20190104202606-0ac367fd6bee+incompatible

replace github.com/docker/docker => github.com/docker/docker v17.12.0-ce-rc1.0.20200309214505-aa6a9891b09c+incompatible

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd v0.0.0-20180202092358-40e2722dffea

require (
	code.cloudfoundry.org/bbs v0.0.0-20200403215808-d7bc971db0db
	code.cloudfoundry.org/cfhttp/v2 v2.0.0 // indirect
	code.cloudfoundry.org/clock v1.0.0 // indirect
	code.cloudfoundry.org/consuladapter v0.0.0-20200131002136-ac1daf48ba97 // indirect
	code.cloudfoundry.org/diego-logging-client v0.0.0-20200130234554-60ef08820a45 // indirect
	code.cloudfoundry.org/executor v0.0.0-20200218194701-024d0bdd52d4 // indirect
	code.cloudfoundry.org/garden v0.0.0-20200224155059-061eda450ad9 // indirect
	code.cloudfoundry.org/go-diodes v0.0.0-20190809170250-f77fb823c7ee // indirect
	code.cloudfoundry.org/go-loggregator v7.4.0+incompatible // indirect
	code.cloudfoundry.org/locket v0.0.0-20200131001124-67fd0a0fdf2d // indirect
	code.cloudfoundry.org/rep v0.0.0-20200325195957-1404b978e31e // indirect
	code.cloudfoundry.org/rfc5424 v0.0.0-20180905210152-236a6d29298a // indirect
	code.cloudfoundry.org/tlsconfig v0.0.0-20200131000646-bbe0f8da39b3 // indirect
	github.com/DataDog/agent-payload v0.0.0-20200331195314-0715d2b415ff
	github.com/DataDog/datadog-go v3.5.0+incompatible
	github.com/DataDog/gohai v0.0.0-20200124154531-8cbe900337f1
	github.com/DataDog/gopsutil v0.0.0-20191127151039-7e1a4eadb59e
	github.com/DataDog/mmh3 v0.0.0-20200316233529-f5b682d8c981 // indirect
	github.com/DataDog/watermarkpodautoscaler v0.1.0
	github.com/DataDog/zstd v1.4.5
	github.com/beevik/ntp v0.3.0
	github.com/benesch/cgosymbolizer v0.0.0-20190515212042-bec6fe6e597b
	github.com/bhmj/jsonslice v0.0.0-20200323023432-92c3edaad8e2
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/clbanning/mxj v1.8.4
	github.com/cockroachdb/cmux v0.0.0-20170110192607-30d10be49292 // indirect
	github.com/containerd/cgroups v0.0.0
	github.com/containerd/containerd v1.2.13
	github.com/containerd/continuity v0.0.0-20200228182428-0f16d7a0959c // indirect
	github.com/containerd/fifo v0.0.0-20191213151349-ff969a566b00 // indirect
	github.com/coreos/go-systemd v0.0.0
	github.com/docker/docker v17.12.0-ce-rc1.0.20200309214505-aa6a9891b09c+incompatible
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/florianl/go-conntrack v0.1.0
	github.com/frankban/quicktest v1.9.0 // indirect
	github.com/go-ini/ini v1.55.0
	github.com/go-test/deep v1.0.5 // indirect
	github.com/gogo/googleapis v1.3.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/gopacket v1.1.17
	github.com/gorilla/mux v1.7.4
	github.com/hashicorp/consul/api v1.4.0
	github.com/hectane/go-acl v0.0.0-20190604041725-da78bae5fc95
	github.com/ianlancetaylor/cgosymbolizer v0.0.0-20200321040036-d43e30eacb43 // indirect
	github.com/iovisor/gobpf v0.0.0-20200329161226-8b2cce9dac28
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/kubernetes-incubator/custom-metrics-apiserver v0.0.0-20200323093244-5046ce1afe6b
	github.com/lxn/walk v0.0.0-20191128110447-55ccb3a9f5c1
	github.com/lxn/win v0.0.0-20191128105842-2da648fda5b4
	github.com/mdlayher/netlink v1.1.0
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/oliveagle/jsonpath v0.0.0-20180606110733-2e52cf6e6852 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/pierrec/lz4 v2.5.0+incompatible // indirect
	github.com/prometheus/client_golang v1.5.1
	github.com/samuel/go-zookeeper v0.0.0-20190923202752-2cc03de413da
	github.com/shirou/gopsutil v2.20.3+incompatible
	github.com/spf13/cobra v0.0.7
	github.com/spf13/viper v1.6.2
	github.com/tedsuo/ifrit v0.0.0-20191009134036-9a97d0632f00 // indirect
	github.com/tinylib/msgp v1.1.2
	github.com/twmb/murmur3 v1.1.3
	github.com/ulikunitz/xz v0.5.7 // indirect
	github.com/vito/go-sse v1.0.0 // indirect
	github.com/zorkian/go-datadog-api v2.28.0+incompatible // indirect
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	gopkg.in/Knetic/govaluate.v3 v3.0.0 // indirect
	gopkg.in/ini.v1 v1.55.0 // indirect
	gopkg.in/zorkian/go-datadog-api.v2 v2.28.0
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/apiserver v0.17.3
	k8s.io/metrics v0.17.3
)

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/Masterminds/vcs v1.13.1 // indirect
	github.com/Microsoft/go-winio v0.4.11
	github.com/aws/aws-sdk-go v1.25.11
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/containerd/typeurl v0.0.0-20190228175220-2a93cfde8c20
	github.com/coreos/etcd v3.3.15+incompatible
	github.com/coreos/go-semver v0.3.0
	github.com/docker/go-connections v0.3.0
	github.com/fatih/color v1.9.0
	github.com/fullstorydev/grpcui v0.2.1 // indirect
	github.com/go-ole/go-ole v1.2.4
	github.com/godbus/dbus v0.0.0-20151105175453-c7fdd8b5cd55
	github.com/gogo/protobuf v1.3.1
	github.com/goreleaser/goreleaser v0.131.1 // indirect
	github.com/hashicorp/golang-lru v0.5.1
	github.com/jmank88/nuts v0.4.0 // indirect
	github.com/json-iterator/go v1.1.9
	github.com/miekg/dns v1.0.14
	github.com/nightlyone/lockfile v1.0.0 // indirect
	github.com/opencontainers/runtime-spec v1.0.2 // indirect
	github.com/openshift/api v3.9.1-0.20190424152011-77b8897ec79a+incompatible
	github.com/pkg/errors v0.9.1
	github.com/sdboyer/constext v0.0.0-20170321163424-836a14457353 // indirect
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4
	github.com/spf13/afero v1.2.2
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.5.1
	github.com/tebeka/go2xunit v1.4.10 // indirect
	github.com/urfave/negroni v1.0.0
	github.com/vishvananda/netns v0.0.0-20171111001504-be1fbeda1936
	golang.org/x/lint v0.0.0-20190930215403-16217165b5de
	golang.org/x/mobile v0.0.0-20190312151609-d3739f865fa6
	golang.org/x/perf v0.0.0-20200318175901-9c9101da8316 // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd
	golang.org/x/text v0.3.2
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.27.1
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/cri-api v0.0.0
	k8s.io/kubernetes v1.15.4
)
