package main

import (
	"fmt"

	"github.com/gorilla/securecookie"
)

var keyLengths = [...]int{16, 24, 32, 64}

func main() {
	fmt.Println("----- Keys -----")
	for _, keyLength := range keyLengths {
		key := securecookie.GenerateRandomKey(keyLength)

		fmt.Printf("%d bytes: ", keyLength)
		for _, ch := range key {
			hex := fmt.Sprintf("%x", ch)
			if len(hex) < 2 {
				hex = "0" + hex
			}

			fmt.Printf("\\x%s", hex)
		}

		fmt.Println("\n----------------")
	}
}
