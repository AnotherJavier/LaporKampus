package helper

import (
	"testing"
	"time"
)

func TestCursorRoundTrip(t *testing.T) {
	createdAt := time.Date(2026, 11, 18, 9, 30, 0, 0, time.UTC)

	encoded := EncodeCursor(createdAt, 42)
	decoded, err := DecodeCursor(encoded)
	if err != nil {
		t.Fatalf("DecodeCursor gagal: %v", err)
	}

	if !decoded.CreatedAt.Equal(createdAt) {
		t.Fatalf("created_at = %v, ingin %v", decoded.CreatedAt, createdAt)
	}
	if decoded.ID != 42 {
		t.Fatalf("id = %d, ingin 42", decoded.ID)
	}
}

func TestDecodeCursorInvalid(t *testing.T) {
	cases := []string{
		"",
		"  ",
		"bukan-base64!",
		"YWJj", // "abc" — tidak memuat pemisah
		"MTIz", // "123" — tidak memuat pemisah
	}

	for _, input := range cases {
		if _, err := DecodeCursor(input); err == nil {
			t.Fatalf("DecodeCursor(%q) seharusnya gagal", input)
		}
	}
}
