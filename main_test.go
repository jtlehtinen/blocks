package main

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestDoubleSha(t *testing.T) {
	want, _ := hex.DecodeString("9595c9df90075148eb06860365df33584b75bff782a510c6cd4883a419833d50")
	got := DoubleSha256([]byte("hello"))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\n\tgot  %q\n\twant %q", hex.EncodeToString(got), hex.EncodeToString(want))
	}
}
