[![okp4 github banner](https://raw.githubusercontent.com/okp4/okp4d/main/etc/okp4-banner.png)](https://okp4.network)

<p align="center">
  <a href="https://discord.gg/okp4"><img src="https://img.shields.io/badge/Discord-7289DA?style=for-the-badge&logo=discord&logoColor=white" /></a> &nbsp;
  <a href="https://www.linkedin.com/company/okp4-open-knowledge-protocol-for"><img src="https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white" /></a> &nbsp;
  <a href="https://twitter.com/OKP4_Protocol"><img src="https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white" /></a> &nbsp;
  <a href="https://medium.com/okp4"><img src="https://img.shields.io/badge/Medium-12100E?style=for-the-badge&logo=medium&logoColor=white" /></a> &nbsp;
  <a href="https://www.youtube.com/channel/UCiOfcTaUyv2Szv4OQIepIvg"><img src="https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white" /></a>
</p>

# 𝒐𝒌𝒑4𝒅

[![lint](https://img.shields.io/github/workflow/status/okp4/okp4d/Lint?label=lint&style=for-the-badge&logo=github)](https://github.com/okp4/okp4d/actions/workflows/lint.yml)
[![build](https://img.shields.io/github/workflow/status/okp4/okp4d/Build?label=build&style=for-the-badge&logo=github)](https://github.com/okp4/okp4d/actions/workflows/build.yml)
[![test](https://img.shields.io/github/workflow/status/okp4/okp4d/Test?label=test&style=for-the-badge&logo=github)](https://github.com/okp4/okp4d/actions/workflows/test.yml)
[![discord](https://img.shields.io/discord/946759919678406696.svg?label=discord&logo=discord&logoColor=white&style=for-the-badge)](https://discord.gg/GHNZh4SaJ3)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge&logo=conventionalcommits)](https://conventionalcommits.org)
[![license](https://img.shields.io/github/license/okp4/okp4d.svg?label=License&style=for-the-badge)](https://opensource.org/licenses/Apache-2.0)

> `OKP4` is a public [dPoS](https://en.bitcoinwiki.org/wiki/DPoS) layer 1 specifically designed to enable communities to trustlessly share data, algorithms and resources to build the Dataverse - An open world where everybody can create or participate in custom ecosystems (with common governance mechanisms, sharing rules, business models...) to build a new generation of dApps way beyond Decentralized Finance.

## The protocol

`okp4d` is the node of the [OKP4](https://okp4.network) network built on the [Cosmos SDK] 💫 & [Tendermint] consensus, and designed to become a hub of incentivized data providers, developers, data scientists & users collaborating to generate value from data and algorithms.

Read more in the [introduction blog post](https://blog.okp4.network/what-is-okp4-b6bb058ae221). For a high-level overview of the OKP4 protocol and network economics, check out the [white paper](https://docs.okp4.network/docs/whitepaper/abstract).

## OKP4 Networks

| **Network** | Status | Description                                                                                                                                                                                |    Chain-ID     | URL                                                                                                  |
|-------------|:------:|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:---------------:|------------------------------------------------------------------------------------------------------|
| `devnet`    |   ✅    | Developers-focused network of our blockchain for build/test purposes.                                                                                                                      | `okp4-devnet-1` | api:&nbsp;<https://api.devnet.okp4.network> <br/>rpc:&nbsp;<https://api.devnet.okp4.network:443/rpc> |
| `testnet`   |   🚧   | The easiest place to create test Data Spaces if you're not quite ready yet to dive into building applications on our main network. It is also the ideal place to start as a node operator. |        -        | -                                                                                                    |
| `mainnet`   |   🔜   |                                                                                                                                                                                            |        -        | -                                                                                                    |

## Want to become a validator?

Validators are responsible for securing the okp4 network. Validator responsibilities include maintaining a functional [node](https://docs.okp4.network/docs/nodes/run-node) with constant uptime and providing a sufficient amount of $KNOW as stake. In exchange for this service, validators receive block rewards and transaction fees.

Want to become a validator? 👉 [Checkout the documentation!](https://docs.okp4.network/docs/nodes/introduction)

## Developing & contributing

`okp4d` is written in [Go] and built using [Cosmos SDK].

### Prerequisites

- install [Go] `1.18+` following instructions from the [official Go documentation](https://golang.org/doc/install);
- verify that [Docker] is properly installed and if not, follow the [instructions](https://docs.docker.com) for your environment;
- the project comes with a convenient `Makefile` so verify that [`make`](https://fr.wikipedia.org/wiki/Make) is properly installed.

### Build

To build the `okp4d` node, invoke the goal `build` of the `Makefile`:

```sh
make build
```

The binary will be generated under the folder `target/dist`.

## Supported platforms

The `okp4d` blockchain currently supports the following builds:

| **Platform** | **Arch** |       **Status**       |
|--------------|----------|:----------------------:|
| Darwin       | amd64    |           ✅            |
| Darwin       | arm64    |           ✅            |
| Linux        | amd64    |           ✅            |
| Linux        | arm64    |           ✅            |
| Windows      | amd64    | ️🚫<br/> Not supported |

> Note: as the blockchain depends on [CosmWasm/wasmvm](https://github.com/CosmWasm/wasmvm), we only support the targets
> supported by this project.

## Releases

All releases can be found [here](https://github.com/okp4/okp4d/releases).

`okp4d` follows the [Semantic Versioning 2.0.0](https://semver.org/) to determine when and how the version changes, and
we also apply the philosophical principles of [release early - release often](https://en.wikipedia.org/wiki/Release_early,_release_often).

## Docker image

A docker image is available on [Docker hub](https://hub.docker.com/r/okp4/okp4d):

```bash
docker pull okp4/okp4d:latest
```

### Get the documentation

```bash
docker run okp4/okp4d:latest --help
```

### Query a running network

Example:

```bash
API_URL=https://api.devnet.okp4.network:443/rpc
WALLET=okp41pmkq300lrngpkeprygfrtag0xpgp9z92c7eskm

docker run okp4/okp4d:latest query bank balances $WALLET --chain-id okp4-devnet-1 --node $API_URL
 ```

### Create a wallet

```bash
docker run -v $(pwd)/home:/home okp4/okp4d:latest keys add my-wallet --keyring-backend test --home /home 
```

### Start a node

Everything you need to start a node and even more is explained here : <https://docs.okp4.network/docs/nodes/run-node>

```bash
MONIKER=node-in-my-name
CHAIN_ID=localnet-okp4-1

docker run -v $(pwd)/home:/home okp4/okp4d:latest init $MONIKER --chain-id $CHAIN_ID --home /home 
```

This will create a home folder, you can then update the `config/genesis.json` with one of this ones : <https://github.com/okp4/okp4d/tree/main/chains/>

#### Join a running network

Set `persistent_peers` in `config/config.toml` file.

#### Then start the node

```bash
docker run -v $(pwd)/home:/home okp4/okp4d:latest start --home /home
```

## Bug reports & feature requests

If you notice anything not behaving how you expected, if you would like to make a suggestion or would like
to request a new feature, please open a [**new issue**](https://github.com/okp4/okp4d/issues/new/choose). We appreciate any help
you're willing to give!

> Don't hesitate to ask if you are having trouble setting up your project repository, creating your first branch or
> configuring your development environment. Mentors and maintainers are here to help!

## Community

The [**OKP4 Discord Server**](https://discord.gg/GHNZh4SaJ3) is our primary chat channel for the open-source community,
software developers and node operators.

Please reach out to us and say hi 👋, we're happy to help there.

[Cosmos SDK]: https://v1.cosmos.network/sdk
[Docker]: https://www.docker.com/
[Go]: https://go.dev
[Tendermint]: https://tendermint.com/
