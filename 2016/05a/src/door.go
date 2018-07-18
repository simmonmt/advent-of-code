package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %v door_id", os.Args[0])
	}

	doorID := os.Args[1]

	pw := []byte{}
	for i := 0; ; i++ {
		if i != 0 && i%1000000 == 0 {
			fmt.Printf("i=%v\n", i)
		}

		data := fmt.Sprintf("%s%d", doorID, i)
		hash := md5.Sum([]byte(data))

		if !(hash[0] == 0 && hash[1] == 0 && (hash[2]&0xf0) == 0) {
			continue
		}

		pw = append(pw, hash[2]&0x0f)
		fmt.Printf("pwlen now %v, pw is %v\n", len(pw), pw)
		if len(pw) == 8 {
			break
		}
	}

	for _, b := range pw {
		fmt.Printf("%x", b)
	}
	fmt.Println()
}
