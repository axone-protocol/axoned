# ØKP4 protocol changelog

## [2.1.1](https://github.com/okp4/okp4d/compare/v2.1.0...v2.1.1) (2022-10-10)


### Bug Fixes

* **mint:** provide annual inflation rate ([608af3f](https://github.com/okp4/okp4d/commit/608af3f20ccc0afedfd518b56ae84c679df21502))
* **mint:** set initial inflation to 7.5% instead of 15% ([7bbd048](https://github.com/okp4/okp4d/commit/7bbd04802ce402967a56550f9c843c608fd91877))

## [2.1.0](https://github.com/okp4/okp4d/compare/v2.0.0...v2.1.0) (2022-10-05)


### Features

* **cliff:** add cliff cmd on vesting transaction ([ccff37c](https://github.com/okp4/okp4d/commit/ccff37ce34941297fc924b277381d08c7a48a367))
* **cliff:** add vesting-cliff-time flags on add-genesis account cmd ([ea3e2c5](https://github.com/okp4/okp4d/commit/ea3e2c5abb28ae24fa7f6633c3a3f5f95196f92f))
* **cliff:** override add-genesis-account ([434d418](https://github.com/okp4/okp4d/commit/434d4181a114675cd3c894378d0179266ec1fd3b))
* **cliff:** register cliff vesting account msg ([9106919](https://github.com/okp4/okp4d/commit/910691901185a32a2699c67fe5560f55e37fddba))
* implment okp4 inflaction calculation function ([2e95801](https://github.com/okp4/okp4d/commit/2e958010bdbaa78af0f0a1958c9c0493c2fbfed7))
* use okp4 inflation calculation fn (instead of default one) ([bdca893](https://github.com/okp4/okp4d/commit/bdca89319b5546fca1aabe2ff5bece1c05258558))
* use okp4 vesting module ([7493de9](https://github.com/okp4/okp4d/commit/7493de90a0fb601b9ab7b10eb0526a59652d5eef))
* use third party to generate proto ([cb4f5bb](https://github.com/okp4/okp4d/commit/cb4f5bb36fee1cc54d59bf902eae229fe2421bd6))


### Bug Fixes

* **cliff:** improve verification on cliff msg tx ([4dfaa5b](https://github.com/okp4/okp4d/commit/4dfaa5b2624fcd39b12459d83275027f8d09507d))
* **ibc:** ensure ibc fees are managed ([c26b0db](https://github.com/okp4/okp4d/commit/c26b0db71dbd5b71788a7f333b673bc75ac80198))
* make linter happy ([584851c](https://github.com/okp4/okp4d/commit/584851c5b248112fc10ea3ff0a3f3873a3878a06))

## [2.0.0](https://github.com/okp4/okp4d/compare/v1.3.0...v2.0.0) (2022-09-23)


### ⚠ BREAKING CHANGES

* reboot chain with ignite cli v0.24.0

### Features

* add logic module params to genesis files ([9ac7ef8](https://github.com/okp4/okp4d/commit/9ac7ef8a7ecbe72fd8efd983ed138fd375aaf1fa))
* scaffold logic module using ignite ([81ee269](https://github.com/okp4/okp4d/commit/81ee26997d3e941fa24655ea89ce3dad93b6cfd2))
* update openapi documentation (synced with code) ([787ff01](https://github.com/okp4/okp4d/commit/787ff01148f1721f75594824c6da715ade5f18ea))
* **wasm:** prepare ante handler with wasm decorators ([afb4748](https://github.com/okp4/okp4d/commit/afb4748f5ea7b3b6723226235d9d93e171605c9f))
* **wasm:** wire wasm module in app ([b163790](https://github.com/okp4/okp4d/commit/b163790cb60c2221793742eed6fecc22d4254315))


### Bug Fixes

* fix (pre-)genesis files after 0.46 cosmos sdk migration ([da284a0](https://github.com/okp4/okp4d/commit/da284a0d4993f2231794ddbfae16b8151b23fbfc))
* use proper versions of buf protoc plugins ([1ca5e1d](https://github.com/okp4/okp4d/commit/1ca5e1d6127d255cf8f11375a7256e5c07efedf0))
* **workflow:** use secret for dockerhub user ([0d95c94](https://github.com/okp4/okp4d/commit/0d95c945e7c0e858ae22a2e7d573853e889d4473))


### Code Refactoring

* reboot chain with ignite cli v0.24.0 ([423179e](https://github.com/okp4/okp4d/commit/423179e5028de57b7859bbd6ebcb5d12a4b42fb5))

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
