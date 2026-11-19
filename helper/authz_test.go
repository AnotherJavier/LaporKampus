package helper

import "testing"

func TestPermissionSet(t *testing.T) {
	perms := NewPermissionSet(map[string][]string{
		"admin":      {"ticket.read", "ticket.assign"},
		"supervisor": {"ticket.read"},
	})

	if !perms.Can("admin", "ticket.assign") {
		t.Fatal("admin seharusnya boleh ticket.assign")
	}
	if perms.Can("supervisor", "ticket.assign") {
		t.Fatal("supervisor tidak boleh ticket.assign")
	}
	if perms.Can("mahasiswa", "ticket.read") {
		t.Fatal("role tak dikenal seharusnya ditolak")
	}

	var nilSet *PermissionSet
	if nilSet.Can("admin", "ticket.read") {
		t.Fatal("PermissionSet nil harus menolak semua")
	}
}
