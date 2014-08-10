package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/lib/pq/hstore"
)

type Person struct {
	Id     int
	UserId int `db:"user_id"`
	Name   string
	Meta   hstore.Hstore
	Color  sql.NullInt64
}

type personJSON struct {
	Id     int               `json:"id"`
	UserId int               `json:"user_id"`
	Name   string            `json:"name"`
	Meta   map[string]string `json:"meta"`
	Color  json.RawMessage   `json:"color"`
}

func (p *Person) MarshalJSON() ([]byte, error) {
	colorVal, _ := p.Color.Value()
	colorJSON := []byte(Jsonify(colorVal))

	metaVal := make(map[string]string)
	if p.Meta.Map != nil {
		for key, val := range p.Meta.Map {
			if val.Valid {
				metaVal[key] = val.String
			}
		}
	}

	pJson := personJSON{
		Id:     p.Id,
		UserId: p.UserId,
		Name:   p.Name,
		Meta:   metaVal,
		Color:  colorJSON,
	}

	return json.Marshal(&pJson)
}

func (p *Person) UnmarshalJSON(b []byte) error {
	pJson := new(personJSON)

	err := json.Unmarshal(b, pJson)
	if err != nil {
		return err
	}

	p.Id = pJson.Id
	p.UserId = pJson.UserId
	p.Name = pJson.Name

	p.Meta.Map = make(map[string]sql.NullString)

	for key, val := range pJson.Meta {
		p.Meta.Map[key] = sql.NullString{val, true}
	}

	if bytes.Equal(pJson.Color, []byte("null")) {
		p.Color = sql.NullInt64{0, false}
	} else {
		var colorVal int64
		json.Unmarshal(pJson.Color, &colorVal)
		p.Color = sql.NullInt64{colorVal, true}
	}

	return nil
}

func (s *pgDbService) GetPerson(userId, id int) (*Person, error) {
	person := new(Person)

	err := s.db.Get(person, s.db.Rebind(`SELECT * FROM "person" WHERE id=? AND user_id=?`), id, userId)
	if err != nil {
		return nil, err
	}
	return person, nil
}

func (s *pgDbService) GetPeople(userId int) ([]Person, error) {
	people := []Person{}

	err := s.db.Select(&people, s.db.Rebind(`SELECT * FROM "person" WHERE user_id=?`), userId)
	if err != nil {
		return nil, err
	}

	return people, nil
}
