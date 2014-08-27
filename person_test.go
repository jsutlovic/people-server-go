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

func TestGetPeopleError(t *testing.T) {
	pgdbs := NewPgDbService("mock", "")

	userId := 1

	sqlmock.ExpectQuery(`SELECT \* FROM "person" WHERE user_id=\?`).
		WithArgs(userId).
		WillReturnError(errors.New("Could not find person (list)"))

	pp, err := pgdbs.GetPeople(userId)
	if !assert.Nil(t, pp, "Person list should be null") {
		return
	}

	if !assert.NotNil(t, err, "Should have an error") {
		return
	}

	assert.Equal(t, err.Error(), "Could not find person (list)")
}

func TestGetPeople(t *testing.T) {
	pgdbs := NewPgDbService("mock", "")

	userId := 1
	personId1 := 2
	personId2 := 3

	cols := []string{"id", "user_id", "name", "meta", "color"}

	name1 := "Person 1"
	name2 := "Person 2"

	meta1 := hstore.Hstore{nil}
	meta2 := hstore.Hstore{map[string]sql.NullString{"type": sql.NullString{"asdf", true}}}

	metaVal1, _ := meta1.Value()
	metaVal2, _ := meta2.Value()

	color1 := sql.NullInt64{0, false}
	color2 := sql.NullInt64{1, true}

	colorVal1, _ := color1.Value()
	colorVal2, _ := color2.Value()

	sqlmock.ExpectQuery(`SELECT \* FROM "person" WHERE user_id=\?`).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows(cols).
		AddRow(personId1, userId, name1, metaVal1, colorVal1).
		AddRow(personId2, userId, name2, metaVal2, colorVal2))

	pp, err := pgdbs.GetPeople(userId)
	if !assert.Nil(t, err, "Query should not error") {
		return
	}

	if !assert.NotNil(t, pp, "People result should not be nil") {
		return
	}

	if !assert.Len(t, pp, 2, "People result should have a length of 2") {
		return
	}

	p1 := pp[0]
	p2 := pp[1]

	assert.Equal(t, p1.Id, personId1)
	assert.Equal(t, p2.Id, personId2)

	assert.Equal(t, p1.UserId, userId)
	assert.Equal(t, p2.UserId, userId)

	assert.Equal(t, p1.Name, name1)
	assert.Equal(t, p2.Name, name2)

	assert.Equal(t, p1.Meta, meta1)
	assert.Equal(t, p2.Meta, meta2)

	assert.Equal(t, p1.Color, color1)
	assert.Equal(t, p2.Color, color2)
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
	personJSONAnomaloustests := []struct {
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

func TestPersonUnmarshalAnomalousJSON(t *testing.T) {
	var personUnmarshalAnomalousTests = []struct {
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
					"other": sql.NullString{"", true},
				}},
				Color: sql.NullInt64{1, true},
			},
			json: `{"id":1,"user_id":1,"name":"Test 1","meta":{"type":"asdf", "other": null},"color":1}`,
		},
		{
			p: Person{
				Id:     2,
				UserId: 2,
				Name:   "Test 2",
				Meta:   hstore.Hstore{map[string]sql.NullString{}},
				Color:  sql.NullInt64{0, false},
			},
			json: `{"id":2,"user_id":2,"name":"Test 2","meta":null,"color":null}`,
		},
	}

	for _, test := range personUnmarshalAnomalousTests {
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
		`"id": 1`,
	}

	for _, test := range personUnmarshalErrorTests {
		unmarshaled := Person{}
		err := unmarshaled.UnmarshalJSON([]byte(test))
		assert.NotNil(t, err)
	}
}

func TestCreatePersonEmptyName(t *testing.T) {
	pgdbs := NewPgDbService("mock", "")

	userId := 1
	meta := map[string]string{}
	color := sql.NullInt64{1, true}

	names := []string{"", " ", "\t", "\n"}

	for _, name := range names {
		_, err := pgdbs.CreatePerson(userId, name, meta, color)
		assert.Error(t, err, "Empty name should cause error")
		assert.Equal(t, err.Error(), "Person name cannot be empty")
	}
}
