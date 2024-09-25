package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// dont do this, see above edit
func Print(b []byte) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out.Bytes()))
}
