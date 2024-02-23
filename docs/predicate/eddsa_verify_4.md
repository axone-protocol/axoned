---
sidebar_position: 11
---
[//]: # (This file is auto-generated. Please do not modify it yourself.)

# eddsa_verify/4

## Description

`eddsa_verify/4` determines if a given signature is valid as per the EdDSA algorithm for the provided data, using the specified public key.

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

## Examples

```text
# Verify a signature for a given hexadecimal data.
- eddsa_verify([127, ...], '9b038f8ef6918cbb56040dfda401b56b...', [23, 56, ...], [encoding(hex), type(ed25519)])

# Verify a signature for binary data.
- eddsa_verify([127, ...], [56, 90, ..], [23, 56, ...], [encoding(octet), type(ed25519)])
```
