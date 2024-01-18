## okp4d tx interchain-accounts host generate-packet-data

Generates protobuf or proto3 JSON encoded ICA packet data.

### Synopsis

generate-packet-data accepts a message string and serializes it (depending on the
encoding parameter) using protobuf or proto3 JSON into packet data which is outputted to stdout.
It can be used in conjunction with send-tx which submits pre-built packet data containing messages
to be executed on the host chain. The default encoding format is protobuf if none is specified;
otherwise the encoding flag can be used in combination with either "proto3" or "proto3json".

```
okp4d tx interchain-accounts host generate-packet-data [message] [flags]
```

### Examples

```
okp4d tx interchain-accounts host generate-packet-data '{
    "@type":"/cosmos.bank.v1beta1.MsgSend",
    "from_address":"cosmos15ccshhmp0gsx29qpqq6g4zmltnnvgmyu9ueuadh9y2nc5zj0szls5gtddz",
    "to_address":"cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
    "amount": [
        {
            "denom": "stake",
            "amount": "1000"
        }
    ]
}' --memo memo --encoding proto3json


okp4d tx interchain-accounts host generate-packet-data '[{
    "@type":"/cosmos.bank.v1beta1.MsgSend",
    "from_address":"cosmos15ccshhmp0gsx29qpqq6g4zmltnnvgmyu9ueuadh9y2nc5zj0szls5gtddz",
    "to_address":"cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw",
    "amount": [
        {
            "denom": "stake",
            "amount": "1000"
        }
    ]
},
{
	"@type": "/cosmos.staking.v1beta1.MsgDelegate",
	"delegator_address": "cosmos15ccshhmp0gsx29qpqq6g4zmltnnvgmyu9ueuadh9y2nc5zj0szls5gtddz",
	"validator_address": "cosmosvaloper1qnk2n4nlkpw9xfqntladh74w6ujtulwnmxnh3k",
	"amount": {
		"denom": "stake",
		"amount": "1000"
	}
}]'
```

### Options

```
      --encoding string   optional encoding format of the messages in the interchain accounts packet data
  -h, --help              help for generate-packet-data
      --memo string       optional memo to be included in the interchain accounts packet data
```

### SEE ALSO

* [okp4d tx interchain-accounts host](okp4d_tx_interchain-accounts_host.md)	 - IBC interchain accounts host transaction subcommands
