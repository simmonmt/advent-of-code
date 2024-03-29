// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
