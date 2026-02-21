---
sidebar_position: 14
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# crypto_data_hash/3

## Description

`crypto_data_hash/3` is a predicate that computes the Hash of the given Data using different algorithms.

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

## Examples

```text
# Compute the SHA-256 hash of the given data and unify it with the given Hash.
- crypto_data_hash('Hello AXONE', Hash).

# Compute the SHA-256 hash of the given hexadecimal data and unify it with the given Hash.
- crypto_data_hash('9b038f8ef6918cbb56040dfda401b56b...', Hash, encoding(hex)).

# Compute the SHA-256 hash of the given hexadecimal data and unify it with the given Hash.
- crypto_data_hash([127, ...], Hash, encoding(octet)).
```
