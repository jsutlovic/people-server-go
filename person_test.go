package main

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq/hstore"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPersonFields(t *testing.T) {
	personType := reflect.TypeOf(Person{})

	fieldCount := personType.NumField()
	assert.Equal(t, fieldCount, 5)

	_, idExists := personType.FieldByName("Id")
	_, userIdExists := personType.FieldByName("UserId")
	_, nameExists := personType.FieldByName("Name")
	_, metaExists := personType.FieldByName("Meta")
	_, colorExists := personType.FieldByName("Color")

	assert.True(t, idExists)
	assert.True(t, userIdExists)
	assert.True(t, nameExists)
	assert.True(t, metaExists)
	assert.True(t, colorExists)
}

func TestPersonFieldsJson(t *testing.T) {
	personType := reflect.TypeOf(Person{})

	idField, _ := personType.FieldByName("Id")
	userIdField, _ := personType.FieldByName("UserId")
	nameField, _ := personType.FieldByName("Name")
	metaField, _ := personType.FieldByName("Meta")
	colorField, _ := personType.FieldByName("Color")

	assert.Equal(t, idField.Tag.Get("json"), "id")
	assert.Equal(t, userIdField.Tag.Get("json"), "user_id")
	assert.Equal(t, nameField.Tag.Get("json"), "name")
	assert.Equal(t, metaField.Tag.Get("json"), "meta")
	assert.Equal(t, colorField.Tag.Get("json"), "color")
}

func TestPersonFieldsDb(t *testing.T) {
	personType := reflect.TypeOf(Person{})

	idField, _ := personType.FieldByName("Id")
	userIdField, _ := personType.FieldByName("UserId")
	nameField, _ := personType.FieldByName("Name")
	metaField, _ := personType.FieldByName("Meta")
	colorField, _ := personType.FieldByName("Color")

	assert.Equal(t, idField.Tag.Get("db"), "")
	assert.Equal(t, userIdField.Tag.Get("db"), "user_id")
	assert.Equal(t, nameField.Tag.Get("db"), "")
	assert.Equal(t, metaField.Tag.Get("db"), "")
	assert.Equal(t, colorField.Tag.Get("db"), "")
}

func TestGetPersonNotFound(t *testing.T) {
	pgdbs := NewPgDbService("mock", "")

	personId := 1
	userId := 2

	sqlmock.ExpectQuery(`SELECT \* FROM "person" WHERE id=\? AND user_id=\?`).
		WithArgs(personId, userId).
		WillReturnError(errors.New("Could not find person"))

	p, err := pgdbs.GetPerson(userId, personId)
	if !assert.Nil(t, p, "Person should be nil") {
		return
	}

	if !assert.NotNil(t, err, "Should have an error") {
		return
	}

	assert.Equal(t, err.Error(), "Could not find person")
}

func TestGetPerson(t *testing.T) {
	pgdbs := NewPgDbService("mock", "")

	personId := 1
	userId := 2
	cols := []string{"id", "user_id", "name", "meta", "color"}

	meta := new(hstore.Hstore)
	meta.Scan([]byte(`"type"=>"asdf"`))

	color := new(sql.NullInt64)
	color.Scan(1)

	sqlmock.ExpectQuery(`SELECT \* FROM "person" WHERE id=\? AND user_id=\?`).
		WithArgs(personId, userId).
		WillReturnRows(sqlmock.NewRows(cols).AddRow(1, userId, "Person 1", meta, color))

	p, err := pgdbs.GetPerson(userId, personId)
	if !assert.Nil(t, err, "Query should not error") {
		return
	}

	if !assert.NotNil(t, p, "Person should not be nil") {
		return
	}

	assert.Equal(t, p.Id, personId)
}
