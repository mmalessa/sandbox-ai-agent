package appdebug

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(p interface{}) {
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(string(b))
}
