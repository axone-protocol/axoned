# AXONE protocol changelog

## [13.0.1](https://github.com/axone-protocol/axoned/compare/v13.0.0...v13.0.1) (2026-02-01)

## [13.0.0](https://github.com/axone-protocol/axoned/compare/v12.0.0...v13.0.0) (2025-09-27)


### ⚠ BREAKING CHANGES

* **deps:** bump wasmd to v0.61.2 and wasmvm to v3.0.2
* **cli:** remove custom credential and did key commands
* **app:** drop x/crisis module and keeper
* **deps:** bump cosmos-sdk to v0.53.4 and ibc-go to v8.7.0
* **vesting:** back to canonical cosmos vesting
* **ibc:** replug IBC modules right now 🔌

### Code Refactoring

* **app:** drop x/crisis module and keeper ([78e57f2](https://github.com/axone-protocol/axoned/commit/78e57f2043e8243cfda9050d164739ba838f8df3))
* **cli:** remove custom credential and did key commands ([9906178](https://github.com/axone-protocol/axoned/commit/9906178e659e9d003fadfc9e265081d605acacab))
* **ibc:** replug IBC modules right now 🔌 ([569348b](https://github.com/axone-protocol/axoned/commit/569348be8dcc87333bde26fa381eae6e49b75255))
* **vesting:** back to canonical cosmos vesting ([32e117c](https://github.com/axone-protocol/axoned/commit/32e117c7a20b57c1d2a12e2328dda1732ac0552f))


### Build System

* **deps:** bump cosmos-sdk to v0.53.4 and ibc-go to v8.7.0 ([efd9426](https://github.com/axone-protocol/axoned/commit/efd942680dcaf0bc2fb0ba3b3ffbefb467b697fa))
* **deps:** bump wasmd to v0.61.2 and wasmvm to v3.0.2 ([1a6f327](https://github.com/axone-protocol/axoned/commit/1a6f32791a6d31968b4f0ac53660772457ce40b7))

## [12.0.0](https://github.com/axone-protocol/axoned/compare/v11.0.1...v12.0.0) (2025-07-06)


### ⚠ BREAKING CHANGES

* **ibc:** unplug IBC for now 🔌

### Features

* **app:** add MaxWasmSize parameter for configurable Wasm code limits ([caa0403](https://github.com/axone-protocol/axoned/commit/caa0403b192d7ed3e2742440e6985c335776c5c3))
* **logic:** add base64_encoded/3 precicate ([d1a75ee](https://github.com/axone-protocol/axoned/commit/d1a75eea1197a40cb7ae55d73799d1ba8f4b4419))
* **logic:** add base64/2 predicate ([7e55d19](https://github.com/axone-protocol/axoned/commit/7e55d1942bb5db1da8152dde6515166f8adc8e4b))
* **logic:** add base64url/2 predicate ([3176ff7](https://github.com/axone-protocol/axoned/commit/3176ff7b2100e7007ddb31b9ce11c70f988075f0))


### Code Refactoring

* **ibc:** unplug IBC for now 🔌 ([324685f](https://github.com/axone-protocol/axoned/commit/324685f9c6962d5f66799f4df751b29e99360675))

## [11.0.1](https://github.com/axone-protocol/axoned/compare/v11.0.0...v11.0.1) (2024-12-18)

### Build

* **project:** update -tags format to use comma-separated values ([f4dad40](https://github.com/axone-protocol/axoned/commit/f4dad40660c391b5ec9b1aa524e4270f84547648))

### Docs

* **cli:** update CLI docs following Cosmos SDK v0.50.11 upgrade ([0ddeec2](https://github.com/axone-protocol/axoned/commit/0ddeec2d149084173dd5355338a6323573fd8409))

## [11.0.0](https://github.com/axone-protocol/axoned/compare/v10.0.0...v11.0.0) (2024-12-11)


### ⚠ BREAKING CHANGES

* **logic:** change Limits and GasPolicy properties to use uint64
* **logic:** use '=' functor for encoding key-value pairs in json_prolog/2
* **logic:** remove @([]) empty list encoding in json_prolog/2
* **logic:** preserve object key order in json_prolog/2
* **logic:** add decimal number support for json_prolog/2

### Features

* **logic:** accept JSON as text for json_prolog/2 ([83aa22c](https://github.com/axone-protocol/axoned/commit/83aa22c11fbdae76121697175da3dfccd5dfe04d))
* **logic:** add atomic_list_concat/2 predicate ([44daf95](https://github.com/axone-protocol/axoned/commit/44daf95f7c082d5a941fc919e63f761343f78bf6))
* **logic:** add atomic_list_concat/3 predicate ([f394074](https://github.com/axone-protocol/axoned/commit/f394074d8ec4802381445a825541688a59b7e575))
* **logic:** add block_header/1 predicate ([94a94a3](https://github.com/axone-protocol/axoned/commit/94a94a374cb186d7ad2312b9ffc72bcbbc82b8e0))
* **logic:** add json_read/2 predicate ([a9902a3](https://github.com/axone-protocol/axoned/commit/a9902a32a7e4c180e55609cd445fc6fe24823047))
* **logic:** add json_write/2 predicate ([1675c34](https://github.com/axone-protocol/axoned/commit/1675c34b3b10b3a3898a1fc2b59765e00b8bb7eb))
* **logic:** add term_to_atom/2 predicate ([96043c2](https://github.com/axone-protocol/axoned/commit/96043c23fb4943825c1f86d57feb188904a4a9b6))
* **logic:** deprecate the block_height/1 predicate ([0467fb5](https://github.com/axone-protocol/axoned/commit/0467fb5237abc705d023b4d3d3ec789c37351899))
* **logic:** deprecate the block_time/1 predicate ([404fdc2](https://github.com/axone-protocol/axoned/commit/404fdc2691d1fc5732ede10d84fd1ec2a5ea0884))
* **logic:** deprecate the chain_id/1 predicate ([8a6144e](https://github.com/axone-protocol/axoned/commit/8a6144e66cf7d3ca96e1c019ffbdd041470e5f0c))
* **logic:** generate default values for the interpreter in genesis file ([6146f20](https://github.com/axone-protocol/axoned/commit/6146f2007528092db5f4c89b9e9702ebb9e980b2))
* **logic:** implement chain migration to v11 ([da93e81](https://github.com/axone-protocol/axoned/commit/da93e8146a2fe127d43a5f1f886e78b1b077b890))
* **logic:** implement telemetry for predicate execution duration metrics ([f550a53](https://github.com/axone-protocol/axoned/commit/f550a535dbb0327876c64b6d0e93eb579df70af2))
* **logic:** introduce resource_context/1 composite for missing context value retrieval ([1132405](https://github.com/axone-protocol/axoned/commit/1132405069638be018c7792900862efd2e10f9f0))
* **logic:** introduce telemetry to track predicate call counts ([d0fd0c7](https://github.com/axone-protocol/axoned/commit/d0fd0c7bed42205c70b613b7de0838c4787e232b))
* **logic:** use the default cost params when lookup gas cost ([350b78a](https://github.com/axone-protocol/axoned/commit/350b78a087c2ce0dbd20e7efa1e293506ff0fd8f))


### Bug Fixes

* **logic:** add decimal number support for json_prolog/2 ([4e7375b](https://github.com/axone-protocol/axoned/commit/4e7375b4bb3714ef0ed706f7628fe030847a0d2d))
* **prolog:** implement fine-grained gas metering for predicates ([239a4e0](https://github.com/axone-protocol/axoned/commit/239a4e04074b518d764644b0649e85827a0365dc))


### Code Refactoring

* **logic:** change Limits and GasPolicy properties to use uint64 ([db18c9a](https://github.com/axone-protocol/axoned/commit/db18c9a38d80c4e882f410fafca876679a8cffaa))
* **logic:** preserve object key order in json_prolog/2 ([6ea2633](https://github.com/axone-protocol/axoned/commit/6ea26335aadb9cd62bda9ea1653c43bde05c9a9c))
* **logic:** remove @([]) empty list encoding in json_prolog/2 ([c3e5781](https://github.com/axone-protocol/axoned/commit/c3e57816bab683bb29b3f160f16d566b355ec84f))
* **logic:** use '=' functor for encoding key-value pairs in json_prolog/2 ([37d4518](https://github.com/axone-protocol/axoned/commit/37d4518c9d6ec6693ef4632af67c0c1b74c00d3a))

## [10.0.0](https://github.com/axone-protocol/axoned/compare/v9.0.0...v10.0.0) (2024-08-29)


### ⚠ BREAKING CHANGES

* **deps:** bump docker image wasmvm
* **deps:** bump wasmd to v0.53.0

### Features

* **credential:** allow choose proof purpose when signing ([b30ff57](https://github.com/axone-protocol/axoned/commit/b30ff57c5791b2cbdd41f141b1d631fa5a7905b3))


### Bug Fixes

* **credential:** correctly wraps error on parsing arg ([643dab4](https://github.com/axone-protocol/axoned/commit/643dab4ed6cdf8bc53ca2e2754ca74f39acfda25))
* **deps:** bump docker image wasmvm ([3fb333c](https://github.com/axone-protocol/axoned/commit/3fb333cc89834da0cc8932edd5670460b14e15de))
* **deps:** bump wasmd to v0.53.0 ([1283233](https://github.com/axone-protocol/axoned/commit/1283233eb6e02300b775275b0c2c38f2d18c7bc8))

## [9.0.0](https://github.com/axone-protocol/axoned/compare/v8.0.0...v9.0.0) (2024-08-02)


### ⚠ BREAKING CHANGES

* **logic:** remove unsafe prolog predicates from registry
* **logic:** add max_variables limits params
* **logic:** change interpreter by our fork axone-protocol/prolog
* **logic:** remove max gas module parameter
* **wasm:** bump wasm to 0.52.0

### Features

* **app:** update application init ([cee814a](https://github.com/axone-protocol/axoned/commit/cee814a391f1d1ec6fa6a73408bf90e9bc831d84))
* handle MaxVariables error limit exceeded ([040a765](https://github.com/axone-protocol/axoned/commit/040a765ea934eb5c5dd0fa466758bd0385234a7e))
* **logic:** add max_variables limits params ([94e54e0](https://github.com/axone-protocol/axoned/commit/94e54e04ee22c73d7fbf4fcddf20097b9b727816))
* **logic:** mark rpc queries as safe ([66b1b8c](https://github.com/axone-protocol/axoned/commit/66b1b8cc4b4c35d94115b8fd80f7c3418de9ca9a))
* **vesting:** update to align with cosmos-sdk v0.50.4 ([7536aee](https://github.com/axone-protocol/axoned/commit/7536aee0d8f84e5dfb11d26d3819289c45a2013f))
* **wasm:** bump wasm to 0.52.0 ([c1e404c](https://github.com/axone-protocol/axoned/commit/c1e404c7ebd6684a9cb048046c905e1d71f399de))
* **wasm:** use wasm provided build in capabilities ([77c1063](https://github.com/axone-protocol/axoned/commit/77c10635a901963768f6f2e55f011994235cfd2f))


### Bug Fixes

* **docs:** generate docs with module subfolder on proto ([fde166c](https://github.com/axone-protocol/axoned/commit/fde166cda8dec105b958a4f2ba4b5f1a15284b34))
* **gov:** properly configure legacy props router ([74df1c7](https://github.com/axone-protocol/axoned/commit/74df1c73637329be6af5ec7e637321b01d0201a1))
* **logic:** handle and return out of gas error at GRPC level ([ff3de76](https://github.com/axone-protocol/axoned/commit/ff3de76a201811d44eb8ca80b80c892b9feed767))
* **logic:** remove max gas module parameter ([db9164b](https://github.com/axone-protocol/axoned/commit/db9164b7037e91c656d738d0f9f6b9c6f3027a3e))
* **logic:** remove unsafe prolog predicates from registry ([293da10](https://github.com/axone-protocol/axoned/commit/293da10e8c24ec1712c8c700330beb9e213c3819))
* **logic:** use ordered predicate register ([0c99e9c](https://github.com/axone-protocol/axoned/commit/0c99e9c3ae6d367f5870bab484694e031348516f))
* **proto:** remove unused proto def from global registry ([798fdcd](https://github.com/axone-protocol/axoned/commit/798fdcded13ee8c4336b408dabd0d490bfc6d95e))
* use new method for instantiate interpreter ([e08560d](https://github.com/axone-protocol/axoned/commit/e08560d02564c5d0fa6e21360703c8229e6897dc))
* **vesting:** format and lint ([c018447](https://github.com/axone-protocol/axoned/commit/c018447c95832720e419282c7695d0655765a775))
* **vesting:** reintroduce the add-genesis-account custom command to include cliff vesting account ([4b4bafb](https://github.com/axone-protocol/axoned/commit/4b4bafbc61a8d2ffb46fa1a4ca28eb35435b0d70))


### Performance Improvements

* **logic:** reduce read params in ask query ([92e6d05](https://github.com/axone-protocol/axoned/commit/92e6d054d06cd3096a199d3a1ba8d504ec4a3520))


### Code Refactoring

* **logic:** change interpreter by our fork axone-protocol/prolog ([c1b3035](https://github.com/axone-protocol/axoned/commit/c1b30352f084aeb6e91347d3c274f5a361bc6858))

## [8.0.0](https://github.com/axone-protocol/axoned/compare/v7.1.0...v8.0.0) (2024-05-22)


### ⚠ BREAKING CHANGES

* **renaming:** change denom name from uknow to uaxone

### Features

* **cli:** update narative of axone command description ([4ab2b38](https://github.com/axone-protocol/axoned/commit/4ab2b38bd1b868534d0e2856f8781c6421517b33))
* **logic:** introduce "limit" flag to command "ask" ([58b4c45](https://github.com/axone-protocol/axoned/commit/58b4c45855c21885e65292482fddfe55d91ebaa5))


### Code Refactoring

* **renaming:** change denom name from uknow to uaxone ([88b9c92](https://github.com/axone-protocol/axoned/commit/88b9c92c6cfb559470f8b6bace299d8272c05bc4))

# ØKP4 protocol changelog

## [7.1.0](https://github.com/axone-protocol/axoned/compare/v7.0.1...v7.1.0) (2024-04-02)


### Features

* **logic:** expose open/3 predicate ([b71bf4f](https://github.com/axone-protocol/axoned/commit/b71bf4f3089239da24df5d52f055f354e58e0159))
* **logic:** introduce parameter "limit" for queries ([cc86031](https://github.com/axone-protocol/axoned/commit/cc860317cd40cc8462e462267e31c527b0e09d25))
* update generated proto & docs ([2d833c8](https://github.com/axone-protocol/axoned/commit/2d833c8164efb627e347d8d7d497d4c4ea42995a))
* **wasm:** update cosmwasm to v2 ([5640041](https://github.com/axone-protocol/axoned/commit/564004197d555855935d2aebcb61851b083bf7c2))


### Bug Fixes

* **logic:** make the user output stream always initialized ([2345934](https://github.com/axone-protocol/axoned/commit/2345934f62c6519a951031a71768390f681f7cc6))
* **logic:** prevent non-bound substitution variables to be returned ([9250bda](https://github.com/axone-protocol/axoned/commit/9250bda2ddca9119c9cb18014f3e2008f30c4e0c))
* **sign:** properly expand input path ([a783c02](https://github.com/axone-protocol/axoned/commit/a783c0228d34f0e4d2b21c2b457ab478b38fd667))
* **sign:** use right signature algorithm ([5bb6d2d](https://github.com/axone-protocol/axoned/commit/5bb6d2dced285b3519f012b840fef8192993193d))

## [7.0.1](https://github.com/axone-protocol/axoned/compare/v7.0.0...v7.0.1) (2024-03-05)

## [7.0.0](https://github.com/axone-protocol/axoned/compare/v6.0.0...v7.0.0) (2024-02-29)


### ⚠ BREAKING CHANGES

* **logic:** improve predicate call policy (blacklist, gas)
* **logic:** enhance answer responses for error cases
* **logic:** adopt new exposed error terms from ichiban/prolog
* **logic:** adopt unstructured expression for substitutions
* **logic:** did_components/2 now returns encoded components

### Features

* **app:** register v7.0.0 software upgrade ([4358459](https://github.com/axone-protocol/axoned/commit/4358459650b186effc5754da8c1a69b139e5c5a8))
* **cli:** add --date to credential sign command ([993bd68](https://github.com/axone-protocol/axoned/commit/993bd68c7d75a6cc8a4267695870fb75effec973))
* **cli:** add credential sign command ([01083a4](https://github.com/axone-protocol/axoned/commit/01083a4b460260acbd541c6ae839b35e2012ca1c))
* **cli:** add did output only to show key command ([7a6eee2](https://github.com/axone-protocol/axoned/commit/7a6eee24627390708f209f6eff36df1c03c157e3))
* **cli:** add did:key to keys list command ([30a5615](https://github.com/axone-protocol/axoned/commit/30a5615ee462c8c48e82a0a31bc1f6575b04181c))
* **cli:** add schema-map flag to credential sign command ([bdd93a3](https://github.com/axone-protocol/axoned/commit/bdd93a3a14865de1a6cf14abb5056ebe4aa63e7a))
* **cli:** extend keys show command to display did:key info ([631a636](https://github.com/axone-protocol/axoned/commit/631a636f6bebefc14b7cbce21d8bb85b29d24a9d))
* **cli:** introduce "keys did" command ([39c87b7](https://github.com/axone-protocol/axoned/commit/39c87b75a7b8a863076f87b4e449d9fe0537af40))
* **mint:** add optional inflation bounds params ([fd77aab](https://github.com/axone-protocol/axoned/commit/fd77aaba811536fd3002fc5ecdb2ff2ba64b7035))
* **wasm:** enable some cosmwasm features ([19f93e8](https://github.com/axone-protocol/axoned/commit/19f93e8e2c941867630445e6bb4f7a1955981c57))


### Bug Fixes

* **app:** ignore tink crypto proto from checks at app init ([aa36819](https://github.com/axone-protocol/axoned/commit/aa36819fb5f678e2807dd61f0adca28bd697d52d))
* **cli:** go backwards by remvoving logic query autocli ([46d4b66](https://github.com/axone-protocol/axoned/commit/46d4b6649670242ae2b30621b411b606e2cc497e))
* **logic:** fix variable resolution ([d1f0ee0](https://github.com/axone-protocol/axoned/commit/d1f0ee02bf825765aec68ce5df437bb6a783bbe0))
* **logic:** re-introduce descriptor for CLI (lost with migration) ([060972d](https://github.com/axone-protocol/axoned/commit/060972dd38ce692aa180f310410dde1ba0c8a9da))
* wire config options of transactions for autocli ([b95c89c](https://github.com/axone-protocol/axoned/commit/b95c89caa164d53afb87f7c1d7198b882f192513)), closes [#551](https://github.com/axone-protocol/axoned/issues/551)


### Code Refactoring

* **logic:** adopt new exposed error terms from ichiban/prolog ([fd4a231](https://github.com/axone-protocol/axoned/commit/fd4a231685a98a7c90f889a981179c6d25acc796))
* **logic:** adopt unstructured expression for substitutions ([50f9d7f](https://github.com/axone-protocol/axoned/commit/50f9d7fc8d81b4c2a2a4a0b7bd99e412d87f2f57))
* **logic:** did_components/2 now returns encoded components ([e6cd0fc](https://github.com/axone-protocol/axoned/commit/e6cd0fc9f94aeb46738352957224eb7882179f9d))
* **logic:** enhance answer responses for error cases ([16698c3](https://github.com/axone-protocol/axoned/commit/16698c31bc3e9cee130241e62e83a5f317706a92))
* **logic:** improve predicate call policy (blacklist, gas) ([17fdaed](https://github.com/axone-protocol/axoned/commit/17fdaede65ae7ccdb714624741336e2e9089fb54))

## [6.0.0](https://github.com/axone-protocol/axoned/compare/v5.0.0...v6.0.0) (2024-01-19)


### ⚠ BREAKING CHANGES

* remove upgrades
* **app:** update app for v0.50.x sdk changes
* **logic:** standardize uri_encoded/3 errors
* **logic:** add error field to anwser proto type
* **logic:** standardize predicates errors
* **logic:** standardize error messages for some predicates
* **logic:** remove sha_hash/2 predicate
* **DIDComponents:** restrict type-checking for expected atom components

### Features

* **app:** re implement app state export ([85362d4](https://github.com/axone-protocol/axoned/commit/85362d4851d68efef88b8bfcaab155e2798d7270))
* **app:** update app for v0.50.x sdk changes ([8a39efc](https://github.com/axone-protocol/axoned/commit/8a39efc27918d44a091b491920df5008ddf8041d))
* **app:** update old migrations ([ab068fc](https://github.com/axone-protocol/axoned/commit/ab068fc2039ff63c6e84e49aa0770e4f4f98ee19))
* **cmd:** update cmd wiring ([dca3eb2](https://github.com/axone-protocol/axoned/commit/dca3eb2103e751f23c0bd27e05d6af6296214a94))
* **logic:** add crypto_data_hash/3 predicate with SHA-256 support ([edb3d55](https://github.com/axone-protocol/axoned/commit/edb3d55efb280c46111edaa308f8f62ae815d5e1))
* **logic:** add MD5 algorithm to crypto_data_hash/3 ([7e17001](https://github.com/axone-protocol/axoned/commit/7e17001873be050c17d1f7be5e63ab89e6ba7672))
* **logic:** add Secp256k1 support ([d612819](https://github.com/axone-protocol/axoned/commit/d61281999cb19e67ee493823077f0fcbc105cb85))
* **logic:** add Secp256k1 support for ECDSAVerify predicate ([b9f930b](https://github.com/axone-protocol/axoned/commit/b9f930bab5627d4c8ad111d20c93538755ee10f6))
* **logic:** add SHA-512 algorithm to crypto_data_hash/3 ([93b2d59](https://github.com/axone-protocol/axoned/commit/93b2d59cbffcec95ee590f50e550fc1b30d3f47b))
* **logic:** add string_bytes/3 predicate ([2e2bac2](https://github.com/axone-protocol/axoned/commit/2e2bac26069fb7cf17d071c460948f1d4056539d))
* **logic:** add universal Hash function ([157ff21](https://github.com/axone-protocol/axoned/commit/157ff21a1491dda9e3034882a28d7373d54d5f57))
* **logic:** add utf8 encoding option for TermToBytes ([3b3ca89](https://github.com/axone-protocol/axoned/commit/3b3ca89a961ab7d5ee8dc4d7bda84e7ddc364aae))
* **logic:** deprecate sha_hash/2 predicate ([9be1b56](https://github.com/axone-protocol/axoned/commit/9be1b562d34e6a7c0c8c30bcbb4956357a29ec64))
* **logic:** handle encoding option ([d83b5d3](https://github.com/axone-protocol/axoned/commit/d83b5d37ca7d3167db33534ebc1a5b1e8b2fa563))
* **logic:** implement crypto signature verification ([a92ec86](https://github.com/axone-protocol/axoned/commit/a92ec86ca74594755547818a7ae9e95617222624))
* **logic:** manage ecda secp256k1 key for secp_verify/4 ([268e493](https://github.com/axone-protocol/axoned/commit/268e493f41a561e1578216e22dc9280f36d5ee14))
* **logic:** migrate logic module ([a2b5286](https://github.com/axone-protocol/axoned/commit/a2b5286729f0c95064783a9883528385e3d09eab))
* **logic:** reigster eddsa_verify/4 predicate ([2a19e91](https://github.com/axone-protocol/axoned/commit/2a19e9160871b7a631b53f096f84aecc6c80aa8a))
* **logic:** remove sha_hash/2 predicate ([605d0cd](https://github.com/axone-protocol/axoned/commit/605d0cd1ccc98210a2cc6a97501eba698abbbab0))
* **mint:** add new mint function calculation ([a2bc105](https://github.com/axone-protocol/axoned/commit/a2bc105de2898a23e743b33c51138a6387328677))
* **mint:** create migration for v3 ([3df4557](https://github.com/axone-protocol/axoned/commit/3df45574c96cca5b753e465bcea2991459c9b80a))
* **mint:** implement token model v2 ([473221d](https://github.com/axone-protocol/axoned/commit/473221d1fbd0d2c9b0a536b294cb52d5354f8458))
* **mint:** instantiate default new mint params ([6484ce4](https://github.com/axone-protocol/axoned/commit/6484ce4a6f9dc90baa8f11a54bf6c0704b06dea7))
* **mint:** migrate mint module ([ca000b2](https://github.com/axone-protocol/axoned/commit/ca000b29e68a0a9cf531c05c68eb37c4c62f9454))
* **mint:** migrate minter params by removing not used key ([e62215a](https://github.com/axone-protocol/axoned/commit/e62215aa4bdec783d047c0f17373b4e707d6bbd3))
* **mint:** new mint params ([90456d0](https://github.com/axone-protocol/axoned/commit/90456d035830368fb33008f1ae58d5bdcbb67c96))
* **mint:** register migration handler ([e38ae03](https://github.com/axone-protocol/axoned/commit/e38ae03fea72b3715593287766c08a7d68940d77))
* **mint:** remove legacy store migration ([83bfc00](https://github.com/axone-protocol/axoned/commit/83bfc0089e199755363c819e75503349b4b6d892))
* **mint:** specify token model v2 schemas ([36a965f](https://github.com/axone-protocol/axoned/commit/36a965f0c0a4adc8257f66fd1ab023edb8056d2b))
* **prolog:** add default value support to GetOption func ([9a2dee9](https://github.com/axone-protocol/axoned/commit/9a2dee961e76f0397ee021e1001ddd0d9e417d43))
* remove upgrades ([be741a3](https://github.com/axone-protocol/axoned/commit/be741a3672fa6c4fadbda843666c2ef0a322972e))
* **vesting:** update module ([96f970c](https://github.com/axone-protocol/axoned/commit/96f970cdf1483e2beefcd00653b27edf52208f88))
* **wasm:** upgrade lib wasmvm in docker build ([5c01946](https://github.com/axone-protocol/axoned/commit/5c019469da6e7d146615ec58679c2858524be56f))


### Bug Fixes

* **app:** do not try to migrate consensus params ([61c6ddf](https://github.com/axone-protocol/axoned/commit/61c6ddfe491abfa5c33825caa490cfb3a716a9ca))
* **build:** update chain-upgrade makefile task to make it work with latest version ([3ea9694](https://github.com/axone-protocol/axoned/commit/3ea9694788438a7256235ebdd804937c1bf3a7bb))
* **lint:** fix unused linter directive ([49c1202](https://github.com/axone-protocol/axoned/commit/49c1202ebd634a0c3ec31d9e50d70d03f75582c7))
* **logic:** add return statement for hash errors ([9996654](https://github.com/axone-protocol/axoned/commit/9996654f59416cd6f5fc284f9c142016a1c9139b))
* **logic:** fix incorrect handling of empty list ([25e476d](https://github.com/axone-protocol/axoned/commit/25e476d898563f87967fda0accdbb758631b2f80))
* **logic:** fix some typos in naming and description ([e966ad6](https://github.com/axone-protocol/axoned/commit/e966ad6116a4dc583457c637cbf31f6d0baaa967))
* **mint:** calculate inflation in percent instead of permille ([c555990](https://github.com/axone-protocol/axoned/commit/c555990962ce4dcc8fdb0a3f54c59c678604df76))
* **mint:** migration get old mint denom ([8bae962](https://github.com/axone-protocol/axoned/commit/8bae96245500965027871f6968dfe02cd867e771))
* **store:** add ibcfee type storekey to kv store keys ([#419](https://github.com/axone-protocol/axoned/issues/419)) ([56a5f92](https://github.com/axone-protocol/axoned/commit/56a5f92f97db86e861778ab44158184a9cf8be13))
* **test:** allow default genesis creation in testutil ([bb35f38](https://github.com/axone-protocol/axoned/commit/bb35f38b473a0b9ea4025bf7a7fe9b85056d0a39))
* **wasm:** bump wasmvm to fixed version ([77c58cb](https://github.com/axone-protocol/axoned/commit/77c58cbae401958e0cb8ba38b1af2e68dfc9fd29))


### Code Refactoring

* **DIDComponents:** restrict type-checking for expected atom components ([ba5b2a3](https://github.com/axone-protocol/axoned/commit/ba5b2a357e7891ab91ecc64b2e9ead5bfdfbd72b))
* **logic:** add error field to anwser proto type ([b70ecc9](https://github.com/axone-protocol/axoned/commit/b70ecc97fbb68570145c06439753aea5f3482ab6))
* **logic:** standardize error messages for some predicates ([864b059](https://github.com/axone-protocol/axoned/commit/864b05998435c776d8111a9a2875887e977120ef))
* **logic:** standardize predicates errors ([ea7b38d](https://github.com/axone-protocol/axoned/commit/ea7b38defedf647c538c1378e8ae860a65535036))
* **logic:** standardize uri_encoded/3 errors ([b8d8133](https://github.com/axone-protocol/axoned/commit/b8d8133b4a2c20f11831b3e441f26b588a7e7776))

## [5.0.0](https://github.com/axone-protocol/axoned/compare/v4.1.0...v5.0.0) (2023-06-27)


### ⚠ BREAKING CHANGES

* **logic:** implement our own open/4 predicate
* **logic:** add base64Decode uri key on cosmwasm uri
* **logic:** specify whitelist / blacklist in its own type
* bump cosmos-sdk to 0.47.1
* **mint:** move x/param to x/mint state
* **logic:** allow predicates blacklisting configuration

### Features

* generate protobuf code ([c2bb70a](https://github.com/axone-protocol/axoned/commit/c2bb70ae57d702edf72e86d15a97e9aaac2671d2))
* implement migration from v2 to v3 ([19717ce](https://github.com/axone-protocol/axoned/commit/19717ce4642f9b5d25d07f5e23e9db2b334a8b70))
* **logic:** add base64Decode uri key on cosmwasm uri ([0516290](https://github.com/axone-protocol/axoned/commit/051629060d2d42dca736dfae91b46b6ecb672ab7))
* **logic:** add bounded buffer utility ([0ae4d43](https://github.com/axone-protocol/axoned/commit/0ae4d43bccd270a3b3dc47f6d29525b02419d8d1))
* **logic:** add convenient function to check nil pointers ([3d29c12](https://github.com/axone-protocol/axoned/commit/3d29c1286d5fbbee7b5724ed62472c5ca29693af))
* **logic:** add default_predicate_cost parameter ([52eac51](https://github.com/axone-protocol/axoned/commit/52eac5188de21c14d1c85b265df07209ed0f96a1))
* **logic:** add filtered virtual FS ([d35673d](https://github.com/axone-protocol/axoned/commit/d35673d0f2238874d9d12e85202b21140648ebb0))
* **logic:** add functional functions ([11f1738](https://github.com/axone-protocol/axoned/commit/11f17381edb5a4c42cfa0e48afa0d9ee821c78ef))
* **logic:** add gas policy parameters for the logic module ([2697fcb](https://github.com/axone-protocol/axoned/commit/2697fcb2230cb96bd51745d5103886f539684bc4))
* **logic:** add option for files whitelist and backlist on interpreter params ([fec5745](https://github.com/axone-protocol/axoned/commit/fec574518b7d7fdeca594b41858da5c11adfc32f))
* **logic:** add some functions to deal with urls ([5f9c8a1](https://github.com/axone-protocol/axoned/commit/5f9c8a16b38b0ac7546cfde9c4ef72051c329d10))
* **logic:** add source_files/1 predicate ([b89718c](https://github.com/axone-protocol/axoned/commit/b89718c5d59bda8fc73b3a3045de002523452c60))
* **logic:** add util func to extract json object attribute ([9b58497](https://github.com/axone-protocol/axoned/commit/9b58497948ba6d4d687fb80269bb5e03edef3550))
* **logic:** add v2 protobuf types for migration ([e00048a](https://github.com/axone-protocol/axoned/commit/e00048a4260316993fc7d77934cdf13477e0198a))
* **logic:** apply predicate costs ([9f83562](https://github.com/axone-protocol/axoned/commit/9f83562fc9f12e24554188e7cab963854049a350))
* **logic:** convert basic json object into prolog ([310459a](https://github.com/axone-protocol/axoned/commit/310459a2e351720c6497d78fb9462a5462cff9b7))
* **logic:** convert json string to terms ([e43010d](https://github.com/axone-protocol/axoned/commit/e43010d6869204e21b392ce54a7be362c359e522))
* **logic:** implement default predicate cost ([9681681](https://github.com/axone-protocol/axoned/commit/968168166df065e182d6201b06900b66b8779eea))
* **logic:** implement max length on read_string/3 predicate ([5c2a834](https://github.com/axone-protocol/axoned/commit/5c2a8345d3f1ec4f3d7eb1f2c7ec5397f01cc300))
* **logic:** implement our own open/4 predicate ([ec04a14](https://github.com/axone-protocol/axoned/commit/ec04a14b55a2c1a7c37ba7cc12e8db1cbfbf81ed))
* **logic:** implement preciate cost policy ([0ac899d](https://github.com/axone-protocol/axoned/commit/0ac899d45be8fb525ae1614b4477acc5ba882399))
* **logic:** implement support for user output ([0a65522](https://github.com/axone-protocol/axoned/commit/0a65522e223d2d149f6497c439baa07229a66aca))
* **logic:** implement virtual FS white/black list ([81c5c3e](https://github.com/axone-protocol/axoned/commit/81c5c3ee5b8c71d91bb70f0f2c6e063b91208f3d))
* **logic:** improve params definition ([6aa3d49](https://github.com/axone-protocol/axoned/commit/6aa3d4932675964265463c05aefa9b81d75609c2))
* **logic:** include json_prolog/2 into the registry ([5b85987](https://github.com/axone-protocol/axoned/commit/5b85987f102fb8efe57cf70c442d46d62fff1484))
* **logic:** introduce weithted_meter gas meter ([f668dd1](https://github.com/axone-protocol/axoned/commit/f668dd14a88ce3ef4ea71e1a750f91c4a7c78f53))
* **logic:** json_prolog/2 handle boolean ([7679f94](https://github.com/axone-protocol/axoned/commit/7679f94b132e251170b8c245d5d8b4bb87cd2438))
* **logic:** json_prolog/2 handle boolean and null ([9c3b7f8](https://github.com/axone-protocol/axoned/commit/9c3b7f8364e0ca8e7552b821873d76e3eefcfec0))
* **logic:** json_prolog/2 handle integer number ([a60f332](https://github.com/axone-protocol/axoned/commit/a60f3324c587826a5c470b6b985c527aed3474a8))
* **logic:** json_prolog/2 handle json array ([ff1f248](https://github.com/axone-protocol/axoned/commit/ff1f2481941c4681cd0e98e1d197b23823520d69))
* **logic:** json_prolog/2 handle json term to json object ([4bb3f9a](https://github.com/axone-protocol/axoned/commit/4bb3f9a0ad6825af1033ffed9888710efd353951))
* **logic:** json_prolog/2 handle list term to json array ([4e2b8b6](https://github.com/axone-protocol/axoned/commit/4e2b8b6bf590b03c1defc930d7d6082832c0b3d5))
* **logic:** json_prolog/2 handle null json value ([94e9c5b](https://github.com/axone-protocol/axoned/commit/94e9c5b85fbae0c16383d379c8e174b06be28c62))
* **logic:** json_prolog/2 handle string term to json string ([c0b5a6c](https://github.com/axone-protocol/axoned/commit/c0b5a6c5290f3aedbbd0b9f7d9f7a712f8d773f7))
* **logic:** regenerate protos ([cdfd71a](https://github.com/axone-protocol/axoned/commit/cdfd71abcc8b1cf63972c974a3abb489149010c6))
* **logic:** register msg to update params ([16435d2](https://github.com/axone-protocol/axoned/commit/16435d27cd7bf9988b48eb4a965ebdc396648447))
* **logic:** register our own open/4 predicate ([9af2390](https://github.com/axone-protocol/axoned/commit/9af23908c7a9d811b92bbd8c2870832f9b8d5b26))
* **logic:** register uri_encoded/3 on registry ([43b4cdf](https://github.com/axone-protocol/axoned/commit/43b4cdfe8aa3c21fc22f3ceb1b0858722c952ebb))
* **logic:** return rpc error when interpreter enconter an error ([2aed7e7](https://github.com/axone-protocol/axoned/commit/2aed7e70006569532501fba31bae534412089106))
* **logic:** specify user output support ([6343864](https://github.com/axone-protocol/axoned/commit/63438641aab39b347672fd67e8da64ee442b70a3))
* **logic:** update for sdk047 ([58c7efd](https://github.com/axone-protocol/axoned/commit/58c7efdd7d37840ad08b6876b3e1f5245a8455fe))
* **logic:** uri_encoded/3 detect component used ([276bf77](https://github.com/axone-protocol/axoned/commit/276bf77a815da6773b48c10a793ac304607c22ce))
* **logic:** uri_encoded/3 implement encoding component ([e02ca30](https://github.com/axone-protocol/axoned/commit/e02ca308f65d1b666ba37f13e3ad08d5ad236d14))
* **logic:** uri_encoded/3 implement unescape component ([9eeb484](https://github.com/axone-protocol/axoned/commit/9eeb48439cee13dde81c93eedfcd0b792ef81de4))
* **mint:** update mint module for sdk047 ([e8d4f90](https://github.com/axone-protocol/axoned/commit/e8d4f90555105e8cf87e8874f2a4626f1448d3ea))
* **mint:** update mint module for sdk047 ([fdae447](https://github.com/axone-protocol/axoned/commit/fdae4477dce464bb033bffadde37d66eac75fd49))
* **prolog:** implement read_string predicate ([a3601a0](https://github.com/axone-protocol/axoned/commit/a3601a0ca28dd3f5688022946fb90a5472ae2e3f))
* **vesting:** update to sdk047 ([d823c27](https://github.com/axone-protocol/axoned/commit/d823c27914aae040e95478ca127d4deed9fef453))


### Bug Fixes

* **ci:** fix linter ([5b60b4c](https://github.com/axone-protocol/axoned/commit/5b60b4c3d09d91cab17c340ff273e6290aa8e692))
* **lint:** add nolint for deprecated func in migration ([1b55d0d](https://github.com/axone-protocol/axoned/commit/1b55d0dc07a1f0197eabda1b036b24bf73a437ce))
* **lint:** gci import typo ([fe19886](https://github.com/axone-protocol/axoned/commit/fe198865a1d3559219686ca9dd527608b36c148e))
* **lint:** handle error ([8e28027](https://github.com/axone-protocol/axoned/commit/8e2802786bc6e0de566275fcf66de26c89aaec99))
* **lint:** make read_string gci-ed ([aee8b8b](https://github.com/axone-protocol/axoned/commit/aee8b8b2ac595c233e14b06a060b0afc450054cf))
* **lint:** reapply good gci import order ([765a2ed](https://github.com/axone-protocol/axoned/commit/765a2ed41651e93da0606603b0e8db4e56c3cde7))
* **logic:** avoid killing querying goroutine on gas limit exceeded ([86a184a](https://github.com/axone-protocol/axoned/commit/86a184ad27d9948bf6248752e4f9b7e118fb3688))
* **logic:** correct error messages ([f67c48b](https://github.com/axone-protocol/axoned/commit/f67c48bb3fb043ddb2ed073810326ef3ed48d82f))
* **logic:** do not convert empty string to Variable ([5c0d0fc](https://github.com/axone-protocol/axoned/commit/5c0d0fcc6efc18d1b5e9b8596464d43e56001158))
* **logic:** fix empty array management in json predicate ([054f854](https://github.com/axone-protocol/axoned/commit/054f85423bf4b83b39deaab7e096b615dffd491f))
* **logic:** fix error reported on url parse failure ([715baef](https://github.com/axone-protocol/axoned/commit/715baef3d5423ff1e6482233470578d436943186))
* proto linter ([79d96f0](https://github.com/axone-protocol/axoned/commit/79d96f09562f0b64304f6cb910708dfdfb038c43))
* remove unused proto import ([824c912](https://github.com/axone-protocol/axoned/commit/824c9126dcd1aca7bb85fd1cd3b4e8195dcd7663))
* **sdk:** solves barberry issue by updating the sdk ([25a3c6f](https://github.com/axone-protocol/axoned/commit/25a3c6f9ea8c45a3cc5d23c1d2c1b73800bf77d8))
* **test:** register get_char on testutil ([a5ef9ee](https://github.com/axone-protocol/axoned/commit/a5ef9ee468d58d087c0eb75633b8704f8e01fd9d))
* **test:** use good fs import on logic test ([0d6a7f4](https://github.com/axone-protocol/axoned/commit/0d6a7f49485c65b0f5e79076779994d8bb839c17))
* **upgrade:** stakingKeeper instance ([0dd7980](https://github.com/axone-protocol/axoned/commit/0dd7980f09db6868ed260d03bd78599021360080))


### Miscellaneous Chores

* bump cosmos-sdk to 0.47.1 ([492122d](https://github.com/axone-protocol/axoned/commit/492122dbaf874892579d477eea40cc22c5f9a7ad))


### Code Refactoring

* **logic:** allow predicates blacklisting configuration ([ec80998](https://github.com/axone-protocol/axoned/commit/ec80998d448b432603908fdef77281a1dbe0d70e))
* **logic:** specify whitelist / blacklist in its own type ([a8b2500](https://github.com/axone-protocol/axoned/commit/a8b2500203435e680e93c485bfc285d28140b0b8))
* **mint:** move x/param to x/mint state ([fe7c618](https://github.com/axone-protocol/axoned/commit/fe7c6184d108685a0060d79c8f6c22532d005363))

## [4.1.0](https://github.com/axone-protocol/axoned/compare/v4.0.0...v4.1.0) (2023-03-17)


### Features

* **logic:** add crypto_hash/2 predicate ([5c70aba](https://github.com/axone-protocol/axoned/commit/5c70aba06fbccf0ac73893ae41a883e62d23f248))
* **logic:** add hex_bytes/2 predicate ([eb167ee](https://github.com/axone-protocol/axoned/commit/eb167ee688e0200a4449ba1d76ce7f7a041242d5))
* **logic:** add wasm keeper interface ([4ccc32b](https://github.com/axone-protocol/axoned/commit/4ccc32bd7b71556330a2997b0b3dcd0c268d1f3e))
* **logic:** bech32_address/3 predicate conversion in one way ([ba1195a](https://github.com/axone-protocol/axoned/commit/ba1195a6dbb6bb9a23c9377913e488522a7208a2))
* **logic:** call wasm contract from fileSystem ([4eb6b47](https://github.com/axone-protocol/axoned/commit/4eb6b478cd9ab1b3d4268028611abfa86a3381e2))
* **logic:** convert base64 bech32 to bech32 encoded string ([7b24610](https://github.com/axone-protocol/axoned/commit/7b24610494bf338b2a47ec9d4c96a24caca812f3))
* **logic:** create custom file system handler ([b9cd4fb](https://github.com/axone-protocol/axoned/commit/b9cd4fb287baac37cbe0c1b90a8923523a03f7a3))
* **logic:** handle wasm uri on consult/1 predicate ([bbc7aae](https://github.com/axone-protocol/axoned/commit/bbc7aaedb828959ee1865ab24042e1ec2c15366d))
* **logic:** impl Read on Object file ([02cc0d1](https://github.com/axone-protocol/axoned/commit/02cc0d17e4bb7adf77f9a51cc442f9782ad0347e))
* **logic:** implements source_file/1 predicate ([8ceede1](https://github.com/axone-protocol/axoned/commit/8ceede17e4bb17b64eee8c26e17b174363781d8a))
* **logic:** return List of byte for crypto_hash/2 ([6534b9b](https://github.com/axone-protocol/axoned/commit/6534b9b5afc1a6f40680b159b149d4d4fc6f149c))


### Bug Fixes

* **linter:** remove unsused directive linter ([d61a5d6](https://github.com/axone-protocol/axoned/commit/d61a5d68ff19d0c4e38eedcb1c3bf7d1b0b55643))
* **logic:** add type safety on interface and rename it corrctly ([010bcc7](https://github.com/axone-protocol/axoned/commit/010bcc708c0724b404ddbd576fcf61dae75d4406))
* **logic:** check scheme on wasm fs ([a6f522b](https://github.com/axone-protocol/axoned/commit/a6f522b9afca923e4a2a4814736b956a231609c5))
* **logic:** comment typo ([04dcf6e](https://github.com/axone-protocol/axoned/commit/04dcf6ef416e7aa01056d6738abd827030c22169))
* **logic:** file time is the block height time ([a869ec6](https://github.com/axone-protocol/axoned/commit/a869ec691f8e4c84847a76116cf8bf31766c0f9c))
* **logic:** fix linter and empty path error ([db727a0](https://github.com/axone-protocol/axoned/commit/db727a06dce0859a22f8ebdd3856b9d3fafc0aad))
* **logic:** fix out of gas on goroutine ([b38ab90](https://github.com/axone-protocol/axoned/commit/b38ab902d643855e10723f135c27815ed8b71d26))
* **logic:** fix test after reviews ([bd57659](https://github.com/axone-protocol/axoned/commit/bd57659aac7c747ad6e3493da7ab6b0b1ddf72d5))
* **logic:** implement open on fs ([9bb1c10](https://github.com/axone-protocol/axoned/commit/9bb1c10649a5b19dd2750759542247af82fe7981))
* **logic:** linter and unit tests ([1ac14dd](https://github.com/axone-protocol/axoned/commit/1ac14dd826858dd9e777d82a1830c65a10da18cc))
* **logic:** make addressPairToBech32 private ([c31ba77](https://github.com/axone-protocol/axoned/commit/c31ba77dbd2b82c000078bdf8c89f8ce4fa432d3))
* **logic:** make source_file return multiple results instead of list ([d7e9526](https://github.com/axone-protocol/axoned/commit/d7e952666a10b86003a1ab033a7e8d84e0dc3449))
* **logic:** remove chek axone scheme ([dc53827](https://github.com/axone-protocol/axoned/commit/dc5382701292a12a13481752a16853f5720a0315))
* **logic:** remove unsued wasm on context ([306e364](https://github.com/axone-protocol/axoned/commit/306e3648ee0ed65fbeeadb95cb11f5996bad6b32))
* **logic:** remove unused files ([71f94b4](https://github.com/axone-protocol/axoned/commit/71f94b44a81ad661fbdaf675a9b042583f19b842))

## [4.0.0](https://github.com/axone-protocol/axoned/compare/v3.0.0...v4.0.0) (2023-02-15)


### ⚠ BREAKING CHANGES

* **proto:** align naming of a proto field in yaml marshal

### Features

* add utilitary functions ([830fe64](https://github.com/axone-protocol/axoned/commit/830fe64e586d7ccc6303b53dde8d71de84001b19))
* **buf:** generate new proto ([af9e24d](https://github.com/axone-protocol/axoned/commit/af9e24df210833ed5150be0fae56012d808ce1c8))
* **buf:** remove third party proto ([f65ba19](https://github.com/axone-protocol/axoned/commit/f65ba19dfc9fa12d908067a7e4b984383dd13b58))
* **buf:** use buf deps instead of third party ([bbcde9e](https://github.com/axone-protocol/axoned/commit/bbcde9ed8d112ba7a6aeb4eb7a826051d505cde4))
* compute total gas (sdk + interpreter) ([cd260df](https://github.com/axone-protocol/axoned/commit/cd260df292f2ed402d63eacd759af05f81a5ef37))
* implement grpc ask service ([cab9522](https://github.com/axone-protocol/axoned/commit/cab95228b39438c7b49c1849705531ab830ac515))
* implement logic business ([c4693bb](https://github.com/axone-protocol/axoned/commit/c4693bb54fcd30cb50d660f1122fd7a45eeef75b))
* improve command description and example ([2be2ee8](https://github.com/axone-protocol/axoned/commit/2be2ee8c5de18277a654cee94a2c96c8918e9271))
* **ledger:** fix Ledger build tag definition ([12cd92a](https://github.com/axone-protocol/axoned/commit/12cd92aa9340a615063fdf38b55e1f179193eae6))
* **logic:** add bank_balances predicate ([b0cc5cc](https://github.com/axone-protocol/axoned/commit/b0cc5cc9ba2a390061bd5d55623e3f5133a112b1))
* **logic:** add bank_spendable_coin predicate ([e7acefa](https://github.com/axone-protocol/axoned/commit/e7acefacfe999cd241854f140e04e2977654787a))
* **logic:** add block_height/1 predicate ([70b0bc0](https://github.com/axone-protocol/axoned/commit/70b0bc063af66239e3f225cbb2da935d681dc86f))
* **logic:** add block_time/1 predicate ([cc52351](https://github.com/axone-protocol/axoned/commit/cc523512347cb673518b4642d8d35955518c3489))
* **logic:** add chain_id/1 predicate ([eaac24b](https://github.com/axone-protocol/axoned/commit/eaac24bb4641faed945dc167c48bf9e644157c47))
* **logic:** add context extraction util ([64a5523](https://github.com/axone-protocol/axoned/commit/64a5523ba8e22da4f130efd300e51aed3ae9baa6))
* **logic:** add did_components/2 predicate ([09976d9](https://github.com/axone-protocol/axoned/commit/09976d9c41ea19207528c24ca3cd2a7a3b1e3030))
* **logic:** add go-routine safe version of GasMeter ([5c1b4b9](https://github.com/axone-protocol/axoned/commit/5c1b4b9356e3c7ab92ab87363b3dd7f4a8045e15))
* **logic:** add limit context ([3569103](https://github.com/axone-protocol/axoned/commit/3569103e364093862b6fa8a922fd48f0984bba3b))
* **logic:** add locked coins method on expected bank keeper ([48b10e5](https://github.com/axone-protocol/axoned/commit/48b10e5499905b67a9eeeea76f966c5172f187ae))
* **logic:** add locked coins predicate implementation ([7a926c5](https://github.com/axone-protocol/axoned/commit/7a926c5937399be5b9cb11119c570169ee759b54))
* **logic:** allow return all spendable coins balances ([e0a7de5](https://github.com/axone-protocol/axoned/commit/e0a7de5bfd0917a4480c16e0674f2a28c59298c4))
* **logic:** call ask query from cli ([d8f343d](https://github.com/axone-protocol/axoned/commit/d8f343d95e92bc7220f5c3e6877e022816411d7d))
* **logic:** change type of limit params as *Uint allowing nil value ([a8f5a60](https://github.com/axone-protocol/axoned/commit/a8f5a60dd2a6b830cddda2e56ff56def8625c466))
* **logic:** decouple wasm ask response from grpc type ([03128f5](https://github.com/axone-protocol/axoned/commit/03128f5f3368af3a0104cd38356c4622fe94bbe3))
* **logic:** improve error messages ([5f2028e](https://github.com/axone-protocol/axoned/commit/5f2028e617c6821e7966699fcbd90de1f9967589))
* **logic:** improve parameters configuration ([d1396bb](https://github.com/axone-protocol/axoned/commit/d1396bb43339b25823f4179df994ae34a3d9685a))
* **logic:** inject auth and bank keeper ([578fd39](https://github.com/axone-protocol/axoned/commit/578fd39fd21d3e17665ecc6b3e9779ebaf4c9b16))
* **logic:** move logic query into a dedicated file ([9a4a047](https://github.com/axone-protocol/axoned/commit/9a4a0472dc1b4b367f90f4f65da218f5eca5e2e7))
* **logic:** register params for genesis ([acb64a9](https://github.com/axone-protocol/axoned/commit/acb64a95390c7276f07c2c39938eea4d1ee66a6b))
* **logic:** register the locked coins predicated ([b61ce52](https://github.com/axone-protocol/axoned/commit/b61ce526e0a5e1a38a8504f5cd96e3e79819b284))
* **logic:** simplify wasm custom query integration ([383a0e7](https://github.com/axone-protocol/axoned/commit/383a0e7162a640c7799aa3ace2db4a3c5862a0c6))
* **logic:** specify logic query operation ([8b385e0](https://github.com/axone-protocol/axoned/commit/8b385e0e1dc9eddca59daf1d42ef06fb3033ae27))
* **logic:** specify parameters for module logic ([6297da0](https://github.com/axone-protocol/axoned/commit/6297da050c989ab5c11788cc2b3c8d987401c51f))
* **upgrade:** allow add custom proposal file on chain-upgrade sript ([ef68aa4](https://github.com/axone-protocol/axoned/commit/ef68aa4592a2aa431fc2cff4575ff3ba21cfd8b5))
* **upgrade:** create package for register upgrades ([ef21308](https://github.com/axone-protocol/axoned/commit/ef213086d5cf0f9df82a235cc64a42291aa3f065))
* **wasm:** implements CustomQuerier with logic module ([3c68496](https://github.com/axone-protocol/axoned/commit/3c684960e4f2cfae8d278bb64db1e5ff2a2d4dbe))
* **wasm:** wire the wasm CustomQuerier in the app ([9435cf6](https://github.com/axone-protocol/axoned/commit/9435cf6e0bdf9b6df6b8a418350623344172e41e))


### Bug Fixes

* **ci:** add build before install for test blockchain ([81c98c5](https://github.com/axone-protocol/axoned/commit/81c98c50afbca8ee83577cf4d1a9c0e18babc891))
* **ci:** fix changed file conditions for run test workflows ([3496204](https://github.com/axone-protocol/axoned/commit/349620421a1f63265472acd714ecc2676e7a57f9))
* **docs:** change boolean as string for trigger updtae doc ([8850aee](https://github.com/axone-protocol/axoned/commit/8850aee5c4e7ec78e8a033d4aaae1391f45b0753))
* **docs:** fix linter generation ([0ffa36e](https://github.com/axone-protocol/axoned/commit/0ffa36eb6e70aaa659c08e1333975e0bf15337c6))
* **docs:** set the workflow id instead of name to fix not found ([9040df8](https://github.com/axone-protocol/axoned/commit/9040df81f0feac0ad2b5d51fac312e7a45b187c7))
* don't load program from filesystem ([eaef7b3](https://github.com/axone-protocol/axoned/commit/eaef7b3d80462393a411d29b65154519a301bd58))
* fix error message (predicate name was incorrect) ([d4d6a2d](https://github.com/axone-protocol/axoned/commit/d4d6a2da1cf707eb31d8e4afb6f1ae1dc6d2a052))
* fix typo in predicate name ([a80dd89](https://github.com/axone-protocol/axoned/commit/a80dd8913b091c40fffcbb4fd8e00757b6c744cd))
* fix wrong types (was problematic for type assertions) ([2b7e7bc](https://github.com/axone-protocol/axoned/commit/2b7e7bcefee08315f82b82a365fbca8245ba6254))
* **lint:** add updated generated doc ([0280a12](https://github.com/axone-protocol/axoned/commit/0280a12ce9c7fe0700401ec105de7ce6be6948c8))
* **lint:** fix golangci lint ([cbbe5db](https://github.com/axone-protocol/axoned/commit/cbbe5dbecb61aad0bc46f7b3685996d6d4a9c20d))
* **lint:** remove lint error for upgrade ([d5bb919](https://github.com/axone-protocol/axoned/commit/d5bb91990f0d0acef5db70e743377f6dc39f66b5))
* **logic:** ensure keepers in interpreter exec ctx ([901d7f2](https://github.com/axone-protocol/axoned/commit/901d7f213c2b745744b46dcf72738bdbfb6c24dc))
* **logic:** fix the description inversion of the flags ([5117430](https://github.com/axone-protocol/axoned/commit/5117430785a6930146259f928e118b074dbece39))
* **logic:** insert bankKeeper and accountKeeper into context ([e5338c1](https://github.com/axone-protocol/axoned/commit/e5338c19a8871b55f622c27fad8298e9276e6061))
* **logic:** register bank_balances predicate ([3ad6317](https://github.com/axone-protocol/axoned/commit/3ad6317717623ef5c61f69722e22c7f09dd64413))
* **logic:** remove gocognit linter for tests ([28904d4](https://github.com/axone-protocol/axoned/commit/28904d49008d897581b90a797e5ee5ec8b32bc5f))
* **logic:** sort result for locked coin denom ([bc5e867](https://github.com/axone-protocol/axoned/commit/bc5e8676f493ee4dc853d1aaedafa4533dfd1451))
* **logic:** typo in doc comment ([d098256](https://github.com/axone-protocol/axoned/commit/d098256ed7c81d29b17bf413aa2448d6a42db004))
* **proto:** align naming of a proto field in yaml marshal ([4fe9a67](https://github.com/axone-protocol/axoned/commit/4fe9a672f9d881b39394daf3b1a44b471e10f2c9))

## [3.0.0](https://github.com/axone-protocol/axoned/compare/v2.2.0...v3.0.0) (2022-11-30)


### ⚠ BREAKING CHANGES

* **mint:** configure annual provision and target supply on first block

### Features

* **docs:** trigger docs version update workflow on docs repo ([a224a10](https://github.com/axone-protocol/axoned/commit/a224a106a75188b00a62d8c936cca67da2629890))
* **docs:** trigger the docs workflow to update documentation ([e0558aa](https://github.com/axone-protocol/axoned/commit/e0558aa62ecd58ef46dccc0677b8ec91e026a1a8))
* **mint:** add target_supply on proto ([3c198f7](https://github.com/axone-protocol/axoned/commit/3c198f7fd1b7794c6eeb197934f677bb6bde339b))
* **mint:** configure annual provision and target supply on first block ([31d5884](https://github.com/axone-protocol/axoned/commit/31d5884fc0a7ebbbdf23ce6921726f0489775359))
* **mint:** implement inflation calculation ([42bfa4c](https://github.com/axone-protocol/axoned/commit/42bfa4c8fb5de4bb269970bf852a18708d96aeff))
* **mint:** move annual reduction factor from minter to minter params ([f731e66](https://github.com/axone-protocol/axoned/commit/f731e6699fcae0c9c48977d5641a25182c976c2f))
* **mint:** remove axone old inflation calc func ([9956d1b](https://github.com/axone-protocol/axoned/commit/9956d1b5c8fe8e0f06ed4aacf8abc799ba322d5a))
* **mint:** set mint param on proto ([ade514e](https://github.com/axone-protocol/axoned/commit/ade514ecf383430dd0dfeaf283b4882659d3afcc))
* **mint:** use local proto ([cbf22f6](https://github.com/axone-protocol/axoned/commit/cbf22f60eff9e2bfed4ca1945a9ee24e56c0bb0d))
* **mint:** use own mint module ([af1386e](https://github.com/axone-protocol/axoned/commit/af1386e01a9750d0807c0291507de6751345344d))


### Bug Fixes

* **docs:** change comments syntax in markdown template ([9e8496b](https://github.com/axone-protocol/axoned/commit/9e8496b6279bd5361587027cd0f36d0230980e1e))
* **docs:** fix linting issue ([13709c4](https://github.com/axone-protocol/axoned/commit/13709c4011473e8e11f1040ca4b688da4b45a463))
* **docs:** ignore linting of generated protobuf docs ([84aaab2](https://github.com/axone-protocol/axoned/commit/84aaab20bc0870f71efbb3095cce10fd44078e20))
* **mint:** avoid return negative coin ([b25b1f3](https://github.com/axone-protocol/axoned/commit/b25b1f3b1054c4bb53a940c2eb672a594c47e384))
* **mint:** make linter more happy ([35bed9f](https://github.com/axone-protocol/axoned/commit/35bed9f864c19229736bf4b71e37cd3f04a2f7fe))
* **mint:** spelling mistake ([d043d1b](https://github.com/axone-protocol/axoned/commit/d043d1b49b83028602c10c586500065bb34b34c3))

## [2.2.0](https://github.com/axone-protocol/axoned/compare/v2.1.1...v2.2.0) (2022-10-13)


### Features

* **ledger:** add build dependancies and follow sdk standards ([dcd4135](https://github.com/axone-protocol/axoned/commit/dcd41350615132d983edf40ea0e8a490f0f5589f))
* **ledger:** bump ledger to v0.9.3 ([bff91ba](https://github.com/axone-protocol/axoned/commit/bff91baf73aa839f2beb33d43e9568ad7122daf4))
* **ledger:** fix install dependancies in dockerfile ([5f90752](https://github.com/axone-protocol/axoned/commit/5f9075231ff4e195dacd49b56a6402ee3acccd88))
* **ledger:** update ledger-go to support Ledger Nano S Plus ([4dc7f4d](https://github.com/axone-protocol/axoned/commit/4dc7f4d08d09f1b864010ab89df286bf80cf3ee8))

## [2.1.1](https://github.com/axone-protocol/axoned/compare/v2.1.0...v2.1.1) (2022-10-10)


### Bug Fixes

* **mint:** provide annual inflation rate ([608af3f](https://github.com/axone-protocol/axoned/commit/608af3f20ccc0afedfd518b56ae84c679df21502))
* **mint:** set initial inflation to 7.5% instead of 15% ([7bbd048](https://github.com/axone-protocol/axoned/commit/7bbd04802ce402967a56550f9c843c608fd91877))

## [2.1.0](https://github.com/axone-protocol/axoned/compare/v2.0.0...v2.1.0) (2022-10-05)


### Features

* **cliff:** add cliff cmd on vesting transaction ([ccff37c](https://github.com/axone-protocol/axoned/commit/ccff37ce34941297fc924b277381d08c7a48a367))
* **cliff:** add vesting-cliff-time flags on add-genesis account cmd ([ea3e2c5](https://github.com/axone-protocol/axoned/commit/ea3e2c5abb28ae24fa7f6633c3a3f5f95196f92f))
* **cliff:** override add-genesis-account ([434d418](https://github.com/axone-protocol/axoned/commit/434d4181a114675cd3c894378d0179266ec1fd3b))
* **cliff:** register cliff vesting account msg ([9106919](https://github.com/axone-protocol/axoned/commit/910691901185a32a2699c67fe5560f55e37fddba))
* implment axone inflaction calculation function ([2e95801](https://github.com/axone-protocol/axoned/commit/2e958010bdbaa78af0f0a1958c9c0493c2fbfed7))
* use axone inflation calculation fn (instead of default one) ([bdca893](https://github.com/axone-protocol/axoned/commit/bdca89319b5546fca1aabe2ff5bece1c05258558))
* use axone vesting module ([7493de9](https://github.com/axone-protocol/axoned/commit/7493de90a0fb601b9ab7b10eb0526a59652d5eef))
* use third party to generate proto ([cb4f5bb](https://github.com/axone-protocol/axoned/commit/cb4f5bb36fee1cc54d59bf902eae229fe2421bd6))


### Bug Fixes

* **cliff:** improve verification on cliff msg tx ([4dfaa5b](https://github.com/axone-protocol/axoned/commit/4dfaa5b2624fcd39b12459d83275027f8d09507d))
* **ibc:** ensure ibc fees are managed ([c26b0db](https://github.com/axone-protocol/axoned/commit/c26b0db71dbd5b71788a7f333b673bc75ac80198))
* make linter happy ([584851c](https://github.com/axone-protocol/axoned/commit/584851c5b248112fc10ea3ff0a3f3873a3878a06))

## [2.0.0](https://github.com/axone-protocol/axoned/compare/v1.3.0...v2.0.0) (2022-09-23)


### ⚠ BREAKING CHANGES

* reboot chain with ignite cli v0.24.0

### Features

* add logic module params to genesis files ([9ac7ef8](https://github.com/axone-protocol/axoned/commit/9ac7ef8a7ecbe72fd8efd983ed138fd375aaf1fa))
* scaffold logic module using ignite ([81ee269](https://github.com/axone-protocol/axoned/commit/81ee26997d3e941fa24655ea89ce3dad93b6cfd2))
* update openapi documentation (synced with code) ([787ff01](https://github.com/axone-protocol/axoned/commit/787ff01148f1721f75594824c6da715ade5f18ea))
* **wasm:** prepare ante handler with wasm decorators ([afb4748](https://github.com/axone-protocol/axoned/commit/afb4748f5ea7b3b6723226235d9d93e171605c9f))
* **wasm:** wire wasm module in app ([b163790](https://github.com/axone-protocol/axoned/commit/b163790cb60c2221793742eed6fecc22d4254315))


### Bug Fixes

* fix (pre-)genesis files after 0.46 cosmos sdk migration ([da284a0](https://github.com/axone-protocol/axoned/commit/da284a0d4993f2231794ddbfae16b8151b23fbfc))
* use proper versions of buf protoc plugins ([1ca5e1d](https://github.com/axone-protocol/axoned/commit/1ca5e1d6127d255cf8f11375a7256e5c07efedf0))
* **workflow:** use secret for dockerhub user ([0d95c94](https://github.com/axone-protocol/axoned/commit/0d95c945e7c0e858ae22a2e7d573853e889d4473))


### Code Refactoring

* reboot chain with ignite cli v0.24.0 ([423179e](https://github.com/axone-protocol/axoned/commit/423179e5028de57b7859bbd6ebcb5d12a4b42fb5))

# [1.3.0](https://github.com/axone-protocol/axoned/compare/v1.2.0...v1.3.0) (2022-07-08)


### Bug Fixes

* fix genesis file after cosmos SDK update ([ff03ba9](https://github.com/axone-protocol/axoned/commit/ff03ba906313f8cbdeb77b3fb158cdf188f197fb))
* generate pre-genesis file with uknow unit ([677fc0c](https://github.com/axone-protocol/axoned/commit/677fc0cd16f18f9e650d2df4a669932da9e8200c))
* make it start 🚀 ([0cdd4db](https://github.com/axone-protocol/axoned/commit/0cdd4db2bdd16fcbe2788d5b6802d7ed9402a2bb))
* make linters happy ([309193a](https://github.com/axone-protocol/axoned/commit/309193aae612a94bc8d1ee6f557af204a7720044))
* references all modules in SetOrder* functions ([84f9fe2](https://github.com/axone-protocol/axoned/commit/84f9fe201d0214b32a0ff3d2f4e4893bf879ba49))


### Features

* add a proper description for the axone CLI ([bc74f2c](https://github.com/axone-protocol/axoned/commit/bc74f2c4f158895cf45534b33e0c32a2fb92cd90))
* **denom:** add uknow & know denoms metadata ([55d52ef](https://github.com/axone-protocol/axoned/commit/55d52efceff6cbd9714dc7a6014029824716f138))
* handle wasm proposals in app ([44b10c0](https://github.com/axone-protocol/axoned/commit/44b10c0b91cf10ff8b17b86737067ff1ee47586a))
* implement genesis account cmd ([0368a11](https://github.com/axone-protocol/axoned/commit/0368a111f374de99c6e3cb0a2ecc1d8c641e7187))
* implement genesis wasm cmd ([b4ac0bc](https://github.com/axone-protocol/axoned/commit/b4ac0bc82b3cdf3c6b7c0f24cd5739eb7f245ae8))
* implement axone encoding config ([acf6a26](https://github.com/axone-protocol/axoned/commit/acf6a26fc5d0d42df5e77b3bca65f4e663a87200))
* implement root cmd ([d065735](https://github.com/axone-protocol/axoned/commit/d0657358ee65dbd1c516e81eedb56a5ed6c63fd9))
* prepare ante handler with wasm decorators ([fb32135](https://github.com/axone-protocol/axoned/commit/fb32135288dcff37b6a0a25897210ab78461b745))
* provide default app encoding config ([602d2db](https://github.com/axone-protocol/axoned/commit/602d2db7dad855871d550d849b964e35ca2f26a1))
* re-sync openapi specification (after monitoringp removal) ([96dba69](https://github.com/axone-protocol/axoned/commit/96dba6985cbce0def82cb56671a265a915bb3116))
* remove ignite dependencies ([8cf0c4b](https://github.com/axone-protocol/axoned/commit/8cf0c4b2a9314f42b244193e63de5b3d90ab8a42))
* remove monitoringp module ([f97b8b9](https://github.com/axone-protocol/axoned/commit/f97b8b94895a7aa6f66bfce6880d263b67cf8cec))
* remove unwanted monitoringp module ([6be3175](https://github.com/axone-protocol/axoned/commit/6be317582c22ea5172932126e33f9efefc725d66))
* replace ignite root cmd by ours ([e92cfe4](https://github.com/axone-protocol/axoned/commit/e92cfe43d94243ddc1faabdd27f13bf401827297))
* wire wasm module in app ([ac56cbe](https://github.com/axone-protocol/axoned/commit/ac56cbe76bb5f16c75c63cfb1675450217ab23a2))

# [1.2.0](https://github.com/axone-protocol/axoned/compare/v1.1.0...v1.2.0) (2022-04-21)


### Features

* check uri format on trigger service msg ([b9ae987](https://github.com/axone-protocol/axoned/commit/b9ae987b081ee7b2dca5da09cabe088d19c2a728))
* extend debug cmd adding decode-blocks ([2959444](https://github.com/axone-protocol/axoned/commit/29594444e8c749f975e28316dddb0b1322bfebfa))
* implements debug blocks proto base64 decode ([1749e70](https://github.com/axone-protocol/axoned/commit/1749e70151ae16b0abfe326b70501ab64ca42d2a))
* scaffold trigger-service message ([57a54eb](https://github.com/axone-protocol/axoned/commit/57a54eb49cff5ee66957a1d2d795af5edab45bde))

# [1.1.0](https://github.com/axone-protocol/axoned/compare/v1.0.0...v1.1.0) (2022-03-16)


### Bug Fixes

* fix missing `id` in message ([c8f81fe](https://github.com/axone-protocol/axoned/commit/c8f81febc74622397ffb2b2a4e103ef28511bb92))


### Features

* implement dataspace creation ([8ec1074](https://github.com/axone-protocol/axoned/commit/8ec10742fd6271720f98fce93e20626447008841))
* scaffold module knowledge ([82c95e7](https://github.com/axone-protocol/axoned/commit/82c95e7acebf30df7ca924383ad8e5de3da533bc))
* update openapi documentation ([03a541c](https://github.com/axone-protocol/axoned/commit/03a541c732f8fd94701015317c5d1fcb0574c688))

# 1.0.0 (2022-02-08)


### Bug Fixes

* Add make init and fix start ([30468fe](https://github.com/axone-protocol/axoned/commit/30468fed2d6b8b86e60f0cb44479547b05bded18))


### Code Refactoring

* rename repository from "axone" to "axoned" ([1a3b24f](https://github.com/axone-protocol/axoned/commit/1a3b24f383e7ef18f0bd2a2e93b774813e824339))


### Features

* add scaffolded code ([232da23](https://github.com/axone-protocol/axoned/commit/232da23ee3428c3aec0d68a41771832aa6502cb3))


### BREAKING CHANGES

* this change has a few implications, such as changing
the name of the published docker image.
