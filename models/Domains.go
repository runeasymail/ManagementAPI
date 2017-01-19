package models

import (
	"errors"
	"github.com/runeasymail/ManagementAPI/helpers"
)

type Domains struct {
	Id         uint64 `db:"id" json:"id" form:"id"`
	Name       string `db:"name" json:"name" form:"name" validation:"required"`
	UsersCount uint64 `db:"users_count" json:"users_count"`
}

// Get All domains orderder by id DESC
// with users count
func GetDomains() (result []Domains) {
	sql := `SELECT 
				virtual_domains.id,
				virtual_domains.name,
				(
				SELECT COUNT(id)
				FROM virtual_users
				WHERE virtual_users.domain_id = virtual_domains.id) AS users_count
			FROM 
			virtual_domains
			ORDER BY id DESC`
	helpers.MyDB.Unsafe().Select(&result, sql)
	return
}

func AddNewDomain(domain string, username string, password string) (result bool, err error) {

	var count_domains uint64
	sql := `select count(id) from virtual_domains where name = ?`
	helpers.MyDB.Unsafe().Get(&count_domains, sql, domain)

	if count_domains != 0 {
		err = errors.New("Domain is already exist")
		return
	}

	sql = `INSERT INTO virtual_domains(name) values(?)`
	res, err := helpers.MyDB.Unsafe().Exec(sql, domain)

	if err != nil {
		return
	}

	id, _ := res.LastInsertId()

	// add new User
	userData := Users{DomainID: uint64(id), Password: password, Email: username}
	_, err = AddNewUser(userData)

	if err != nil {
		return
	}

	result = true
	return
}
