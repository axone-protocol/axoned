[//]: # (This file is auto-generated. Please do not modify it yourself.)

# Predicates documentation

## bank_balances/2

bank_balances/2 is a predicate which unifies the given terms with the list of balances \(coins\) of the given account.

The signature is as follows:

```text
bank_balances(?Account, ?Balances)
```

where:

- Account represents the account address \(in Bech32 format\).
- Balances represents the balances of the account as a list of pairs of coin denomination and amount.

Examples:

```text
# Query the balances of the account.
- bank_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', X).

# Query the balances of all accounts. The result is a list of pairs of account address and balances.
- bank_balances(X, Y).

# Query the first balance of the given account by unifying the denomination and amount with the given terms.
- bank_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [-(D, A), _]).
```

## bank_locked_balances/2

bank_locked_balances/2 is a predicate which unifies the given terms with the list of locked coins of the given account.

The signature is as follows:

```text
bank_locked_balances(?Account, ?Balances)
```

where:

- Account represents the account address \(in Bech32 format\).
- Balances represents the locked balances of the account as a list of pairs of coin denomination and amount.

Examples:

```text
# Query the locked coins of the account.
- bank_locked_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', X).

# Query the locked balances of all accounts. The result is a list of pairs of account address and balances.
- bank_locked_balances(X, Y).

# Query the first locked balances of the given account by unifying the denomination and amount with the given terms.
- bank_locked_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [-(D, A), _]).
```

## bank_spendable_balances/2

bank_spendable_balances/2 is a predicate which unifies the given terms with the list of spendable coins of the given account.

The signature is as follows:

```text
bank_spendable_balances(?Account, ?Balances)
```

where:

- Account represents the account address \(in Bech32 format\).
- Balances represents the spendable balances of the account as a list of pairs of coin denomination and amount.

Examples:

```text
# Query the spendable balances of the account.
- bank_spendable_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', X).

# Query the spendable balances of all accounts. The result is a list of pairs of account address and balances.
- bank_spendable_balances(X, Y).

# Query the first spendable balances of the given account by unifying the denomination and amount with the given terms.
- bank_spendable_balances('okp41ffd5wx65l407yvm478cxzlgygw07h79sq0m3fm', [-(D, A), _]).
```

## bech32_address/2

bech32_address/2 is a predicate that convert a [bech32](<https://docs.cosmos.network/main/build/spec/addresses/bech32#hrp-table>) encoded string into [base64](<https://fr.wikipedia.org/wiki/Base64>) bytes and give the address prefix, or convert a prefix \(HRP\) and [base64](<https://fr.wikipedia.org/wiki/Base64>) encoded bytes to [bech32](<https://docs.cosmos.network/main/build/spec/addresses/bech32#hrp-table>) encoded string.

The signature is as follows:

```text
bech32_address(-Address, +Bech32)
bech32_address(+Address, -Bech32)
bech32_address(+Address, +Bech32)
```

where:

- Address is a pair of the HRP \(Human\-Readable Part\) which holds the address prefix and a list of numbers ranging from 0 to 255 that represent the base64 encoded bech32 address string.
- Bech32 is an Atom or string representing the bech32 encoded string address

Examples:

```text
# Convert the given bech32 address into base64 encoded byte by unify the prefix of given address (Hrp) and the
base64 encoded value (Address).
- bech32_address(-(Hrp, Address), 'okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn').

# Convert the given pair of HRP and base64 encoded address byte by unify the Bech32 string encoded value.
- bech32_address(-('okp4', [163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]), Bech32).
```

## block_height/1

block_height/1 is a predicate which unifies the given term with the current block height.

The signature is as follows:

```text
block_height(?Height)
```

where:

- Height represents the current chain height at the time of the query.

Examples:

```text
# Query the current block height.
- block_height(Height).
```

## block_time/1

block_time/1 is a predicate which unifies the given term with the current block time.

The signature is as follows:

```text
block_time(?Time)
```

where:

- Time represents the current chain time at the time of the query.

Examples:

```text
# Query the current block time.
- block_time(Time).
```

## chain_id/1

chain_id/1 is a predicate which unifies the given term with the current chain ID. The signature is:

The signature is as follows:

```text
chain_id(?ID)
```

where:

- ID represents the current chain ID at the time of the query.

Examples:

```text
# Query the current chain ID.
- chain_id(ID).
```

## crypto_data_hash/3

crypto_data_hash/3 is a predicate that computes the Hash of the given Data using different algorithms.

The signature is as follows:

```text
crypto_data_hash(+Data, -Hash, +Options) is det
crypto_data_hash(+Data, +Hash, +Options) is det
```

Where:

- Data represents the data to be hashed, given as an atom, or code\-list.
- Hash represents the Hashed value of Data, which can be given as an atom or a variable.
- Options are additional configurations for the hashing process. Supported options include: encoding\(\+Format\) which specifies the encoding used for the Data, and algorithm\(\+Alg\) which chooses the hashing algorithm among the supported ones \(see below for details\).

For Format, the supported encodings are:

- utf8 \(default\), the UTF\-8 encoding represented as an atom.
- text, the plain text encoding represented as an atom.
- hex, the hexadecimal encoding represented as an atom.
- octet, the raw byte encoding depicted as a list of integers ranging from 0 to 255.

For Alg, the supported algorithms are:

- sha256 \(default\): The SHA\-256 algorithm.
- sha512: The SHA\-512 algorithm.
- md5: \(insecure\) The MD5 algorithm.

Note: Due to the principles of the hash algorithm \(pre\-image resistance\), this predicate can only compute the hash value from input data, and cannot compute the original input data from the hash value.

Examples:

```text
# Compute the SHA-256 hash of the given data and unify it with the given Hash.
- crypto_data_hash('Hello OKP4', Hash).

# Compute the SHA-256 hash of the given hexadecimal data and unify it with the given Hash.
- crypto_data_hash('9b038f8ef6918cbb56040dfda401b56b...', Hash, encoding(hex)).

# Compute the SHA-256 hash of the given hexadecimal data and unify it with the given Hash.
- crypto_data_hash([127, ...], Hash, encoding(octet)).
```

## did_components/2

did_components/2 is a predicate which breaks down a DID into its components according to the [W3C DID](<https://w3c.github.io/did-core>) specification.

The signature is as follows:

```text
did_components(+DID, -Components) is det
did_components(-DID, +Components) is det
```

where:

- DID represent DID URI, given as an Atom, compliant with [W3C DID](<https://w3c.github.io/did-core>) specification.
- Components is a compound Term in the format did\(Method, ID, Path, Query, Fragment\), aligned with the [DID syntax](<https://w3c.github.io/did-core/#did-syntax>), where: Method is the method name, ID is the method\-specific identifier, Path is the path component, Query is the query component and Fragment is the fragment component. Values are given as an Atom and are url encoded. For any component not present, its value will be null and thus will be left as an uninstantiated variable.

Examples:

```text
# Decompose a DID into its components.
- did_components('did:example:123456?versionId=1', did_components(Method, ID, Path, Query, Fragment)).

# Reconstruct a DID from its components.
- did_components(DID, did_components('example', '123456', _, 'versionId=1', _42)).
```

## ecdsa_verify/4

ecdsa_verify/4 determines if a given signature is valid as per the ECDSA algorithm for the provided data, using the specified public key.

The signature is as follows:

```text
ecdsa_verify(+PubKey, +Data, +Signature, +Options), which is semi-deterministic.
```

Where:

- PubKey is the 33\-byte compressed public key, as specified in section 4.3.6 of ANSI X9.62.

- Data is the hash of the signed message, which can be either an atom or a list of bytes.

- Signature represents the ASN.1 encoded signature corresponding to the Data.

- Options are additional configurations for the verification process. Supported options include: encoding\(\+Format\) which specifies the encoding used for the data, and type\(\+Alg\) which chooses the algorithm within the ECDSA family \(see below for details\).

For Format, the supported encodings are:

- hex \(default\), the hexadecimal encoding represented as an atom.
- octet, the plain byte encoding depicted as a list of integers ranging from 0 to 255.
- text, the plain text encoding represented as an atom.
- utf8 \(default\), the UTF\-8 encoding represented as an atom.

For Alg, the supported algorithms are:

- secp256r1 \(default\): Also known as P\-256 and prime256v1.
- secp256k1: The Koblitz elliptic curve used in Bitcoin's public\-key cryptography.

Examples:

```text
# Verify a signature for hexadecimal data using the ECDSA secp256r1 algorithm.
- ecdsa_verify([127, ...], '9b038f8ef6918cbb56040dfda401b56b...', [23, 56, ...], encoding(hex))

# Verify a signature for binary data using the ECDSA secp256k1 algorithm.
- ecdsa_verify([127, ...], [56, 90, ..], [23, 56, ...], [encoding(octet), type(secp256k1)])
```

## eddsa_verify/4

eddsa_verify/4 determines if a given signature is valid as per the EdDSA algorithm for the provided data, using the specified public key.

The signature is as follows:

```text
eddsa_verify(+PubKey, +Data, +Signature, +Options) is semi-det
```

Where:

- PubKey is the encoded public key as a list of bytes.
- Data is the message to verify, represented as either a hexadecimal atom or a list of bytes. It's important that the message isn't pre\-hashed since the Ed25519 algorithm processes messages in two passes when signing.
- Signature represents the signature corresponding to the data, provided as a list of bytes.
- Options are additional configurations for the verification process. Supported options include: encoding\(\+Format\) which specifies the encoding used for the Data, and type\(\+Alg\) which chooses the algorithm within the EdDSA family \(see below for details\).

For Format, the supported encodings are:

- hex \(default\), the hexadecimal encoding represented as an atom.
- octet, the plain byte encoding depicted as a list of integers ranging from 0 to 255.
- text, the plain text encoding represented as an atom.
- utf8 \(default\), the UTF\-8 encoding represented as an atom.

For Alg, the supported algorithms are:

- ed25519 \(default\): The EdDSA signature scheme using SHA\-512 \(SHA\-2\) and Curve25519.

Examples:

```text
# Verify a signature for a given hexadecimal data.
- eddsa_verify([127, ...], '9b038f8ef6918cbb56040dfda401b56b...', [23, 56, ...], [encoding(hex), type(ed25519)])

# Verify a signature for binary data.
- eddsa_verify([127, ...], [56, 90, ..], [23, 56, ...], [encoding(octet), type(ed25519)])
```

## hex_bytes/2

hex_bytes/2 is a predicate that unifies hexadecimal encoded bytes to a list of bytes.

The signature is as follows:

```text
hex_bytes(?Hex, ?Bytes) is det
```

Where:

- Hex is an Atom, string or list of characters in hexadecimal encoding.
- Bytes is the list of numbers between 0 and 255 that represent the sequence of bytes.

Examples:

```text
# Convert hexadecimal atom to list of bytes.
- hex_bytes('2c26b46b68ffc68ff99b453c1d3041341342d706483bfa0f98a5e886266e7ae', Bytes).
```

## json_prolog/2

json_prolog/2 is a predicate that will unify a JSON string into prolog terms and vice versa.

The signature is as follows:

```text
json_prolog(?Json, ?Term) is det
```

Where:

- Json is the string representation of the json
- Term is an Atom that would be unified by the JSON representation as Prolog terms.

In addition, when passing Json and Term, this predicate return true if both result match.

Examples:

```text
# JSON conversion to Prolog.
- json_prolog('{"foo": "bar"}', json([foo-bar])).
```

## open/4

open/4 is a predicate that unify a stream with a source sink on a virtual file system.

The signature is as follows:

```text
open(+SourceSink, +Mode, ?Stream, +Options)
```

Where:

- SourceSink is an atom representing the source or sink of the stream. The atom typically represents a resource that can be opened, such as a URI. The URI scheme determines the type of resource that is opened.
- Mode is an atom representing the mode of the stream \(read, write, append\).
- Stream is the stream to be opened.
- Options is a list of options. No options are currently defined, so the list should be empty.

Examples:

```text
# open/4 a stream from a cosmwasm query.
# The Stream should be read as a string with a read_string/3 predicate, and then closed with the close/1 predicate.
- open('cosmwasm:okp4-objectarium:okp412kgx?query=%7B%22object_data%22%3A%7B%...4dd539e3%22%7D%7D', 'read', Stream, [])
```

## read_string/3

read_string/3 is a predicate that reads characters from the provided Stream and unifies them with String. Users can optionally specify a maximum length for reading; if the stream reaches this length, the reading stops. If Length remains unbound, the entire Stream is read, and upon completion, Length is unified with the count of characters read.

The signature is as follows:

```text
read_string(+Stream, ?Length, -String) is det
```

Where:

- Stream is the input stream to read from.
- Length is the optional maximum number of characters to read from the Stream. If unbound, denotes the full length of Stream.
- String is the resultant string after reading from the Stream.

Examples:

```text
# Given a file `foo.txt` that contains `Hello World`:

file_to_string(File, String, Length) :-

open(File, read, In),
read_string(In, Length, String),
close(Stream).

# It gives:
?- file_to_string('path/file/foo.txt', String, Length).

String = 'Hello World'
Length = 11
```

## source_file/1

source_file/1 is a predicate that unify the given term with the currently loaded source file.

The signature is as follows:

```text
source_file(?File).
```

Where:

- File represents a loaded source file.

Examples:

```text
# Query all the loaded source files, in alphanumeric order.
- source_file(File).

# Query the given source file is loaded.
- source_file('foo.pl').
```

## string_bytes/3

string_bytes/3 is a predicate that unifies a string with a list of bytes, returning true when the \(Unicode\) String is represented by Bytes in Encoding.

The signature is as follows:

```text
string_bytes(?String, ?Bytes, +Encoding)
```

Where:

- String is the string to convert to bytes. It can be an Atom, string or list of characters codes.
- Bytes is the list of numbers between 0 and 255 that represent the sequence of bytes.
- Encoding is the encoding to use for the conversion.

Encoding can be one of the following: \- 'text' considers the string as a sequence of Unicode characters. \- 'octet' considers the string as a sequence of bytes. \- '\<encoding\>' considers the string as a sequence of characters in the given encoding.

At least one of String or Bytes must be instantiated.

Examples:

```text
# Convert a string to a list of bytes.
- string_bytes('Hello World', Bytes, octet).

# Convert a list of bytes to a string.
- string_bytes(String, [72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100], octet).
```

## uri_encoded/3

uri_encoded/3 is a predicate that unifies the given URI component with the given encoded or decoded string.

The signature is as follows:

```text
uri_encoded(+Component, +Value, -Encoded) is det
uri_encoded(+Component, -Value, +Encoded) is det
```

Where:

- Component represents the component of the URI to be escaped. It can be the atom query, fragment, path or segment.
- Decoded represents the decoded string to be escaped.
- Encoded represents the encoded string.

For more information on URI encoding, refer to [RFC 3986](<https://datatracker.ietf.org/doc/html/rfc3986#section-2.1>).

Examples:

```text
# Escape the given string to be used in the path component.
- uri_encoded(path, "foo/bar", Encoded).

# Unescape the given string to be used in the path component.
- uri_encoded(path, Decoded, foo%2Fbar).
```
