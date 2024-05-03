[![axone github banner](https://raw.githubusercontent.com/axone-protocol/.github/main/profile/static/axone-banner.png)](https://axone.xyz)

<p align="center">
  <a href="https://discord.gg/axone"><img src="https://img.shields.io/badge/Discord-7289DA?style=for-the-badge&logo=discord&logoColor=white" /></a> &nbsp;
  <a href="https://www.linkedin.com/company/axone-protocol/"><img src="https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white" /></a> &nbsp;
  <a href="https://twitter.com/axonexyz"><img src="https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white" /></a> &nbsp;
  <a href="https://blog.axone.xyz"><img src="https://img.shields.io/badge/Medium-12100E?style=for-the-badge&logo=medium&logoColor=white" /></a> &nbsp;
  <a href="https://www.youtube.com/channel/UCiOfcTaUyv2Szv4OQIepIvg"><img src="https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white" /></a>
</p>

# 𝒐𝒌𝒑4𝒅

[![lint](https://img.shields.io/github/actions/workflow/status/axone-protocol/axoned/lint.yml?label=lint&style=for-the-badge&logo=github)](https://github.com/axone-protocol/axoned/actions/workflows/lint.yml)
[![build](https://img.shields.io/github/actions/workflow/status/axone-protocol/axoned/build.yml?label=build&style=for-the-badge&logo=github)](https://github.com/axone-protocol/axoned/actions/workflows/build.yml)
[![test](https://img.shields.io/github/actions/workflow/status/axone-protocol/axoned/test.yml?label=test&style=for-the-badge&logo=github)](https://github.com/axone-protocol/axoned/actions/workflows/test.yml)
[![codecov](https://img.shields.io/codecov/c/github/axone-protocol/axoned?style=for-the-badge&token=O3FJO5QDCA&logo=codecov)](https://codecov.io/gh/axone-protocol/axoned)
[![docker-pull](https://img.shields.io/docker/pulls/axone-protocol/axoned?label=downloads&style=for-the-badge&logo=docker)](https://hub.docker.com/r/axone-protocol/axoned)
[![discord](https://img.shields.io/discord/946759919678406696.svg?label=discord&logo=discord&logoColor=white&style=for-the-badge)](https://discord.gg/axone)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge&logo=conventionalcommits)](https://conventionalcommits.org)
[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg?style=for-the-badge)](https://github.com/semantic-release/semantic-release)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg?style=for-the-badge)](https://github.com/axone-protocol/.github/blob/main/CODE_OF_CONDUCT.md)
[![license](https://img.shields.io/github/license/axone-protocol/axoned.svg?label=License&style=for-the-badge)](https://opensource.org/licenses/Apache-2.0)

> `AXONE` is a public [dPoS](https://en.bitcoinwiki.org/wiki/DPoS) layer 1 specifically designed to enable communities to trustlessly share data, algorithms and resources to build the Dataverse - An open world where everybody can create or participate in custom ecosystems (with common governance mechanisms, sharing rules, business models...) to build a new generation of dApps way beyond Decentralized Finance.

## The protocol

`axoned` is the node of the [AXONE](https://axone.xyz) network built on the [Cosmos SDK] 💫 & [Tendermint] consensus, and designed to become a hub of incentivized data providers, developers, data scientists & users collaborating to generate value from data and algorithms.

Read more in the [introduction blog post](https://blog.axone.xyz/). For a high-level overview of the AXONE protocol and network economics, check out the [white paper](https://docs.axone.xyz/docs/whitepaper/abstract).

## Want to become a validator?

Validators are responsible for securing the axone network. Validator responsibilities include maintaining a functional [node](https://docs.axone.xyz/docs/nodes/run-node) with constant uptime and providing a sufficient amount of $KNOW as stake. In exchange for this service, validators receive block rewards and transaction fees.

Want to become a validator? 👉 [Checkout the documentation!](https://docs.axone.xyz/docs/nodes/introduction)

Looking for a network to join ? 👉 [Checkout the networks!](https://github.com/axone-protocol/networks)

## Supported platforms

The `axoned` blockchain currently supports the following builds:

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

All releases can be found [here](https://github.com/axone-protocol/axoned/releases).

`axoned` follows the [Semantic Versioning 2.0.0](https://semver.org/) to determine when and how the version changes, and
we also apply the philosophical principles of [release early - release often](https://en.wikipedia.org/wiki/Release_early,_release_often).

## Docker image

For a quick start, a docker image is officially available on [Docker hub](https://hub.docker.com/r/axone-protocol/axoned).

```bash
docker pull axone-protocol/axoned:latest
```

### Get the documentation

```bash
docker run axone-protocol/axoned:latest --help
```

### Query a running network

Example:

```bash
API_URL=https://api.devnet.okp4.network:443/rpc
WALLET=okp41pmkq300lrngpkeprygfrtag0xpgp9z92c7eskm

docker run axone-protocol/axoned:latest query bank balances $WALLET --chain-id axone-devnet-1 --node $API_URL
 ```

### Create a wallet

```bash
docker run -v $(pwd)/home:/home axone-protocol/axoned:latest keys add my-wallet --keyring-backend test --home /home 
```

### Start a node

Everything you need to start a node and more is explained here: <https://docs.axone.xyz/docs/nodes/run-node>

```bash
MONIKER=node-in-my-name
CHAIN_ID=localnet-axone-1

docker run -v $(pwd)/home:/home axone-protocol/axoned:latest init $MONIKER --chain-id $CHAIN_ID --home /home 
```

This will create a home folder, you can then update the `config/genesis.json` with one of this ones : <https://github.com/axone-protocol/axoned/tree/main/chains/>

#### Join a running network

Set `persistent_peers` in `config/config.toml` file.

#### Then start the node

```bash
docker run -v $(pwd)/home:/home axone-protocol/axoned:latest start --home /home
```

## Developing & contributing

`axoned` is written in [Go] and built using [Cosmos SDK]. A number of smart contracts are also deployed on the
AXONE blockchain and hosted in the [axone-protocol/contracts](https://github.com/axone-protocol/contracts) project.

### Prerequisites

- install [Go] `1.21+` following instructions from the [official Go documentation](https://golang.org/doc/install);
- use [gofumpt](https://github.com/mvdan/gofumpt) as formatter. You can integrate it in your favorite IDE following these [instructions](https://github.com/mvdan/gofumpt#installation) or invoke the makefile `make format-go`;
- verify that [Docker] is properly installed and if not, follow the [instructions](https://docs.docker.com) for your environment;
- verify that [`make`](https://fr.wikipedia.org/wiki/Make) is properly installed if you intend to use the provided `Makefile`.

### Makefile

The project comes with a convenient `Makefile` that helps you to build, install, lint and test the project.

```text
$ make <target>

Targets:
  Lint:
    lint                Lint all available linters
    lint-go             Lint go source code
    lint-proto          Lint proto files
  Format:
    format              Run all available formatters
    format-go           Format go files
    format-proto        Format proto files
  Build:
    build               Build all available artefacts (executable, docker image, etc.)
    build-go            Build node executable for the current environment (default build)
    build-go-all        Build node executables for all available environments
  Install:
    install             Install node executable
  Test:
    test                Pass all the tests
    test-go             Pass the test for the go source code
  Chain:
    chain-init          Initialize the blockchain with default settings.
    chain-start         Start the blockchain with existing configuration (see chain-init)
    chain-stop          Stop the blockchain
    chain-upgrade       Test the chain upgrade from the given FROM_VERSION to the given TO_VERSION
  Clean:
    clean               Remove all the files from the target folder
  Proto:
    proto               Generate all resources for proto files (go, doc, etc.)
    proto-format        Format Protobuf files
    proto-build         Build all Protobuf files
    proto-gen           Generate all the code from the Protobuf files
  Documentation:
    doc                 Generate all the documentation
    doc-proto           Generate the documentation from the Protobuf files
    doc-command         Generate markdown documentation for the command
  Mock:
    mock                Generate all the mocks (for tests)
  Release:
    release-assets      Generate release assets
  Help:
    help                Show this help.

This Makefile depends on docker. To install it, please follow the instructions:
- for macOS: https://docs.docker.com/docker-for-mac/install/
- for Windows: https://docs.docker.com/docker-for-windows/install/
- for Linux: https://docs.docker.com/engine/install/
```

### Build

To build the `axoned` node, invoke the goal `build` of the `Makefile`:

```sh
make build
```

The binary will be generated under the folder `target/dist`.

## Bug reports & feature requests

If you notice anything not behaving how you expected, if you would like to make a suggestion or would like
to request a new feature, please open a [**new issue**](https://github.com/axone-protocol/axoned/issues/new/choose). We appreciate any help
you're willing to give!

> Don't hesitate to ask if you are having trouble setting up your project repository, creating your first branch or
> configuring your development environment. Mentors and maintainers are here to help!

## Community

The [**AXONE Discord Server**](https://discord.gg/axone) is our primary chat channel for the open-source community,
software developers and node operators.

Please reach out to us and say hi 👋, we're happy to help there.

[Cosmos SDK]: https://v1.cosmos.network/sdk
[Docker]: https://www.docker.com/
[Go]: https://go.dev
[Tendermint]: https://tendermint.com/

## You want to get involved? 😍

So you want to contribute? Great! ❤️ We appreciate any help you're willing to give. Don't hesitate to open issues and/or
submit pull requests.

We believe that collaboration is key to the success of the AXONE project. Join our Community discussions on the [Community Repository](https://github.com/axone-protocol/community) to:

- Engage in conversations with peers and experts.
- Share your insights and experiences with AXONE.
- Learn from others and expand your knowledge of the protocol.

The Community Repository serves as a hub for discussions, questions, and knowledge-sharing related to AXONE. We encourage you to actively participate and contribute to the growth of our community.

Please check out AXONE health files:

- [Contributing](https://github.com/axone-protocol/.github/blob/main/CONTRIBUTING.md)
- [Code of conduct](https://github.com/axone-protocol/.github/blob/main/CODE_OF_CONDUCT.md)
