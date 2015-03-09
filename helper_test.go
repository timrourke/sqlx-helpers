package helper

import (
	"strings"
	"testing"
)

type DeclaredTag struct {
	Id   string `db:"diff_id"`
	Name string `db:"diff_name"`
}

type NoTag struct {
	Id   string
	Name string
}

func TestDeclaredTag(t *testing.T) {
	model := DeclaredTag{}

	data := map[string]interface{}{
		"Id":   []string{"one", "two"},
		"Name": "A name",
	}

	where, expanded := CreateWhere(model, data)

	if !strings.ContainsAny(where, "diff_id in (:Id0,:Id1)") {
		t.Error("Incorrect where: " + where)
	}
	if !strings.ContainsAny(where, "diff_name = :Name") {
		t.Error("Incorrect where: " + where)
	}

	if _, ok := expanded["Id0"]; !ok {
		t.Error("Incorrect expanded")
	}
	if _, ok := expanded["Id1"]; !ok {
		t.Error("Incorrect expanded")
	}
}

func TestNoTags(t *testing.T) {
	model := NoTag{}

	data := map[string]interface{}{
		"Id":   []string{"one", "two"},
		"Name": "A name",
	}

	where, expanded := CreateWhere(model, data)
	if !strings.ContainsAny(where, "Id in (:Id0,:Id1)") {
		t.Error("Incorrect where: " + where)
	}
	if !strings.ContainsAny(where, "Name = :Name") {
		t.Error("Incorrect where: " + where)
	}

	if _, ok := expanded["Id0"]; !ok {
		t.Error("Incorrect expanded")
	}
	if _, ok := expanded["Id1"]; !ok {
		t.Error("Incorrect expanded")
	}
}
