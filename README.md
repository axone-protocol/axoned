# økp4

> Implementation of **økp4**, a blockchain for the decentralized digital commons - built using the
> [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and [Starport](https://starport.com).

[![build](https://github.com/okp4/okp4d/actions/workflows/build.yml/badge.svg)](https://github.com/okp4/okp4d/actions/workflows/build.yml)
[![test](https://github.com/okp4/okp4d/actions/workflows/test.yml/badge.svg)](https://github.com/okp4/okp4d/actions/workflows/test.yml)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![license](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

## Configuring locally an instance with one node

Install locally the okp4d executable with:

```bash
make install
```

The location of okp4d in the system can be displayed with:

```bash
which okp4d
```

Initialize private validator, p2p, genesis, and application configuration files with:

```bash
okp4d --home genesis init localnode
```

Replace the configuration thus generated with a custom configuration as follows:

```bash
cp chains/testnet/pre-genesis.json genesis/config/genesis.json
```

Execute now an ```okp4d keys add``` command as follows to create two sets of keys
called "staker" and "wallet2":

```bash
okp4d --home genesis keys --keyring-backend test add staker
```

```bash
okp4d --home genesis keys --keyring-backend test add wallet2
```

Now you may add accounts from keyring to genesis.json with:

```bash
okp4d --home genesis add-genesis-account staker 1000000000uknow --keyring-backend test
```

```bash
okp4d --home genesis add-genesis-account wallet2 1000000000uknow --keyring-backend test
```

Run the following command to generate a genesis transaction that upgrades "staker" into a validator account with a self-delegation:

```bash
okp4d --home genesis gentx staker 1000000uknow --keyring-backend test --node-id $(okp4d --home genesis tendermint show-node-id) --chain-id okp4-testnet-1
```

Execute the following command to collect all of the transactions and configure them in the genesis.json file:

```bash
okp4d --home genesis collect-gentxs
```

## Starting the instance

Now, our first blockchain's node can be started with:

```bash
okp4d --home genesis start --moniker localnode
```

## Transferring tokens between accounts

Store addresses of accounts into env variables:

```bash
export STAKER_ADDR=$(okp4d --home genesis keys --keyring-backend test show -a staker)
export WALLET2_ADDR=$(okp4d --home genesis keys --keyring-backend test show -a wallet2)
```

Display the current balances on those accounts:

```bash
okp4d --home genesis query bank balances $STAKER_ADDR
okp4d --home genesis query bank balances $WALLET2_ADDR
```

Transfer 1000uknow from staker to wallet2:

```bash
okp4d --home genesis tx bank send $STAKER_ADDR $WALLET2_ADDR 1000uknow --keyring-backend test --chain-id okp4-testnet-1
```

Display the balances again to verify the transfer.
