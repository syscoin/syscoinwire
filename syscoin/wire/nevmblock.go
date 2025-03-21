// Copyright (c) 2013-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
    "io"
    "github.com/btcsuite/btcd/wire"
)
const (
    HASH_SIZE = 32
    MAX_NEVM_BLOCK_SIZE = 33554432 // 32 MB
)

// NEVMAddressEntry represents an entry with an address and collateralHeight.
type NEVMAddressEntry struct {
    Address          []byte
    CollateralHeight uint32
}
// NEVMAddressUpdateEntry represents an update entry with old and new addresses.
type NEVMAddressUpdateEntry struct {
    OldAddress []byte
    NewAddress []byte
    CollateralHeight uint32
}
// NEVMRemoveEntry represents an entry with an address to be removed.
type NEVMRemoveEntry struct {
    Address []byte
}
// NEVMAddressDiff holds the differences in the NEVM address list for masternodes.
type NEVMAddressDiff struct {
    AddedMNNEVM   []NEVMAddressEntry
    UpdatedMNNEVM []NEVMAddressUpdateEntry
    RemovedMNNEVM []NEVMRemoveEntry
}

type NEVMBlockWire struct {
    NEVMBlockHash []byte
    TxRoot        []byte
    ReceiptRoot   []byte
    NEVMBlockData []byte
    SYSBlockHash  []byte
    VersionHashes [][]byte
    Diff          NEVMAddressDiff
}

type NEVMDisconnectBlockWire struct {
    SYSBlockHash  []byte
    Diff          NEVMAddressDiff
}


func (a *NEVMAddressEntry) Deserialize(r io.Reader) error {
    var err error
    a.Address, err = wire.ReadVarBytes(r, 0, HASH_SIZE, "Address")
    if err != nil {
        return err
    }
    a.CollateralHeight, err = binarySerializer.Uint32(r, littleEndian)
    if err != nil {
        return err
    }
    return nil
}

func (a *NEVMAddressEntry) Serialize(w io.Writer) error {
    err := wire.WriteVarBytes(w, 0, a.Address)
    if err != nil {
        return err
    }
    err = binarySerializer.PutUint32(w, littleEndian, a.CollateralHeight)
    if err != nil {
        return err
    }
    return nil
}

func (a *NEVMAddressUpdateEntry) Deserialize(r io.Reader) error {
    var err error
    a.OldAddress, err = wire.ReadVarBytes(r, 0, HASH_SIZE, "OldAddress")
    if err != nil {
        return err
    }
    a.NewAddress, err = wire.ReadVarBytes(r, 0, HASH_SIZE, "NewAddress")
    if err != nil {
        return err
    }
    a.CollateralHeight, err = binarySerializer.Uint32(r, littleEndian)
    if err != nil {
        return err
    }
    return nil
}

func (a *NEVMAddressUpdateEntry) Serialize(w io.Writer) error {
    err := wire.WriteVarBytes(w, 0, a.OldAddress)
    if err != nil {
        return err
    }
    err = wire.WriteVarBytes(w, 0, a.NewAddress)
    if err != nil {
        return err
    }
    err = binarySerializer.PutUint32(w, littleEndian, a.CollateralHeight)
    if err != nil {
        return err
    }
    return nil
}

func (a *NEVMRemoveEntry) Deserialize(r io.Reader) error {
    var err error
    a.Address, err = wire.ReadVarBytes(r, 0, HASH_SIZE, "Address")
    if err != nil {
        return err
    }
    return nil
}

func (a *NEVMRemoveEntry) Serialize(w io.Writer) error {
    err := wire.WriteVarBytes(w, 0, a.Address)
    if err != nil {
        return err
    }
    return nil
}

func (d *NEVMAddressDiff) Deserialize(r io.Reader) error {
    var err error

    // Deserialize AddedMNNEVM
    numAdded, err := wire.ReadVarInt(r, 0)
    if err != nil {
        return err
    }
    d.AddedMNNEVM = make([]NEVMAddressEntry, numAdded)
    for i := range d.AddedMNNEVM {
        err = d.AddedMNNEVM[i].Deserialize(r)
        if err != nil {
            return err
        }
    }

    // Deserialize UpdatedMNNEVM
    numUpdated, err := wire.ReadVarInt(r, 0)
    if err != nil {
        return err
    }
    d.UpdatedMNNEVM = make([]NEVMAddressUpdateEntry, numUpdated)
    for i := range d.UpdatedMNNEVM {
        err = d.UpdatedMNNEVM[i].Deserialize(r)
        if err != nil {
            return err
        }
    }

    // Deserialize RemovedMNNEVM
    numRemoved, err := wire.ReadVarInt(r, 0)
    if err != nil {
        return err
    }
    d.RemovedMNNEVM = make([]NEVMRemoveEntry, numRemoved)
    for i := range d.RemovedMNNEVM {
        err = d.RemovedMNNEVM[i].Deserialize(r)
        if err != nil {
            return err
        }
    }

    return nil
}

func (d *NEVMAddressDiff) Serialize(w io.Writer) error {
    var err error

    // Serialize AddedMNNEVM
    err = wire.WriteVarInt(w, 0, uint64(len(d.AddedMNNEVM)))
    if err != nil {
        return err
    }
    for i := range d.AddedMNNEVM {
        err = d.AddedMNNEVM[i].Serialize(w)
        if err != nil {
            return err
        }
    }

    // Serialize UpdatedMNNEVM
    err = wire.WriteVarInt(w, 0, uint64(len(d.UpdatedMNNEVM)))
    if err != nil {
        return err
    }
    for i := range d.UpdatedMNNEVM {
        err = d.UpdatedMNNEVM[i].Serialize(w)
        if err != nil {
            return err
        }
    }

    // Serialize RemovedMNNEVM
    err = wire.WriteVarInt(w, 0, uint64(len(d.RemovedMNNEVM)))
    if err != nil {
        return err
    }
    for i := range d.RemovedMNNEVM {
        err = d.RemovedMNNEVM[i].Serialize(w)
        if err != nil {
            return err
        }
    }

    return nil
}

func (a *NEVMBlockWire) Deserialize(r io.Reader) error {
    var err error

    // Deserialize NEVMBlockHash
    a.NEVMBlockHash = make([]byte, HASH_SIZE)
    _, err = io.ReadFull(r, a.NEVMBlockHash)
    if err != nil {
        return err
    }

    // Deserialize TxRoot
    a.TxRoot = make([]byte, HASH_SIZE)
    _, err = io.ReadFull(r, a.TxRoot)
    if err != nil {
        return err
    }

    // Deserialize ReceiptRoot
    a.ReceiptRoot = make([]byte, HASH_SIZE)
    _, err = io.ReadFull(r, a.ReceiptRoot)
    if err != nil {
        return err
    }

    // Deserialize NEVMBlockData
    a.NEVMBlockData, err = wire.ReadVarBytes(r, 0, MAX_NEVM_BLOCK_SIZE, "NEVMBlockData")
    if err != nil {
        return err
    }

    // Deserialize SYSBlockHash
    a.SYSBlockHash = make([]byte, HASH_SIZE)
    _, err = io.ReadFull(r, a.SYSBlockHash)
    if err != nil {
        return err
    }

    // Deserialize VersionHashes
    numVH, err := wire.ReadVarInt(r, 0)
    if err != nil {
        return err
    }
    a.VersionHashes = make([][]byte, numVH)
    for i := 0; i < int(numVH); i++ {
        a.VersionHashes[i], err = wire.ReadVarBytes(r, 0, HASH_SIZE, "VersionHash")
        if err != nil {
            return err
        }
    }

    // Deserialize Diff
    a.Diff = NEVMAddressDiff{}
    err = a.Diff.Deserialize(r)
    if err != nil {
        return err
    }

    return nil
}


func (a *NEVMBlockWire) Serialize(w io.Writer) error {
	var err error

	// NEVMBlockHash
	if _, err = w.Write(a.NEVMBlockHash[:]); err != nil {
		return err
	}

	// TxRoot
	if _, err = w.Write(a.TxRoot[:]); err != nil {
		return err
	}

	// ReceiptRoot
	if _, err = w.Write(a.ReceiptRoot[:]); err != nil {
		return err
	}

	// NEVMBlockData
	if err = wire.WriteVarBytes(w, 0, a.NEVMBlockData); err != nil {
		return err
	}

	// SYSBlockHash
	if _, err = w.Write(a.SYSBlockHash[:]); err != nil {
		return err
	}

	// VersionHashes (slice of byte slices)
	if err = wire.WriteVarInt(w, 0, uint64(len(a.VersionHashes))); err != nil {
		return err
	}
	for _, vh := range a.VersionHashes {
		if err = wire.WriteVarBytes(w, 0, vh); err != nil {
			return err
		}
	}

	// Diff
	if err = a.Diff.Serialize(w); err != nil {
		return err
	}

	return nil
}

func (a *NEVMDisconnectBlockWire) Deserialize(r io.Reader) error {
    var err error

    // Deserialize SYSBlockHash
    a.SYSBlockHash = make([]byte, HASH_SIZE)
    _, err = io.ReadFull(r, a.SYSBlockHash)
    if err != nil {
        return err
    }

    // Deserialize Diff
    a.Diff = NEVMAddressDiff{}
    err = a.Diff.Deserialize(r)
    if err != nil {
        return err
    }

    return nil
}
func (a *NEVMDisconnectBlockWire) Serialize(w io.Writer) error {
    var err error

    // Serialize SYSBlockHash
    _, err = w.Write(a.SYSBlockHash[:])
    if err != nil {
        return err
    }

    // Serialize Diff
    err = a.Diff.Serialize(w)
    if err != nil {
        return err
    }

    return nil
}

