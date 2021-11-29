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
