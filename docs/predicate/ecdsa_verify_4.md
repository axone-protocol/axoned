---
sidebar_position: 11
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# ecdsa_verify/4

## Description

`ecdsa_verify/4` determines if a given signature is valid as per the ECDSA algorithm for the provided data, using the specified public key.

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

## Examples

```text
# Verify a signature for hexadecimal data using the ECDSA secp256r1 algorithm.
- ecdsa_verify([127, ...], '9b038f8ef6918cbb56040dfda401b56b...', [23, 56, ...], encoding(hex))

# Verify a signature for binary data using the ECDSA secp256k1 algorithm.
- ecdsa_verify([127, ...], [56, 90, ..], [23, 56, ...], [encoding(octet), type(secp256k1)])
```
