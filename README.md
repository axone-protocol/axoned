[![axone github banner](https://raw.githubusercontent.com/axone-protocol/.github/main/profile/static/axone-banner.png)](https://axone.xyz)

<p align="center">
  <a href="https://discord.gg/axone"><img src="https://img.shields.io/discord/946759919678406696.svg?label=discord&labelColor=7289DA&logo=discord&logoColor=white&color=gray&style=for-the-badge" /></a> &nbsp;
  <a href="https://www.linkedin.com/company/axone-protocol/"><img src="https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white" /></a> &nbsp;
  <a href="https://twitter.com/axonexyz"><img src="https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white" /></a> &nbsp;
  <a href="https://blog.axone.xyz"><img src="https://img.shields.io/badge/Medium-12100E?style=for-the-badge&logo=medium&logoColor=white" /></a> &nbsp;
  <a href="https://www.youtube.com/channel/UCiOfcTaUyv2Szv4OQIepIvg"><img src="https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white" /></a>
</p>

# Axone - Orchestration Layer for AI

[![lint](https://img.shields.io/github/actions/workflow/status/axone-protocol/axoned/lint.yml?label=lint&style=for-the-badge&logo=github)](https://github.com/axone-protocol/axoned/actions/workflows/lint.yml)
[![build](https://img.shields.io/github/actions/workflow/status/axone-protocol/axoned/build.yml?label=build&style=for-the-badge&logo=github)](https://github.com/axone-protocol/axoned/actions/workflows/build.yml)
[![test](https://img.shields.io/github/actions/workflow/status/axone-protocol/axoned/test.yml?label=test&style=for-the-badge&logo=github)](https://github.com/axone-protocol/axoned/actions/workflows/test.yml)
[![codecov](https://img.shields.io/codecov/c/github/axone-protocol/axoned?style=for-the-badge&token=O3FJO5QDCA&logo=codecov)](https://codecov.io/gh/axone-protocol/axoned)
[![Go Report Card](https://goreportcard.com/badge/github.com/axone-protocol/axoned/v13?style=for-the-badge)](https://goreportcard.com/report/github.com/axone-protocol/axoned/v13)
[![docker-pull](https://img.shields.io/docker/pulls/axoneprotocol/axoned?label=downloads&style=for-the-badge&logo=docker)](https://hub.docker.com/r/axoneprotocol/axoned)
[![Godoc Reference](https://img.shields.io/badge/godoc-reference-blue.svg?logo=go&logoColor=white&labelColor=gray&label=&style=for-the-badge)](https://pkg.go.dev/github.com/axone-protocol/axoned/v13)

[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge&logo=conventionalcommits)](https://conventionalcommits.org)
[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg?style=for-the-badge)](https://github.com/semantic-release/semantic-release)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg?style=for-the-badge)](https://github.com/axone-protocol/.github/blob/main/CODE_OF_CONDUCT.md)
[![license](https://img.shields.io/github/license/axone-protocol/axoned.svg?label=License&style=for-the-badge)](https://opensource.org/licenses/Apache-2.0)

> `Axone` is a public [dPoS](https://en.bitcoinwiki.org/wiki/DPoS) layer 1 specifically designed for connecting, sharing, and monetizing any resources in the AI stack. It is an open network dedicated to collaborative AI workflow management that is universally compatible with any data, model, or infrastructure. Data, algorithms, storage, compute, APIs... Anything on-chain and off-chain can be shared.

## The protocol

`axoned` is the node of the [AXONE](https://axone.xyz) network built on the [Cosmos SDK] 💫 & [Tendermint] consensus, which allows companies & individuals to define on-chain rules, share any off-chain resources & create a new generation of applications on top of them.

- 📖 Read the [introduction blog post](https://blog.axone.xyz/) to understand the project’s vision.
- 🧠 Explore the [white paper](https://docs.axone.xyz/docs/whitepaper/abstract) for a deeper look at the protocol architecture and network economics.

Looking for more?

- If you’re searching for Axone' smart contracts → [github.com/axone-protocol/contracts](https://github.com/axone-protocol/contracts)
- For the Axone SDK → [github.com/axone-protocol/axone-sdk](https://github.com/axone-protocol/axone-sdk)
- For the MCP server implementation → [github.com/axone-protocol/axone-mcp](https://github.com/axone-protocol/axone-mcp)

## Want to become a validator?

Validators are responsible for securing the axone network. Validator responsibilities include maintaining a functional [node](https://docs.axone.xyz/docs/nodes/run-node) with constant uptime and providing a sufficient amount of $AXONE as stake. In exchange for this service, validators receive block rewards and transaction fees.

Want to become a validator? 👉 [Checkout the documentation!](https://docs.axone.xyz/docs/nodes/introduction)

Looking for a network to join ? 👉 [Checkout the networks!](https://github.com/axone-protocol/networks)

## Supported platforms

The `axoned` blockchain currently supports the following builds:

| **Platform** | **Arch** |       **Status**       |
| ------------ | -------- | :--------------------: |
| Darwin       | amd64    |           ✅           |
| Darwin       | arm64    |           ✅           |
| Linux        | amd64    |           ✅           |
| Linux        | arm64    |           ✅           |
| Windows      | amd64    | ️🚫<br/> Not supported |

> Note: as the blockchain depends on [CosmWasm/wasmvm](https://github.com/CosmWasm/wasmvm), we only support the targets
> supported by this project.

## Releases

All releases can be found [here](https://github.com/axone-protocol/axoned/releases).

`axoned` follows the [Semantic Versioning 2.0.0](https://semver.org/) to determine when and how the version changes, and
we also apply the philosophical principles of [release early - release often](https://en.wikipedia.org/wiki/Release_early,_release_often).

## Install

### From release

```sh
curl https://i.jpillora.com/axone-protocol/axoned! | bash
```

### From source

```sh
make install
```

### Using docker

```sh
docker run -ti --rm axoneprotocol/axoned --help
```

## Developing & contributing

`axoned` is written in [Go] and built using [Cosmos SDK]. A number of smart contracts are also deployed on the
AXONE blockchain and hosted in the [axone-protocol/contracts](https://github.com/axone-protocol/contracts) project.

### Prerequisites

- install [Go] `1.24+` following instructions from the [official Go documentation](https://golang.org/doc/install);
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
    lint-go-golangci    Lint go source code with golangci-lint
    lint-go-modernize   Lint go source code with modernize
    lint-proto          Lint proto files
  Format:
    format              Run all available formatters
    format-go           Format go files
    format-proto        Format proto files
  Build:
    build               Build all available artefacts (executable, docker image, etc.)
    build-go            Build node executable for the current environment (default build)
    build-go-all        Build node executables for all available environments
    build-docker        Build docker image
  Install:
    install             Install node executable
  Test:
    test                Pass all the tests
    test-go             Pass the test for the go source code
  Chain:
    chain-init          Initialize the blockchain with default settings.
    chain-start         Start the blockchain with existing configuration (see chain-init)
    chain-stop          Stop the blockchain
    chain-upgrade       Test the chain upgrade from the given FROM_VERSION to the given TO_VERSION.
  Clean:
    clean               Remove all the files from the target folder
  Proto:
    proto               Generate all resources for proto files (go, doc, etc.)
    proto-gen           Generate all the code from the Protobuf files
  Documentation:
    doc                 Generate all the documentation
    doc-proto           Generate the documentation from the Protobuf files
    doc-command         Generate markdown documentation for the command
    doc-predicate       Generate markdown documentation for all the predicates (module logic)
  Mock:
    mock                Generate all the mocks (for tests)
  Release:
    release-assets      Generate release assets
  Dependencies:
    deps                Install all the dependencies (tools, etc.)
    deps-tparse         Install tparse (v0.17.0)
    deps-heighliner     Install heighliner (v1.7.4)
    deps-cosmovisor     Install cosmovisor (v1.7.1)
    deps-golangci-lint  Install golangci-lint (v2.4.0)
    deps-modernize      Install modernize (v0.20.0)
  Help:
    help                Show this help.

This Makefile depends on docker. To install it, please follow the instructions:
- for macOS: https://docs.docker.com/docker-for-mac/install/
- for Windows: https://docs.docker.com/docker-for-windows/install/
- for Linux: https://docs.docker.com/engine/install/
```

### Build

To build the `axoned` node, invoke the goal `build-go` of the `Makefile`:

```sh
make build-go
```

The binary will be generated under the folder `target/dist`.

### Build a docker image

This project leverages [heighliner](https://github.com/strangelove-ventures/heighliner) to simplify the management and
creation of production-grade container images. To build a Docker image, use the `build-docker` target in the `Makefile`:

```sh
make build-docker
```

### Run a local network

To initialize a local network configuration, invoke the goal `chain-init` of the `Makefile`:

```sh
make chain-init
```

The node home directory will be generated under the folder `target/deployment/localnet`. The configuration contains a single validator node.

To start the network, invoke the goal `chain-start` of the `Makefile`:

```sh
make chain-start
```

A wallet is preconfigured with some tokens, you can use it as follows:

```sh
axoned --home target/deployment/localnet tx bank send validator [to_address] [amount]
```

## Bug reports & feature requests

If you notice anything not behaving how you expected, if you would like to make a suggestion or would like
to request a new feature, please open a [**new issue**](https://github.com/axone-protocol/axoned/issues/new/choose). We appreciate any help
you're willing to give!

> Don't hesitate to ask if you are having trouble setting up your project repository, creating your first branch or
> configuring your development environment. Mentors and maintainers are here to help!

## Audit

| Date       | Auditor                            | Version                                                                                                     | Report                                                                                                                                                                                    |
| ---------- | ---------------------------------- | ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 2024/08/22 | [BlockApex](https://blockapex.io/) | [2f0f84d (v10.0.0)](https://github.com/axone-protocol/axoned/tree/2f0f84d369852bdb178e299a88c1b8eeb0654b8e) | [Axone Blockchain - Final Audit Report.pdf](https://github.com/BlockApex/Audit-Reports/blob/15d8765ac45b4a83bb2f1446fc9bf869c123f8d2/Axone%20Blockchain%20-%20Final%20Audit%20Report.pdf) |

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

We believe that collaboration is key to the success of the Axone project. Join our Community discussions on the [Community space](https://github.com/orgs/axone-protocol/discussions) to:

- Engage in conversations with peers and experts.
- Share your insights and experiences with Axone.
- Learn from others and expand your knowledge of the protocol.

The Community space serves as a hub for discussions, questions, and knowledge-sharing related to Axone.
We encourage you to actively participate and contribute to the growth of our community.

Please check out Axone health files:

- [Contributing](https://github.com/axone-protocol/.github/blob/main/CONTRIBUTING.md)
- [Code of conduct](https://github.com/axone-protocol/.github/blob/main/CODE_OF_CONDUCT.md)
