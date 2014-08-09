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

	meta := hstore.Hstore{map[string]sql.NullString{"type": sql.NullString{"asdf", true}}}
	metaVal, _ := meta.Value()

	color := sql.NullInt64{1, true}
	colorVal, _ := color.Value()

	sqlmock.ExpectQuery(`SELECT \* FROM "person" WHERE id=\? AND user_id=\?`).
		WithArgs(personId, userId).
		WillReturnRows(sqlmock.NewRows(cols).AddRow(1, userId, "Person 1", metaVal, colorVal))

	p, err := pgdbs.GetPerson(userId, personId)
	if !assert.Nil(t, err, "Query should not error") {
		return
	}

	if !assert.NotNil(t, p, "Person should not be nil") {
		return
	}

	assert.Equal(t, p.Id, personId)
	assert.Equal(t, p.Meta, meta)
	assert.Equal(t, p.Color, color)
}

var personJSONtests = []struct {
	p    Person
	json string
}{
	{
		p: Person{
			Id:     1,
			UserId: 1,
			Name:   "Test 1",
			Meta: hstore.Hstore{map[string]sql.NullString{
				"type": sql.NullString{"asdf", true},
			}},
			Color: sql.NullInt64{1, true},
		},
		json: `{"id":1,"user_id":1,"name":"Test 1","meta":{"type":"asdf"},"color":1}`,
	},
	{
		p: Person{
			Id:     2,
			UserId: 2,
			Name:   "Test 2",
			Meta:   hstore.Hstore{map[string]sql.NullString{}},
			Color:  sql.NullInt64{0, false},
		},
		json: `{"id":2,"user_id":2,"name":"Test 2","meta":{},"color":null}`,
	},
}

var personJSONAnomaloustests = []struct {
	p    Person
	json string
}{
	{
		p: Person{
			Id:     1,
			UserId: 1,
			Name:   "Test 1",
			Meta: hstore.Hstore{map[string]sql.NullString{
				"type":  sql.NullString{"asdf", true},
				"other": sql.NullString{"", false},
			}},
			Color: sql.NullInt64{1, true},
		},
		json: `{"id":1,"user_id":1,"name":"Test 1","meta":{"type":"asdf"},"color":1}`,
	},
	{
		p: Person{
			Id:     2,
			UserId: 2,
			Name:   "Test 2",
			Meta:   hstore.Hstore{nil},
			Color:  sql.NullInt64{0, false},
		},
		json: `{"id":2,"user_id":2,"name":"Test 2","meta":{},"color":null}`,
	},
}

func TestPersonMarshalJSON(t *testing.T) {
	for _, test := range personJSONtests {
		marshaled, err := test.p.MarshalJSON()
		if !assert.Nil(t, err) {
			break
		}
		assert.Equal(t, string(marshaled), test.json)
	}
}

func TestPersonMarshalAnomalousJSON(t *testing.T) {
	for _, test := range personJSONAnomaloustests {
		marshaled, err := test.p.MarshalJSON()
		if !assert.Nil(t, err) {
			break
		}
		assert.Equal(t, string(marshaled), test.json)
	}
}

func TestPersonUnmarshalJSON(t *testing.T) {
	for _, test := range personJSONtests {
		unmarshaled := Person{}
		err := unmarshaled.UnmarshalJSON([]byte(test.json))
		if !assert.Nil(t, err) {
			break
		}
		assert.Equal(t, unmarshaled, test.p)
	}
}

func TestPersonUnmarshalError(t *testing.T) {
	personUnmarshalErrorTests := []string{
		"",
		"{id: 1}",
	}

	for _, test := range personUnmarshalErrorTests {
		unmarshaled := Person{}
		err := unmarshaled.UnmarshalJSON([]byte(test))
		assert.NotNil(t, err)
	}
}
