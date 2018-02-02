package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("missing salt")
	}

	salt := os.Args[1]

	for i := 0; ; i++ {
		if i%1000 == 0 {
			fmt.Printf("%d...\n", i)
		}

		data := []byte(salt + strconv.Itoa(i))
		hash := md5.Sum([]byte(data))

		if hash[0] == 0 && hash[1] == 0 && (hash[2]&0xf0) == 0 {
			fmt.Println(i)
			for _, b := range hash {
				fmt.Printf("%02x ", b)
			}
			fmt.Printf("\n")

			break
		}
	}
}
