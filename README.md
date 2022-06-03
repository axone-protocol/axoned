# økp4

> Implementation of **økp4**, a blockchain for the decentralized digital commons - built using the
> [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and [Starport](https://starport.com).

[![build](https://github.com/okp4/okp4d/actions/workflows/build.yml/badge.svg)](https://github.com/okp4/okp4d/actions/workflows/build.yml)
[![test](https://github.com/okp4/okp4d/actions/workflows/test.yml/badge.svg)](https://github.com/okp4/okp4d/actions/workflows/test.yml)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![license](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

# Starting locally an instance with one node

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

Execute now an ```okp4d keys add``` command as follows to create a new set of keys
called "staker":

```bash
okp4d --home genesis keys --keyring-backend test add staker
```

Now you may add account from keyring to genesis.json with:

```bash
okp4d --home genesis add-genesis-account staker 1000000000uknow --keyring-backend test
```

Run the following command to generate a genesis transaction that creates a validator with a self-delegation:

```bash
okp4d --home genesis gentx staker 1000000uknow --keyring-backend test --node-id $(okp4d --home genesis tendermint show-node-id) --chain-id okp4-testnet-1
```

Execute the following command to collect all of the transactions and configure them in the genesis.json file:

```bash
okp4d --home genesis collect-gentxs
```

Now, our first blockchain's node can be started with:

```bash
okp4d --home genesis start --moniker localnode
```
