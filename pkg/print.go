package pkg

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ingshtrom/dh-debug/pkg/types"
)

func PrintDebugTests(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("error reading debug file: %v", err)
		os.Exit(1)
	}

	var dcs *[]types.DebugCommand

	err = json.Unmarshal(data, &dcs)
	if err != nil {
		fmt.Printf("error parsing debug file: %v", err)
		os.Exit(1)
	}

	for _, dc := range *dcs {
		dc.Print()
	}
}
