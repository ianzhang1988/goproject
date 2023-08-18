module goproject

go 1.16

require (
	github.com/Luzifer/go-openssl/v4 v4.1.0
	github.com/armon/go-radix v1.0.0
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/genjidb/genji v0.15.2
	github.com/go-gota/gota v0.12.0
	github.com/go-logr/stdr v1.2.2
	github.com/gogo/protobuf v1.3.2
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-memdb v1.3.4
	github.com/ipfs/go-datastore v0.6.0
	github.com/ipfs/go-log/v2 v2.5.1
	github.com/libp2p/go-libp2p v0.23.2
	github.com/libp2p/go-libp2p-core v0.20.1
	github.com/libp2p/go-libp2p-gostream v0.5.0
	github.com/libp2p/go-libp2p-http v0.4.0
	github.com/libp2p/go-libp2p-kad-dht v0.18.0
	github.com/libp2p/go-libp2p/examples v0.0.0-20220929192648-7828f3e0797e
	github.com/multiformats/go-multiaddr v0.7.0
	github.com/segmentio/kafka-go v0.4.38
	github.com/spacemonkeygo/openssl v0.0.0-20181017203307-c2dcc5cca94a
	github.com/tetratelabs/wazero v1.3.1
	github.com/vadv/gopher-lua-libs v0.4.1
	github.com/vmihailenco/msgpack/v5 v5.3.5
	github.com/yuin/gopher-lua v1.1.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.32.0
	go.opentelemetry.io/otel v1.11.1
	go.opentelemetry.io/otel/exporters/jaeger v1.11.1
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.10.0
	go.opentelemetry.io/otel/sdk v1.11.1
	go.opentelemetry.io/otel/trace v1.11.1
)

// replace github.com/vadv/gopher-lua-libs => github.com/ianzhang1988/gopher-lua-libs v0.0.0-20230809092812-8444098a793d
replace github.com/vadv/gopher-lua-libs => /data/zhangyang/gopher-lua-libs
