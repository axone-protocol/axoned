# ØKP4 protocol changelog

## [4.1.0](https://github.com/okp4/okp4d/compare/v4.0.0...v4.1.0) (2023-03-17)


### Features

* **logic:** add crypto_hash/2 predicate ([5c70aba](https://github.com/okp4/okp4d/commit/5c70aba06fbccf0ac73893ae41a883e62d23f248))
* **logic:** add hex_bytes/2 predicate ([eb167ee](https://github.com/okp4/okp4d/commit/eb167ee688e0200a4449ba1d76ce7f7a041242d5))
* **logic:** add wasm keeper interface ([4ccc32b](https://github.com/okp4/okp4d/commit/4ccc32bd7b71556330a2997b0b3dcd0c268d1f3e))
* **logic:** bech32_address/3 predicate conversion in one way ([ba1195a](https://github.com/okp4/okp4d/commit/ba1195a6dbb6bb9a23c9377913e488522a7208a2))
* **logic:** call wasm contract from fileSystem ([4eb6b47](https://github.com/okp4/okp4d/commit/4eb6b478cd9ab1b3d4268028611abfa86a3381e2))
* **logic:** convert base64 bech32 to bech32 encoded string ([7b24610](https://github.com/okp4/okp4d/commit/7b24610494bf338b2a47ec9d4c96a24caca812f3))
* **logic:** create custom file system handler ([b9cd4fb](https://github.com/okp4/okp4d/commit/b9cd4fb287baac37cbe0c1b90a8923523a03f7a3))
* **logic:** handle wasm uri on consult/1 predicate ([bbc7aae](https://github.com/okp4/okp4d/commit/bbc7aaedb828959ee1865ab24042e1ec2c15366d))
* **logic:** impl Read on Object file ([02cc0d1](https://github.com/okp4/okp4d/commit/02cc0d17e4bb7adf77f9a51cc442f9782ad0347e))
* **logic:** implements source_file/1 predicate ([8ceede1](https://github.com/okp4/okp4d/commit/8ceede17e4bb17b64eee8c26e17b174363781d8a))
* **logic:** return List of byte for crypto_hash/2 ([6534b9b](https://github.com/okp4/okp4d/commit/6534b9b5afc1a6f40680b159b149d4d4fc6f149c))


### Bug Fixes

* **linter:** remove unsused directive linter ([d61a5d6](https://github.com/okp4/okp4d/commit/d61a5d68ff19d0c4e38eedcb1c3bf7d1b0b55643))
* **logic:** add type safety on interface and rename it corrctly ([010bcc7](https://github.com/okp4/okp4d/commit/010bcc708c0724b404ddbd576fcf61dae75d4406))
* **logic:** check scheme on wasm fs ([a6f522b](https://github.com/okp4/okp4d/commit/a6f522b9afca923e4a2a4814736b956a231609c5))
* **logic:** comment typo ([04dcf6e](https://github.com/okp4/okp4d/commit/04dcf6ef416e7aa01056d6738abd827030c22169))
* **logic:** file time is the block height time ([a869ec6](https://github.com/okp4/okp4d/commit/a869ec691f8e4c84847a76116cf8bf31766c0f9c))
* **logic:** fix linter and empty path error ([db727a0](https://github.com/okp4/okp4d/commit/db727a06dce0859a22f8ebdd3856b9d3fafc0aad))
* **logic:** fix out of gas on goroutine ([b38ab90](https://github.com/okp4/okp4d/commit/b38ab902d643855e10723f135c27815ed8b71d26))
* **logic:** fix test after reviews ([bd57659](https://github.com/okp4/okp4d/commit/bd57659aac7c747ad6e3493da7ab6b0b1ddf72d5))
* **logic:** implement open on fs ([9bb1c10](https://github.com/okp4/okp4d/commit/9bb1c10649a5b19dd2750759542247af82fe7981))
* **logic:** linter and unit tests ([1ac14dd](https://github.com/okp4/okp4d/commit/1ac14dd826858dd9e777d82a1830c65a10da18cc))
* **logic:** make addressPairToBech32 private ([c31ba77](https://github.com/okp4/okp4d/commit/c31ba77dbd2b82c000078bdf8c89f8ce4fa432d3))
* **logic:** make source_file return multiple results instead of list ([d7e9526](https://github.com/okp4/okp4d/commit/d7e952666a10b86003a1ab033a7e8d84e0dc3449))
* **logic:** remove chek okp4 scheme ([dc53827](https://github.com/okp4/okp4d/commit/dc5382701292a12a13481752a16853f5720a0315))
* **logic:** remove unsued wasm on context ([306e364](https://github.com/okp4/okp4d/commit/306e3648ee0ed65fbeeadb95cb11f5996bad6b32))
* **logic:** remove unused files ([71f94b4](https://github.com/okp4/okp4d/commit/71f94b44a81ad661fbdaf675a9b042583f19b842))

## [4.0.0](https://github.com/okp4/okp4d/compare/v3.0.0...v4.0.0) (2023-02-15)


### ⚠ BREAKING CHANGES

* **proto:** align naming of a proto field in yaml marshal

### Features

* add utilitary functions ([830fe64](https://github.com/okp4/okp4d/commit/830fe64e586d7ccc6303b53dde8d71de84001b19))
* **buf:** generate new proto ([af9e24d](https://github.com/okp4/okp4d/commit/af9e24df210833ed5150be0fae56012d808ce1c8))
* **buf:** remove third party proto ([f65ba19](https://github.com/okp4/okp4d/commit/f65ba19dfc9fa12d908067a7e4b984383dd13b58))
* **buf:** use buf deps instead of third party ([bbcde9e](https://github.com/okp4/okp4d/commit/bbcde9ed8d112ba7a6aeb4eb7a826051d505cde4))
* compute total gas (sdk + interpreter) ([cd260df](https://github.com/okp4/okp4d/commit/cd260df292f2ed402d63eacd759af05f81a5ef37))
* implement grpc ask service ([cab9522](https://github.com/okp4/okp4d/commit/cab95228b39438c7b49c1849705531ab830ac515))
* implement logic business ([c4693bb](https://github.com/okp4/okp4d/commit/c4693bb54fcd30cb50d660f1122fd7a45eeef75b))
* improve command description and example ([2be2ee8](https://github.com/okp4/okp4d/commit/2be2ee8c5de18277a654cee94a2c96c8918e9271))
* **ledger:** fix Ledger build tag definition ([12cd92a](https://github.com/okp4/okp4d/commit/12cd92aa9340a615063fdf38b55e1f179193eae6))
* **logic:** add bank_balances predicate ([b0cc5cc](https://github.com/okp4/okp4d/commit/b0cc5cc9ba2a390061bd5d55623e3f5133a112b1))
* **logic:** add bank_spendable_coin predicate ([e7acefa](https://github.com/okp4/okp4d/commit/e7acefacfe999cd241854f140e04e2977654787a))
* **logic:** add block_height/1 predicate ([70b0bc0](https://github.com/okp4/okp4d/commit/70b0bc063af66239e3f225cbb2da935d681dc86f))
* **logic:** add block_time/1 predicate ([cc52351](https://github.com/okp4/okp4d/commit/cc523512347cb673518b4642d8d35955518c3489))
* **logic:** add chain_id/1 predicate ([eaac24b](https://github.com/okp4/okp4d/commit/eaac24bb4641faed945dc167c48bf9e644157c47))
* **logic:** add context extraction util ([64a5523](https://github.com/okp4/okp4d/commit/64a5523ba8e22da4f130efd300e51aed3ae9baa6))
* **logic:** add did_components/2 predicate ([09976d9](https://github.com/okp4/okp4d/commit/09976d9c41ea19207528c24ca3cd2a7a3b1e3030))
* **logic:** add go-routine safe version of GasMeter ([5c1b4b9](https://github.com/okp4/okp4d/commit/5c1b4b9356e3c7ab92ab87363b3dd7f4a8045e15))
* **logic:** add limit context ([3569103](https://github.com/okp4/okp4d/commit/3569103e364093862b6fa8a922fd48f0984bba3b))
* **logic:** add locked coins method on expected bank keeper ([48b10e5](https://github.com/okp4/okp4d/commit/48b10e5499905b67a9eeeea76f966c5172f187ae))
* **logic:** add locked coins predicate implementation ([7a926c5](https://github.com/okp4/okp4d/commit/7a926c5937399be5b9cb11119c570169ee759b54))
* **logic:** allow return all spendable coins balances ([e0a7de5](https://github.com/okp4/okp4d/commit/e0a7de5bfd0917a4480c16e0674f2a28c59298c4))
* **logic:** call ask query from cli ([d8f343d](https://github.com/okp4/okp4d/commit/d8f343d95e92bc7220f5c3e6877e022816411d7d))
* **logic:** change type of limit params as *Uint allowing nil value ([a8f5a60](https://github.com/okp4/okp4d/commit/a8f5a60dd2a6b830cddda2e56ff56def8625c466))
* **logic:** decouple wasm ask response from grpc type ([03128f5](https://github.com/okp4/okp4d/commit/03128f5f3368af3a0104cd38356c4622fe94bbe3))
* **logic:** improve error messages ([5f2028e](https://github.com/okp4/okp4d/commit/5f2028e617c6821e7966699fcbd90de1f9967589))
* **logic:** improve parameters configuration ([d1396bb](https://github.com/okp4/okp4d/commit/d1396bb43339b25823f4179df994ae34a3d9685a))
* **logic:** inject auth and bank keeper ([578fd39](https://github.com/okp4/okp4d/commit/578fd39fd21d3e17665ecc6b3e9779ebaf4c9b16))
* **logic:** move logic query into a dedicated file ([9a4a047](https://github.com/okp4/okp4d/commit/9a4a0472dc1b4b367f90f4f65da218f5eca5e2e7))
* **logic:** register params for genesis ([acb64a9](https://github.com/okp4/okp4d/commit/acb64a95390c7276f07c2c39938eea4d1ee66a6b))
* **logic:** register the locked coins predicated ([b61ce52](https://github.com/okp4/okp4d/commit/b61ce526e0a5e1a38a8504f5cd96e3e79819b284))
* **logic:** simplify wasm custom query integration ([383a0e7](https://github.com/okp4/okp4d/commit/383a0e7162a640c7799aa3ace2db4a3c5862a0c6))
* **logic:** specify logic query operation ([8b385e0](https://github.com/okp4/okp4d/commit/8b385e0e1dc9eddca59daf1d42ef06fb3033ae27))
* **logic:** specify parameters for module logic ([6297da0](https://github.com/okp4/okp4d/commit/6297da050c989ab5c11788cc2b3c8d987401c51f))
* **upgrade:** allow add custom proposal file on chain-upgrade sript ([ef68aa4](https://github.com/okp4/okp4d/commit/ef68aa4592a2aa431fc2cff4575ff3ba21cfd8b5))
* **upgrade:** create package for register upgrades ([ef21308](https://github.com/okp4/okp4d/commit/ef213086d5cf0f9df82a235cc64a42291aa3f065))
* **wasm:** implements CustomQuerier with logic module ([3c68496](https://github.com/okp4/okp4d/commit/3c684960e4f2cfae8d278bb64db1e5ff2a2d4dbe))
* **wasm:** wire the wasm CustomQuerier in the app ([9435cf6](https://github.com/okp4/okp4d/commit/9435cf6e0bdf9b6df6b8a418350623344172e41e))


### Bug Fixes

* **ci:** add build before install for test blockchain ([81c98c5](https://github.com/okp4/okp4d/commit/81c98c50afbca8ee83577cf4d1a9c0e18babc891))
* **ci:** fix changed file conditions for run test workflows ([3496204](https://github.com/okp4/okp4d/commit/349620421a1f63265472acd714ecc2676e7a57f9))
* **docs:** change boolean as string for trigger updtae doc ([8850aee](https://github.com/okp4/okp4d/commit/8850aee5c4e7ec78e8a033d4aaae1391f45b0753))
* **docs:** fix linter generation ([0ffa36e](https://github.com/okp4/okp4d/commit/0ffa36eb6e70aaa659c08e1333975e0bf15337c6))
* **docs:** set the workflow id instead of name to fix not found ([9040df8](https://github.com/okp4/okp4d/commit/9040df81f0feac0ad2b5d51fac312e7a45b187c7))
* don't load program from filesystem ([eaef7b3](https://github.com/okp4/okp4d/commit/eaef7b3d80462393a411d29b65154519a301bd58))
* fix error message (predicate name was incorrect) ([d4d6a2d](https://github.com/okp4/okp4d/commit/d4d6a2da1cf707eb31d8e4afb6f1ae1dc6d2a052))
* fix typo in predicate name ([a80dd89](https://github.com/okp4/okp4d/commit/a80dd8913b091c40fffcbb4fd8e00757b6c744cd))
* fix wrong types (was problematic for type assertions) ([2b7e7bc](https://github.com/okp4/okp4d/commit/2b7e7bcefee08315f82b82a365fbca8245ba6254))
* **lint:** add updated generated doc ([0280a12](https://github.com/okp4/okp4d/commit/0280a12ce9c7fe0700401ec105de7ce6be6948c8))
* **lint:** fix golangci lint ([cbbe5db](https://github.com/okp4/okp4d/commit/cbbe5dbecb61aad0bc46f7b3685996d6d4a9c20d))
* **lint:** remove lint error for upgrade ([d5bb919](https://github.com/okp4/okp4d/commit/d5bb91990f0d0acef5db70e743377f6dc39f66b5))
* **logic:** ensure keepers in interpreter exec ctx ([901d7f2](https://github.com/okp4/okp4d/commit/901d7f213c2b745744b46dcf72738bdbfb6c24dc))
* **logic:** fix the description inversion of the flags ([5117430](https://github.com/okp4/okp4d/commit/5117430785a6930146259f928e118b074dbece39))
* **logic:** insert bankKeeper and accountKeeper into context ([e5338c1](https://github.com/okp4/okp4d/commit/e5338c19a8871b55f622c27fad8298e9276e6061))
* **logic:** register bank_balances predicate ([3ad6317](https://github.com/okp4/okp4d/commit/3ad6317717623ef5c61f69722e22c7f09dd64413))
* **logic:** remove gocognit linter for tests ([28904d4](https://github.com/okp4/okp4d/commit/28904d49008d897581b90a797e5ee5ec8b32bc5f))
* **logic:** sort result for locked coin denom ([bc5e867](https://github.com/okp4/okp4d/commit/bc5e8676f493ee4dc853d1aaedafa4533dfd1451))
* **logic:** typo in doc comment ([d098256](https://github.com/okp4/okp4d/commit/d098256ed7c81d29b17bf413aa2448d6a42db004))
* **proto:** align naming of a proto field in yaml marshal ([4fe9a67](https://github.com/okp4/okp4d/commit/4fe9a672f9d881b39394daf3b1a44b471e10f2c9))

## [3.0.0](https://github.com/okp4/okp4d/compare/v2.2.0...v3.0.0) (2022-11-30)


### ⚠ BREAKING CHANGES

* **mint:** configure annual provision and target supply on first block

### Features

* **docs:** trigger docs version update workflow on docs repo ([a224a10](https://github.com/okp4/okp4d/commit/a224a106a75188b00a62d8c936cca67da2629890))
* **docs:** trigger the docs workflow to update documentation ([e0558aa](https://github.com/okp4/okp4d/commit/e0558aa62ecd58ef46dccc0677b8ec91e026a1a8))
* **mint:** add target_supply on proto ([3c198f7](https://github.com/okp4/okp4d/commit/3c198f7fd1b7794c6eeb197934f677bb6bde339b))
* **mint:** configure annual provision and target supply on first block ([31d5884](https://github.com/okp4/okp4d/commit/31d5884fc0a7ebbbdf23ce6921726f0489775359))
* **mint:** implement inflation calculation ([42bfa4c](https://github.com/okp4/okp4d/commit/42bfa4c8fb5de4bb269970bf852a18708d96aeff))
* **mint:** move annual reduction factor from minter to minter params ([f731e66](https://github.com/okp4/okp4d/commit/f731e6699fcae0c9c48977d5641a25182c976c2f))
* **mint:** remove okp4 old inflation calc func ([9956d1b](https://github.com/okp4/okp4d/commit/9956d1b5c8fe8e0f06ed4aacf8abc799ba322d5a))
* **mint:** set mint param on proto ([ade514e](https://github.com/okp4/okp4d/commit/ade514ecf383430dd0dfeaf283b4882659d3afcc))
* **mint:** use local proto ([cbf22f6](https://github.com/okp4/okp4d/commit/cbf22f60eff9e2bfed4ca1945a9ee24e56c0bb0d))
* **mint:** use own mint module ([af1386e](https://github.com/okp4/okp4d/commit/af1386e01a9750d0807c0291507de6751345344d))


### Bug Fixes

* **docs:** change comments syntax in markdown template ([9e8496b](https://github.com/okp4/okp4d/commit/9e8496b6279bd5361587027cd0f36d0230980e1e))
* **docs:** fix linting issue ([13709c4](https://github.com/okp4/okp4d/commit/13709c4011473e8e11f1040ca4b688da4b45a463))
* **docs:** ignore linting of generated protobuf docs ([84aaab2](https://github.com/okp4/okp4d/commit/84aaab20bc0870f71efbb3095cce10fd44078e20))
* **mint:** avoid return negative coin ([b25b1f3](https://github.com/okp4/okp4d/commit/b25b1f3b1054c4bb53a940c2eb672a594c47e384))
* **mint:** make linter more happy ([35bed9f](https://github.com/okp4/okp4d/commit/35bed9f864c19229736bf4b71e37cd3f04a2f7fe))
* **mint:** spelling mistake ([d043d1b](https://github.com/okp4/okp4d/commit/d043d1b49b83028602c10c586500065bb34b34c3))

## [2.2.0](https://github.com/okp4/okp4d/compare/v2.1.1...v2.2.0) (2022-10-13)


### Features

* **ledger:** add build dependancies and follow sdk standards ([dcd4135](https://github.com/okp4/okp4d/commit/dcd41350615132d983edf40ea0e8a490f0f5589f))
* **ledger:** bump ledger to v0.9.3 ([bff91ba](https://github.com/okp4/okp4d/commit/bff91baf73aa839f2beb33d43e9568ad7122daf4))
* **ledger:** fix install dependancies in dockerfile ([5f90752](https://github.com/okp4/okp4d/commit/5f9075231ff4e195dacd49b56a6402ee3acccd88))
* **ledger:** update ledger-go to support Ledger Nano S Plus ([4dc7f4d](https://github.com/okp4/okp4d/commit/4dc7f4d08d09f1b864010ab89df286bf80cf3ee8))

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
