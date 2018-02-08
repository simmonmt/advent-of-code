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

	//fmt.Printf("decode called with %v %v\n", k, data)

	sum := 0
	switch k {
	case reflect.Map:
		for _, key := range v.MapKeys() {
			mapVal := v.MapIndex(key)
			sum += decode(mapVal.Interface())
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
