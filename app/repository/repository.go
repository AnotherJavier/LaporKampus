// Package repository berisi kontrak akses data dan implementasinya.
//
// Gaya modul: kontrak dan implementasi berada di paket yang sama
// (repository.go + <entitas>_repository.go), sehingga service cukup
// bergantung pada interface di sini dan dapat diuji dengan fake.
package repository