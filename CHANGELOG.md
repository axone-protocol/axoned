# ØKP4 protocol changelog

# [1.3.0](https://github.com/okp4/okp4d/compare/v1.2.0...v1.3.0) (2022-07-08)


### Bug Fixes

* fix genesis file after cosmos SDK update ([ff03ba9](https://github.com/okp4/okp4d/commit/ff03ba906313f8cbdeb77b3fb158cdf188f197fb))
* generate pre-genesis file with uknow unit ([677fc0c](https://github.com/okp4/okp4d/commit/677fc0cd16f18f9e650d2df4a669932da9e8200c))
* make it start 🚀 ([0cdd4db](https://github.com/okp4/okp4d/commit/0cdd4db2bdd16fcbe2788d5b6802d7ed9402a2bb))
* make linters happy ([309193a](https://github.com/okp4/okp4d/commit/309193aae612a94bc8d1ee6f557af204a7720044))
* references all modules in SetOrder* functions ([84f9fe2](https://github.com/okp4/okp4d/commit/84f9fe201d0214b32a0ff3d2f4e4893bf879ba49))


### Features

* add a proper description for the OKP4 CLI ([bc74f2c](https://github.com/okp4/okp4d/commit/bc74f2c4f158895cf45534b33e0c32a2fb92cd90))
* **denom:** add uknow & know denoms metadata ([55d52ef](https://github.com/okp4/okp4d/commit/55d52efceff6cbd9714dc7a6014029824716f138))
* handle wasm proposals in app ([44b10c0](https://github.com/okp4/okp4d/commit/44b10c0b91cf10ff8b17b86737067ff1ee47586a))
* implement genesis account cmd ([0368a11](https://github.com/okp4/okp4d/commit/0368a111f374de99c6e3cb0a2ecc1d8c641e7187))
* implement genesis wasm cmd ([b4ac0bc](https://github.com/okp4/okp4d/commit/b4ac0bc82b3cdf3c6b7c0f24cd5739eb7f245ae8))
* implement okp4 encoding config ([acf6a26](https://github.com/okp4/okp4d/commit/acf6a26fc5d0d42df5e77b3bca65f4e663a87200))
* implement root cmd ([d065735](https://github.com/okp4/okp4d/commit/d0657358ee65dbd1c516e81eedb56a5ed6c63fd9))
* prepare ante handler with wasm decorators ([fb32135](https://github.com/okp4/okp4d/commit/fb32135288dcff37b6a0a25897210ab78461b745))
* provide default app encoding config ([602d2db](https://github.com/okp4/okp4d/commit/602d2db7dad855871d550d849b964e35ca2f26a1))
* re-sync openapi specification (after monitoringp removal) ([96dba69](https://github.com/okp4/okp4d/commit/96dba6985cbce0def82cb56671a265a915bb3116))
* remove ignite dependencies ([8cf0c4b](https://github.com/okp4/okp4d/commit/8cf0c4b2a9314f42b244193e63de5b3d90ab8a42))
* remove monitoringp module ([f97b8b9](https://github.com/okp4/okp4d/commit/f97b8b94895a7aa6f66bfce6880d263b67cf8cec))
* remove unwanted monitoringp module ([6be3175](https://github.com/okp4/okp4d/commit/6be317582c22ea5172932126e33f9efefc725d66))
* replace ignite root cmd by ours ([e92cfe4](https://github.com/okp4/okp4d/commit/e92cfe43d94243ddc1faabdd27f13bf401827297))
* wire wasm module in app ([ac56cbe](https://github.com/okp4/okp4d/commit/ac56cbe76bb5f16c75c63cfb1675450217ab23a2))

# [1.2.0](https://github.com/okp4/okp4d/compare/v1.1.0...v1.2.0) (2022-04-21)


### Features

* check uri format on trigger service msg ([b9ae987](https://github.com/okp4/okp4d/commit/b9ae987b081ee7b2dca5da09cabe088d19c2a728))
* extend debug cmd adding decode-blocks ([2959444](https://github.com/okp4/okp4d/commit/29594444e8c749f975e28316dddb0b1322bfebfa))
* implements debug blocks proto base64 decode ([1749e70](https://github.com/okp4/okp4d/commit/1749e70151ae16b0abfe326b70501ab64ca42d2a))
* scaffold trigger-service message ([57a54eb](https://github.com/okp4/okp4d/commit/57a54eb49cff5ee66957a1d2d795af5edab45bde))

# [1.1.0](https://github.com/okp4/okp4d/compare/v1.0.0...v1.1.0) (2022-03-16)


### Bug Fixes

* fix missing `id` in message ([c8f81fe](https://github.com/okp4/okp4d/commit/c8f81febc74622397ffb2b2a4e103ef28511bb92))


### Features

* implement dataspace creation ([8ec1074](https://github.com/okp4/okp4d/commit/8ec10742fd6271720f98fce93e20626447008841))
* scaffold module knowledge ([82c95e7](https://github.com/okp4/okp4d/commit/82c95e7acebf30df7ca924383ad8e5de3da533bc))
* update openapi documentation ([03a541c](https://github.com/okp4/okp4d/commit/03a541c732f8fd94701015317c5d1fcb0574c688))

# 1.0.0 (2022-02-08)


### Bug Fixes

* Add make init and fix start ([30468fe](https://github.com/okp4/okp4d/commit/30468fed2d6b8b86e60f0cb44479547b05bded18))


### Code Refactoring

* rename repository from "okp4" to "okp4d" ([1a3b24f](https://github.com/okp4/okp4d/commit/1a3b24f383e7ef18f0bd2a2e93b774813e824339))


### Features

* add scaffolded code ([232da23](https://github.com/okp4/okp4d/commit/232da23ee3428c3aec0d68a41771832aa6502cb3))


### BREAKING CHANGES

* this change has a few implications, such as changing
the name of the published docker image.
