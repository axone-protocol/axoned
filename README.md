![okp4 github banner](./docs/okp4-banner.png)

<p align="center">
  <a href="https://discord.gg/GHNZh4SaJ3"><img src="https://img.shields.io/badge/Discord-7289DA?style=for-the-badge&logo=discord&logoColor=white" /></a> &nbsp;
  <a href="https://www.linkedin.com/company/okp4-open-knowledge-protocol-for"><img src="https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white" /></a> &nbsp;
  <a href="https://twitter.com/OKP4_Protocol"><img src="https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white" /></a> &nbsp;
  <a href="https://medium.com/okp4"><img src="https://img.shields.io/badge/Medium-12100E?style=for-the-badge&logo=medium&logoColor=white" /></a> &nbsp;
  <a href="https://www.youtube.com/channel/UCiOfcTaUyv2Szv4OQIepIvg"><img src="https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white" /></a>
</p>

# 𝒐𝒌𝒑4𝒅

[![lint](https://img.shields.io/github/workflow/status/okp4/okp4d/Lint?label=lint&style=for-the-badge)](https://github.com/okp4/okp4d/actions/workflows/lint.yml) [![build](https://img.shields.io/github/workflow/status/okp4/okp4d/Build?label=build&style=for-the-badge)](https://github.com/okp4/okp4d/actions/workflows/build.yml) [![test](https://img.shields.io/github/workflow/status/okp4/okp4d/Test?label=test&style=for-the-badge)](https://github.com/okp4/okp4d/actions/workflows/test.yml) [![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org) [![license](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg?style=for-the-badge)](https://opensource.org/licenses/BSD-3-Clause)

> `OKP4` is a public [dPoS](https://en.bitcoinwiki.org/wiki/DPoS) layer 1 specifically designed to enable communities to trustlessly share data, algorithms and resources to build the Dataverse - An open world where everybody can create or participate in custom ecosystems (with common governance mechanisms, sharing rules, business models...) to build a new generation of dApps way beyond Decentralized Finance.

## The protocol

`okp4d` is the node of the OKP4 network built on the [Cosmos SDK] 💫 & [Tendermint] consensus, and designed to become a hub of incentivized data providers, developers, data scientists & users collaborating to generate value from data and algorithms.

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

[Cosmos SDK]: https://v1.cosmos.network/sdk
[Docker]: https://www.docker.com/
[Go]: https://go.dev
[Tendermint]: https://tendermint.com/
