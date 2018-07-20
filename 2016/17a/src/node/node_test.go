package node

import (
	"reflect"
	"testing"
)

func TestSerialize(t *testing.T) {
	type TestCase struct {
		ser   string
		deser *Node
	}

	testCases := []TestCase{
		TestCase{"0,0,pass,", &Node{0, 0, "pass", ""}},
		TestCase{"1,2,pass,UD", &Node{1, 2, "pass", "UD"}},
		TestCase{"-1,-2,pass,UD", &Node{-1, -2, "pass", "UD"}},
	}

	for _, tc := range testCases {
		t.Run(tc.ser, func(t *testing.T) {
			if deser, err := Deserialize(tc.ser); err == nil && !reflect.DeepEqual(deser, tc.deser) {
				t.Errorf(`Deserialize("%v") = %+v, %v, want %+v, nil`, tc.ser, deser, err, tc.deser)
			}

			if ser := tc.deser.Serialize(); ser != tc.ser {
				t.Errorf(`%+v Serialize() = "%v", want "%v"`, tc.deser, ser, tc.ser)
			}
		})
	}
}
