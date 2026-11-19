package helper

import "testing"

func TestPasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     string
	}{
		{"terlalu pendek", "a1b2c3", "minimal 8 karakter"},
		{"tanpa angka", "passwordku", "harus memuat huruf dan angka"},
		{"tanpa huruf", "12345678", "harus memuat huruf dan angka"},
		{"terlalu umum", "password123", "password terlalu umum"},
		{"valid", "rahasia123", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := passwordStrength(tt.password); got != tt.want {
				t.Fatalf("passwordStrength(%q) = %q, ingin %q", tt.password, got, tt.want)
			}
		})
	}
}
