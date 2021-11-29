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

package object

import "testing"

func TestSerializeGenerator(t *testing.T) {
	o := Generator(2)
	oSer := o.Serialize()
	if oSer != 'B' {
		t.Errorf("Generator(2).Serialize() = %v, want %v", string(oSer), "B")
	}
	oDeser := Deserialize(oSer)
	if oDeser != o {
		t.Errorf("Deserialize(Generator(2).Serialize()) = %v, want %v", oDeser, o)
	}
}

func TestSerializeMicrochip(t *testing.T) {
	o := Microchip(2)
	oSer := o.Serialize()
	if oSer != 'b' {
		t.Errorf("Microchip(2).Serialize() = %v, want %v", string(oSer), "b")
	}
	oDeser := Deserialize(oSer)
	if oDeser != o {
		t.Errorf("Deserialize(Microchip(2).Serialize()) = %v, want %v", oDeser, o)
	}
}
