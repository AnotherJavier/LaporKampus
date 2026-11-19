// Package tools menampung alat bantu pengembangan, bukan bagian dari layanan.
//
// Setiap alat berada di subpaket tersendiri dan dijalankan dengan `go run`:
//
//	go run ./tools/loadtest  -url http://localhost:3000/api/v1/health
//	go run ./tools/redeliver -queue default
package tools
