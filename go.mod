module github.com/okp4/okp4d

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.45.4
	github.com/cosmos/ibc-go/v2 v2.0.3
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/ignite-hq/cli v0.20.3
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.4.0 // indirect
	github.com/stretchr/testify v1.7.1
	github.com/tendermint/spn v0.2.1-0.20220511154430-aeab7a5b2bc0
	github.com/tendermint/tendermint v0.34.19
	github.com/tendermint/tm-db v0.6.7
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd // indirect
	google.golang.org/grpc v1.46.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
