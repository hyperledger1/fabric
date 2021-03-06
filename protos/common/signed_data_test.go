/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"bytes"
	"testing"

	"github.com/golang/protobuf/proto"
)

// More duplicate utility which should go away, but the utils are a bit of a mess right now with import cycles
func marshalOrPanic(msg proto.Message) []byte {
	data, err := proto.Marshal(msg)
	if err != nil {
		panic("Error marshaling")
	}
	return data
}

func TestNilSignedConfigurationItemAsSignedData(t *testing.T) {
	var sci *SignedConfigurationItem
	_, err := sci.AsSignedData()
	if err == nil {
		t.Fatalf("Should have errored trying to convert a nil signed configuration item to signed data")
	}
}

func TestSignedConfigurationItemAsSignedData(t *testing.T) {
	configItemBytes := []byte("Foo")
	signatures := [][]byte{[]byte("Signature1"), []byte("Signature2")}
	identities := [][]byte{[]byte("Identity1"), []byte("Identity2")}

	configSignatures := make([]*ConfigurationSignature, len(signatures))
	for i := range configSignatures {
		configSignatures[i] = &ConfigurationSignature{
			SignatureHeader: marshalOrPanic(&SignatureHeader{
				Creator: identities[i],
			}),
			Signature: signatures[i],
		}
	}

	sci := &SignedConfigurationItem{
		ConfigurationItem: configItemBytes,
		Signatures:        configSignatures,
	}

	signedData, err := sci.AsSignedData()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	for i, sigData := range signedData {
		if !bytes.Equal(sigData.Identity, identities[i]) {
			t.Errorf("Expected identity to match at index %d", i)
		}
		if !bytes.Equal(sigData.Data, append(configItemBytes, configSignatures[i].SignatureHeader...)) {
			t.Errorf("Expected signature over concatenation of config item bytes and signature header")
		}
		if !bytes.Equal(sigData.Signature, signatures[i]) {
			t.Errorf("Expected signature to match at index %d", i)
		}
	}
}

func TestNilEnvelopeAsSignedData(t *testing.T) {
	var env *Envelope
	_, err := env.AsSignedData()
	if err == nil {
		t.Fatalf("Should have errored trying to convert a nil envelope")
	}
}

func TestEnvelopeAsSignedData(t *testing.T) {
	identity := []byte("Foo")
	signature := []byte("Bar")
	env := &Envelope{
		Payload: marshalOrPanic(&Payload{
			Header: &Header{
				SignatureHeader: &SignatureHeader{
					Creator: identity,
				},
			},
		}),
		Signature: signature,
	}

	signedData, err := env.AsSignedData()
	if err != nil {
		t.Fatalf("Unexpected error converting envelope to SignedData: %s", err)
	}

	if len(signedData) != 1 {
		t.Fatalf("Expected 1 entry of signed data, but got %d", len(signedData))
	}

	if !bytes.Equal(signedData[0].Identity, identity) {
		t.Errorf("Wrong identity bytes")
	}
	if !bytes.Equal(signedData[0].Data, env.Payload) {
		t.Errorf("Wrong data bytes")
	}
	if !bytes.Equal(signedData[0].Signature, signature) {
		t.Errorf("Wrong data bytes")
	}
}
