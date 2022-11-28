package helpers

import (
	"fmt"
	"os"
)

func Must(err error) {
	if err == nil {
		return
	}
	fmt.Printf("---\nERROR! Dependencies have not been resolved!\n---\n")
	fmt.Printf("%+v\n", err)
	os.Exit(1)
}
