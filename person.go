package main

import (
	"database/sql"
	"github.com/lib/pq/hstore"
)

type Person struct {
	Id     int           `json:"id"`
	UserId int           `db:"user_id" json:"user_id"`
	Name   string        `json:"name"`
	Meta   hstore.Hstore `json:"meta"`
	Color  sql.NullInt64 `json:"color"`
}

func (s *pgDbService) GetPerson(userId, id int) (*Person, error) {
	person := new(Person)

	err := s.db.Get(person, s.db.Rebind(`SELECT * FROM "person" WHERE id=? AND user_id=?`), id, userId)
	if err != nil {
		return nil, err
	}
	return person, nil
}
