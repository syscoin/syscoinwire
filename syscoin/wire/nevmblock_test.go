package wire

import (
	"bytes"
	"crypto/rand"
	"reflect"
	"testing"
)

func randomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}


func TestNEVMBlockWire_SerializeDeserialize(t *testing.T) {
	original := NEVMBlockWire{
		NEVMBlockHash: randomBytes(HASH_SIZE),
		TxRoot:        randomBytes(HASH_SIZE),
		ReceiptRoot:   randomBytes(HASH_SIZE),
		NEVMBlockData: randomBytes(MAX_NEVM_BLOCK_SIZE),
		SYSBlockHash:  randomBytes(HASH_SIZE),
		VersionHashes: [][]byte{
			randomBytes(HASH_SIZE),
			randomBytes(HASH_SIZE),
		},
		Diff: NEVMAddressDiff{
			AddedMNNEVM: []NEVMAddressEntry{
				{Address: randomBytes(HASH_SIZE), CollateralHeight: 0},
			},
			UpdatedMNNEVM: []NEVMAddressUpdateEntry{
				{OldAddress: randomBytes(HASH_SIZE), NewAddress: randomBytes(HASH_SIZE), CollateralHeight: 4294967295},
			},
			RemovedMNNEVM: []NEVMRemoveEntry{
				{Address: randomBytes(HASH_SIZE)},
			},
		},
	}

	var buf bytes.Buffer
	if err := original.Serialize(&buf); err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	var deserialized NEVMBlockWire
	if err := deserialized.Deserialize(&buf); err != nil {
		t.Fatalf("Deserialize failed: %v", err)
	}

	if !reflect.DeepEqual(original, deserialized) {
		t.Errorf("Mismatch after deserialize. Got %+v, want %+v", deserialized, original)
	}
}

func TestNEVMDisconnectBlockWire_SerializeDeserialize(t *testing.T) {
	original := NEVMDisconnectBlockWire{
		SYSBlockHash: randomBytes(HASH_SIZE),
		Diff: NEVMAddressDiff{
			AddedMNNEVM: []NEVMAddressEntry{{Address: randomBytes(HASH_SIZE), CollateralHeight: 123456}},
			UpdatedMNNEVM: []NEVMAddressUpdateEntry{{OldAddress: randomBytes(HASH_SIZE), NewAddress: randomBytes(HASH_SIZE), CollateralHeight: 654321}},
			RemovedMNNEVM: []NEVMRemoveEntry{{Address: randomBytes(HASH_SIZE)}},
		},
	}

	var buf bytes.Buffer
	if err := original.Serialize(&buf); err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	var deserialized NEVMDisconnectBlockWire
	if err := deserialized.Deserialize(&buf); err != nil {
		t.Fatalf("Deserialize failed: %v", err)
	}

	if !reflect.DeepEqual(original, deserialized) {
		t.Errorf("Mismatch after deserialize. Got %+v, want %+v", deserialized, original)
	}
}
