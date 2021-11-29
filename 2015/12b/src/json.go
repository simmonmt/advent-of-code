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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

func decode(data interface{}) int {
	k := reflect.TypeOf(data).Kind()
	v := reflect.ValueOf(data)

	sum := 0
	switch k {
	case reflect.Map:
		foundRed := false
		for _, key := range v.MapKeys() {
			mapVal := v.MapIndex(key).Interface()

			if reflect.TypeOf(mapVal).Kind() == reflect.String &&
				reflect.ValueOf(mapVal).String() == "red" {
				foundRed = true
			} else {
				sum += decode(mapVal)
			}
		}
		if foundRed {
			sum = 0
		}
		break

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			sum += decode(elem.Interface())
		}
		break

	case reflect.Float64:
		sum = int(v.Float())
		break

	case reflect.String:
		break

	default:
		panic(fmt.Sprintf("unknown kind %v", k))
	}
	return sum
}

func main() {
	var data interface{}
	err := json.NewDecoder(os.Stdin).Decode(&data)
	if err != nil {
		log.Fatalf("failed to parse input: %v", err)
	}

	sum := decode(data)
	fmt.Println(sum)
}
