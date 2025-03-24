package wire

import (
	"bytes"
	"reflect"
	"testing"
)

func TestAssetAllocationType_SerializeDeserialize(t *testing.T) {
	original := AssetAllocationType{
		VoutAssets: []AssetOutType{
			{
				AssetGuid: 123456789,
				Values: []AssetOutValueType{
					{N: 1, ValueSat: 5000000000},
					{N: 2, ValueSat: 2500000000},
				},
			},
			{
				AssetGuid: 987654321,
				Values: []AssetOutValueType{
					{N: 3, ValueSat: 1000000000},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := original.Serialize(&buf); err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	var deserialized AssetAllocationType
	if err := deserialized.Deserialize(&buf); err != nil {
		t.Fatalf("Deserialize failed: %v", err)
	}

	if !reflect.DeepEqual(original, deserialized) {
		t.Errorf("Mismatch after deserialize. Got %+v, want %+v", deserialized, original)
	}
}


func TestMintSyscoinType_SerializeDeserialize(t *testing.T) {
	original := MintSyscoinType{
		Allocation: AssetAllocationType{
			VoutAssets: []AssetOutType{{
				AssetGuid: 999,
				Values:    []AssetOutValueType{{N: 1, ValueSat: 123456}},
			}},
		},
		TxHash:             randomBytes(HASH_SIZE),
		BlockHash:          randomBytes(HASH_SIZE),
		TxPos:              65535,
		TxParentNodes:      randomBytes(MAX_RLP_SIZE),
		TxPath:             randomBytes(MAX_RLP_SIZE),
		TxRoot:             randomBytes(HASH_SIZE),
		ReceiptRoot:        randomBytes(HASH_SIZE),
		ReceiptPos:         65535,
		ReceiptParentNodes: randomBytes(MAX_RLP_SIZE),
	}

	var buf bytes.Buffer
	if err := original.Serialize(&buf); err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	var deserialized MintSyscoinType
	if err := deserialized.Deserialize(&buf); err != nil {
		t.Fatalf("Deserialize failed: %v", err)
	}

	if !reflect.DeepEqual(original, deserialized) {
		t.Errorf("Mismatch after deserialize. Got %+v, want %+v", deserialized, original)
	}
}

func TestSyscoinBurnToEthereumType_SerializeDeserialize(t *testing.T) {
	original := SyscoinBurnToEthereumType{
		Allocation: AssetAllocationType{
			VoutAssets: []AssetOutType{{
				AssetGuid: 888,
				Values:    []AssetOutValueType{{N: 2, ValueSat: 654321}},
			}},
		},
		EthAddress: randomBytes(MAX_GUID_LENGTH),
	}

	var buf bytes.Buffer
	if err := original.Serialize(&buf); err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	var deserialized SyscoinBurnToEthereumType
	if err := deserialized.Deserialize(&buf); err != nil {
		t.Fatalf("Deserialize failed: %v", err)
	}

	if !reflect.DeepEqual(original, deserialized) {
		t.Errorf("Mismatch after deserialize. Got %+v, want %+v", deserialized, original)
	}
}
func TestAssetType_SerializeDeserialize(t *testing.T) {
    original := AssetType{
        Symbol:    []byte("SYS"),
        Precision: 8,
    }

    var buf bytes.Buffer
    if err := original.Serialize(&buf); err != nil {
        t.Fatalf("Serialize failed: %v", err)
    }

    var deserialized AssetType
    if err := deserialized.Deserialize(&buf); err != nil {
        t.Fatalf("Deserialize failed: %v", err)
    }

    if !bytes.Equal(original.Symbol, deserialized.Symbol) || original.Precision != deserialized.Precision {
        t.Errorf("Mismatch after deserialize. Got %+v, want %+v", deserialized, original)
    }
}
