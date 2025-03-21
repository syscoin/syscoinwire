# Syscoin Wire Package

This package provides a comprehensive set of tools for serialization and deserialization of Syscoin blockchain-specific transactions, NEVM blocks, and related structures. It integrates smoothly with existing Bitcoin-based infrastructure such as `btcd`, extending functionality specifically for Syscoin.

## Features

- Serialization and deserialization of Syscoin asset allocations
- Handling of NEVM-specific block structures
- Efficient binary serialization optimized for blockchain data
- Comprehensive unit tests covering edge cases

## Installation

Ensure your project uses Go Modules, then install the package:

```sh
go get github.com/syscoin/syscoinwire/syscoin/wire
```

## Usage

Here's how to include and use Syscoin Wire in your Go project:

```go
package main

import (
	"fmt"
	"github.com/syscoin/syscoinwire/syscoin/wire"
)

func main() {
	assetAllocation := wire.AssetAllocationType{
		VoutAssets: []wire.AssetOutType{
			{
				AssetGuid: 1234,
				Values: []wire.AssetOutValueType{
					{N: 0, ValueSat: 5000},
				},
			},
		},
	}

	fmt.Printf("Asset Allocation: %+v\n", assetAllocation)
}
```

## Running Tests

To run unit tests provided by the package, navigate to the root of your project and run:

```sh
go test ./syscoin/wire
```

Tests include serialization and deserialization validation ensuring integrity and consistency:

- `AssetAllocationType`
- `NEVMBlockWire`
- `NEVMDisconnectBlockWire`

These tests thoroughly cover typical scenarios and edge cases such as empty data, maximum sizes, and boundary numeric values.

## Project Structure

Recommended project structure:

```
syscoinwire/
├── go.mod
├── syscoin
│   └── wire
│       ├── asset.go
│       ├── asset_test.go
│       ├── nevmblock.go
│       ├── nevmblock_test.go
│       └── (other future .go files)
└── (additional future project files)
```

## Contributing

Contributions and improvements are welcome! Feel free to open issues or pull requests on GitHub.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

