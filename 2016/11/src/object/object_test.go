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
