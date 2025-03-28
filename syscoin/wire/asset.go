// Copyright (c) 2013-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
    "io"
    "encoding/binary"
    "github.com/btcsuite/btcd/wire"
)
const (
    MAX_GUID_LENGTH = 20
    MAX_RLP_SIZE = 4096
)
type AssetOutValueType struct {
    N uint32
    ValueSat int64
}
type AssetOutType struct {
    AssetGuid uint64
    Values []AssetOutValueType
}
type AssetAllocationType struct {
    VoutAssets []AssetOutType
}
type AssetType struct {
    Contract []byte
    Symbol []byte
    TotalSupply int64
    MaxSupply int64
    Precision uint8
}
type MintSyscoinType struct {
    Allocation AssetAllocationType
    TxHash []byte
    BlockHash []byte
    TxPos uint16
    TxParentNodes []byte
    TxPath []byte
    TxRoot []byte
    ReceiptRoot []byte
    ReceiptPos uint16
    ReceiptParentNodes []byte
}

type SyscoinBurnToEthereumType struct {
    Allocation AssetAllocationType
    EthAddress []byte `json:"ethAddress,omitempty"`
}

func PutUint(w io.Writer, n uint64) error {
    tmp := make([]uint8, 10)
    var len int8=0
    for  {
        var mask uint64
        if len > 0 {
            mask = 0x80
        }
        tmpI := (n & 0x7F) | mask
        tmp[len] = uint8(tmpI)
        if n <= 0x7F {
            break
        }
        n = (n >> 7) - 1
        len++
    }
    for {
        err := binarySerializer.PutUint8(w, tmp[len])
        if err != nil {
            return err
        }
        len--
        if len < 0 {
            break
        }
    }
    return nil
}
func ReadUint(r io.Reader) (uint64, error) {
    var n uint64 = 0
    for {
        chData, err := binarySerializer.Uint8(r)
        if err != nil {
            return 0, err
        }
        n = (n << 7) | (uint64(chData) & 0x7F)
        if (chData & 0x80) > 0 {
            n++
        } else {
            return n, nil
        }
    }
}
// Amount compression:
// * If the amount is 0, output 0
// * first, divide the amount (in base units) by the largest power of 10 possible; call the exponent e (e is max 9)
// * if e<9, the last digit of the resulting number cannot be 0; store it as d, and drop it (divide by 10)
//   * call the result n
//   * output 1 + 10*(9*n + d - 1) + e
// * if e==9, we only know the resulting number is not zero, so output 1 + 10*(n - 1) + 9
// (this is decodable, as d is in [1-9] and e is in [0-9])

func CompressAmount(n uint64) uint64 {
    if n == 0 {
        return 0
    }
    var e int = 0;
    for ((n % 10) == 0) && e < 9 {
        n /= 10
        e++
    }
    if e < 9 {
        var d int = int(n % 10)
        n /= 10
        return 1 + (n*9 + uint64(d) - 1)*10 + uint64(e)
    } else {
        return 1 + (n - 1)*10 + 9
    }
}

func DecompressAmount(x uint64) uint64 {
    // x = 0  OR  x = 1+10*(9*n + d - 1) + e  OR  x = 1+10*(n - 1) + 9
    if x == 0 {
        return 0
    }
    x--
    // x = 10*(9*n + d - 1) + e
    var e int = int(x % 10)
    x /= 10
    var n uint64 = 0
    if e < 9 {
        // x = 9*n + d - 1
        var d int = int(x % 9) + 1
        x /= 9
        // x = n
        n = x*10 + uint64(d)
    } else {
        n = x+1
    }
    for e > 0 {
        n *= 10
        e--
    }
    return n
}
func (a *AssetAllocationType) Deserialize(r io.Reader) error {
    numAssets, err := wire.ReadVarInt(r, 0)
    if err != nil {
        return err
    }
    a.VoutAssets = make([]AssetOutType, numAssets)
    for i := 0; i < int(numAssets); i++ {
        err = a.VoutAssets[i].Deserialize(r)
        if err != nil {
            return err
        }
    }
    return nil
}
func (a *AssetAllocationType) Serialize(w io.Writer) error {
    lenAssets := len(a.VoutAssets)
    err := wire.WriteVarInt(w, 0, uint64(lenAssets))
    if err != nil {
        return err
    }
    for i := 0; i < lenAssets; i++ {
        err = a.VoutAssets[i].Serialize(w)
        if err != nil {
            return err
        }
    }
    return nil
}
func (a *AssetOutValueType) Serialize(w io.Writer) error {
    err := wire.WriteVarInt(w, 0, uint64(a.N))
    if err != nil {
        return err
    }
    err = PutUint(w, CompressAmount(uint64(a.ValueSat)))
    if err != nil {
        return err
    }
    return nil
}
func (a *AssetOutValueType) Deserialize(r io.Reader) error {
    n, err := wire.ReadVarInt(r, 0)
    if err != nil {
        return err
    }
    a.N = uint32(n)
    valueSat, err := ReadUint(r)
    if err != nil {
        return err
    }
    a.ValueSat = int64(DecompressAmount(valueSat))
    return nil
}
func (a *AssetOutType) Serialize(w io.Writer) error {
    err := PutUint(w, a.AssetGuid)
    if err != nil {
        return err
    }
    lenValues := len(a.Values)
    err = wire.WriteVarInt(w, 0, uint64(lenValues))
    if err != nil {
        return err
    }
    for i := 0; i < lenValues; i++ {
        err = a.Values[i].Serialize(w)
        if err != nil {
            return err
        }
    }
    return nil
}
func (a *AssetOutType) Deserialize(r io.Reader) error {
    var err error
    a.AssetGuid, err = ReadUint(r)
    if err != nil {
        return err
    }
    numOutputs, err := wire.ReadVarInt(r, 0)
    if err != nil {
        return err
    }
    a.Values = make([]AssetOutValueType, numOutputs)
    for i := 0; i < int(numOutputs); i++ {
        err = a.Values[i].Deserialize(r)
        if err != nil {
            return err
        }
    }
    return nil
}

func (a *MintSyscoinType) Deserialize(r io.Reader) error {
    err := a.Allocation.Deserialize(r)
    if err != nil {
        return err
    }
    a.TxHash = make([]byte, HASH_SIZE)
    _, err = io.ReadFull(r, a.TxHash)
    if err != nil {
        return err
    }
    a.BlockHash = make([]byte, HASH_SIZE)
    _, err = io.ReadFull(r, a.BlockHash)
    if err != nil {
        return err
    }
    a.TxPos, err = binarySerializer.Uint16(r, binary.LittleEndian)
    if err != nil {
        return err
    }
    a.TxParentNodes, err = wire.ReadVarBytes(r, 0, MAX_RLP_SIZE, "TxParentNodes")
    if err != nil {
        return err
    }
    a.TxPath, err = wire.ReadVarBytes(r, 0, MAX_RLP_SIZE, "TxPath")
    if err != nil {
        return err
    }
    a.ReceiptPos, err = binarySerializer.Uint16(r, binary.LittleEndian)
    if err != nil {
        return err
    }
    a.ReceiptParentNodes, err = wire.ReadVarBytes(r, 0, MAX_RLP_SIZE, "ReceiptParentNodes")
    if err != nil {
        return err
    }
    a.TxRoot = make([]byte, HASH_SIZE)
    _, err = io.ReadFull(r, a.TxRoot)
    if err != nil {
        return err
    }
    a.ReceiptRoot = make([]byte, HASH_SIZE)
    _, err = io.ReadFull(r, a.ReceiptRoot)
    if err != nil {
        return err
    }
    return nil
}

func (a *MintSyscoinType) Serialize(w io.Writer) error {
    err := a.Allocation.Serialize(w)
    if err != nil {
        return err
    }
    _, err = w.Write(a.TxHash[:])
    if err != nil {
        return err
    }
    _, err = w.Write(a.BlockHash[:])
    if err != nil {
        return err
    }
    err = binarySerializer.PutUint16(w, binary.LittleEndian, a.TxPos)
    if err != nil {
        return err
    }
    err = wire.WriteVarBytes(w, 0, a.TxParentNodes)
    if err != nil {
        return err
    }
    err = wire.WriteVarBytes(w, 0, a.TxPath)
    if err != nil {
        return err
    }
    err = binarySerializer.PutUint16(w, binary.LittleEndian, a.ReceiptPos)
    if err != nil {
        return err
    }
    err = wire.WriteVarBytes(w, 0, a.ReceiptParentNodes)
    if err != nil {
        return err
    }
    _, err = w.Write(a.TxRoot[:])
    if err != nil {
        return err
    }
    _, err = w.Write(a.ReceiptRoot[:])
    if err != nil {
        return err
    }
    return nil
}

func (a *SyscoinBurnToEthereumType) Deserialize(r io.Reader) error {
    err := a.Allocation.Deserialize(r)
    if err != nil {
        return err
    }
    a.EthAddress, err = wire.ReadVarBytes(r, 0, MAX_GUID_LENGTH, "ethAddress")
    if err != nil {
        return err
    }
    return nil
}

func (a *SyscoinBurnToEthereumType) Serialize(w io.Writer) error {
    err := a.Allocation.Serialize(w)
    if err != nil {
        return err
    }
    err = wire.WriteVarBytes(w, 0, a.EthAddress)
    if err != nil {
        return err
    }
    return nil
}

func (a *AssetType) Serialize(w io.Writer) error {
    // Serialize Symbol
    if err := wire.WriteVarBytes(w, 0, a.Symbol); err != nil {
        return err
    }

    // Serialize Precision
    if err := binarySerializer.PutUint8(w, a.Precision); err != nil {
        return err
    }

    return nil
}

func (a *AssetType) Deserialize(r io.Reader) error {
    var err error

    // Deserialize Symbol
    a.Symbol, err = wire.ReadVarBytes(r, 0, MAX_GUID_LENGTH, "Symbol")
    if err != nil {
        return err
    }

    // Deserialize Precision
    a.Precision, err = binarySerializer.Uint8(r)
    if err != nil {
        return err
    }

    return nil
}
