module eth-testnet-tool

go 1.19

// Use some special sauce
replace github.com/attestantio/go-eth2-client => /0xtylerholmes/git/go-eth2-client

require (
	github.com/attestantio/go-eth2-client v0.18.3
	github.com/google/gofuzz v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.29.1
	github.com/stretchr/testify v1.8.4
	github.com/tyler-smith/go-bip39 v1.1.0
	github.com/wealdtech/go-eth2-types/v2 v2.8.2
	github.com/wealdtech/go-eth2-util v1.8.2
	golang.org/x/sync v0.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/ferranbt/fastssz v0.1.3 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/goccy/go-yaml v1.9.2 // indirect
	github.com/herumi/bls-eth-go-binary v1.31.0 // indirect
	github.com/holiman/uint256 v1.2.2 // indirect
	github.com/huandu/go-clone v1.6.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prysmaticlabs/go-bitfield v0.0.0-20210809151128-385d8c5e3fb7 // indirect
	github.com/r3labs/sse/v2 v2.10.0 // indirect
	github.com/wealdtech/go-bytesutil v1.2.1 // indirect
	go.opentelemetry.io/otel v1.16.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	go.opentelemetry.io/otel/trace v1.16.0 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/cenkalti/backoff.v1 v1.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
