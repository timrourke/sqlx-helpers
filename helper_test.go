package helper

import (
	"fmt"
	"strings"
	"testing"
)

type DeclaredTag struct {
	Id     string `db:"diff_id"`
	Name   string `db:"diff_name"`
	Ignore string `db:"-"`
}

type NoTag struct {
	Id   string
	Name string
}

func TestWhereDeclaredTag(t *testing.T) {
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

func TestWhereNoTags(t *testing.T) {
	model := NoTag{}

	data := map[string]interface{}{
		"Id":   []string{"one", "two"},
		"Name": "A name",
	}

	where, expanded := CreateWhere(model, data)
	if !strings.ContainsAny(where, "id in (:Id0,:Id1)") {
		t.Error("Incorrect where: " + where)
	}
	if !strings.ContainsAny(where, "name = :Name") {
		t.Error("Incorrect where: " + where)
	}

	if _, ok := expanded["Id0"]; !ok {
		t.Error("Incorrect expanded")
	}
	if _, ok := expanded["Id1"]; !ok {
		t.Error("Incorrect expanded")
	}
}

func TestInsertDeclaredTags(t *testing.T) {
	model := DeclaredTag{}

	insert := CreateInsert(model)

	fmt.Println(insert)

	if strings.Contains(insert, ":-") {
		t.Error("Incorrect insert: " + insert)
	}

	if !strings.Contains(insert, "(diff_id, diff_name) values (:diff_id, :diff_name)") {
		t.Error("Incorrect insert: " + insert)
	}
}

func TestInsertNoTags(t *testing.T) {
	model := NoTag{}

	insert := CreateInsert(model)

	if !strings.Contains(insert, "(id, name) values (:Id, :Name)") {
		t.Error("Incorrect insert: " + insert)
	}
}

func TestUpdateDeclaredTags(t *testing.T) {
	model := DeclaredTag{}

	update := CreateUpdate(model)

	if !strings.Contains(update, "diff_id=:diff_id, diff_name=:diff_name") {
		t.Error("Incorrect update: " + update)
	}
}

func TestUpdateNoTags(t *testing.T) {
	model := NoTag{}

	update := CreateUpdate(model)

	if !strings.Contains(update, "id=:Id, name=:Name") {
		t.Error("Incorrect update: " + update)
	}
}

func TestExclude(t *testing.T) {
	model := DeclaredTag{}

	update := CreateUpdate(model, "diff_id")

	if strings.Contains(update, "diff_id") {
		t.Error("Incorrect update: " + update)
	}

	insert := CreateInsert(model, "diff_name")

	if strings.Contains(insert, "diff_name") {
		t.Error("Incorrect insert: " + insert)
	}
}
