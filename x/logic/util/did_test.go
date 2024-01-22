package util

import (
	"encoding/hex"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateDIDKeyByPubKey(t *testing.T) {
	Convey("Given test cases", t, func() {
		tests := []struct {
			name       string
			pubKey     string
			algo       KeyAlg
			wantResult string
			wantError  error
		}{
			{
				name:       "did:key from secp256k1 pubkey",
				pubKey:     "02d0fe99b214aaeeb5e46ae4d65a1623a95d9d0cecd57d673f7e4c9a19ac0752bc",
				algo:       KeyAlgSecp256k1,
				wantResult: "did:key:zQ3shbUcjA9ptj3HaEWV8PCJPsCHA62YUWyrcmFhuWG4ciAbD",
			},
			{
				name:       "did:key from ed25519 pubkey",
				pubKey:     "ec437b6f607c95267597b76f65317446b4cba6a400254f84027b7c99cdd1ab27",
				algo:       KeyAlgEd25519,
				wantResult: "did:key:z6MkvMXwgwJTJacfBGk5fxr3d4k3uzh4eHTi3oFagNyK55Tt",
			},
			{
				name:      "did:key with incorrect ed25519 pubkey",
				pubKey:    "02d0fe99b214aaeeb5e46ae4d65a1623a95d9d0cecd57d673f7e4c9a19ac0752bc",
				algo:      KeyAlgEd25519,
				wantError: fmt.Errorf("invalid pubkey size; expected 32, got 33"),
			},
			{
				name:      "did:key with incorrect secp256k1 pubkey",
				pubKey:    "ec437b6f607c95267597b76f65317446b4cba6a400254f84027b7c99cdd1ab27",
				algo:      KeyAlgSecp256k1,
				wantError: fmt.Errorf("invalid pubkey size; expected 33, got 32"),
			},
			{
				name:      "did:key with unsupported pubkey (for now)",
				pubKey:    "ec437b6f607c95267597b76f65317446b4cba6a400254f84027b7c99cdd1ab27",
				algo:      KeyAlgSecp256r1,
				wantError: fmt.Errorf("invalid pubkey type: secp256r1; expected oneof [\"secp256k1\" \"ed25519\"]"),
			},
		}
		for _, tt := range tests {
			Convey(fmt.Sprintf("Given test case %s", tt.name), func() {
				pubKeyBs, err := hex.DecodeString(tt.pubKey)
				So(err, ShouldBeNil)
				So(pubKeyBs, ShouldNotBeNil)

				Convey("When calling the DIDKeyByPubKey function", func() {
					var got string
					pubKey, err := BytesToPubKey(pubKeyBs, tt.algo)
					if err == nil {
						got, err = CreateDIDKeyByPubKey(pubKey)
					}

					Convey("Then we should get the expected result", func() {
						if tt.wantError != nil {
							So(err, ShouldNotBeNil)
							So(err.Error(), ShouldEqual, tt.wantError.Error())
						} else {
							So(err, ShouldBeNil)
							So(got, ShouldEqual, tt.wantResult)
						}
					})
				})
			})
		}
	})
}
