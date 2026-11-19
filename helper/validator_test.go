package helper

import (
	"testing"

	"github.com/vyerru/lapor-kampus/app/model"
)

func TestValidateStruct(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		req := model.RegisterRequest{
			Username: "budi.santoso",
			Email:    "budi@example.com",
			Password: "rahasia123",
		}
		if fields := ValidateStruct(req); fields != nil {
			t.Fatalf("seharusnya valid, dapat %v", fields)
		}
	})

	t.Run("field tidak valid", func(t *testing.T) {
		req := model.RegisterRequest{
			Username: "budi santoso",
			Email:    "bukan-email",
			Password: "123",
		}
		fields := ValidateStruct(req)
		for _, key := range []string{"username", "email", "password"} {
			if fields[key] == "" {
				t.Fatalf("field %q seharusnya punya pesan, dapat %v", key, fields)
			}
		}
	})
}
