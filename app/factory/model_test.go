package factory

import (
	"encoding/base64"
	"testing"
)

// TestMacToMagicBytesReference verifies that computing the payload for the
// originally captured MAC reproduces the captured base64 blob exactly.
func TestMacToMagicBytesReference(t *testing.T) {
	const original = "AAAAAGAIAACTBwAAOggAALoAAACQBwAAxAcAAMoGAACVBAAATggAAM0BAAAnCA=="
	got := MacToMagicBytes(magicPayloadRefMAC)
	if base64.StdEncoding.EncodeToString(got) != original {
		t.Fatalf("reference MAC did not reproduce original payload:\n got=%s\nwant=%s",
			base64.StdEncoding.EncodeToString(got), original)
	}
}

// TestMacToMagicBytesDeterministic verifies the output is stable for a given MAC.
func TestMacToMagicBytesDeterministic(t *testing.T) {
	mac := [6]byte{0xde, 0xad, 0xbe, 0xef, 0x00, 0x01}
	a := MacToMagicBytes(mac)
	b := MacToMagicBytes(mac)
	if string(a) != string(b) {
		t.Fatal("MacToMagicBytes is not deterministic for a fixed MAC")
	}
	if len(a) != 46 {
		t.Fatalf("unexpected payload length: %d, want 46", len(a))
	}
}

// TestMacToMagicBytesOtherMACs verifies payload generation for several other
// MACs, so the derivation (frame XOR mac) can be checked against a real device.
func TestMacToMagicBytesOtherMACs(t *testing.T) {
	cases := []struct {
		mac [6]byte
		want string
	}{
		{[6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, "AAAAAGAIAACUAAAAPQ8AAJMpAAC5LgAAkVIAAJ9TAACgMQAAez0AAJpWAABwXw=="},
		{[6]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, "//8AAJ/3AABr/wAAwvAAAGzWAABG0QAAbq0AAGCsAABfzgAAhMIAAGWpAACPoA=="},
		{[6]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}, "EREAAHEZAAC2IgAAHy0AAKAaAACKHQAA1RYAANsXAAD1ZAAALmgAAPwwAAAWOQ=="},
		{[6]byte{0xde, 0xad, 0xbe, 0xef, 0x00, 0x01}, "3t4AAL7WAAA5rQAAkKIAAC2XAAAHkAAAfr0AAHC8AACgMQAAez0AAJtXAABxXg=="},
		{[6]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}, "qqoAAMqiAAAvuwAAhrQAAF/lAAB14gAATI8AAEKOAABO3wAAldMAAGWpAACPoA=="},
	}
	for _, tc := range cases {
		got := base64.StdEncoding.EncodeToString(MacToMagicBytes(tc.mac))
		if got != tc.want {
			t.Errorf("MacToMagicBytes(%x) = %s, want %s", tc.mac, got, tc.want)
		}
	}
}
