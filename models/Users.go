package models

import (
	"errors"
	"github.com/runeasymail/ManagementAPI/helpers"
	"os/exec"
	"strings"
)

type Users struct {
	Id       uint64 `db:"id" json:"id" form:"id"`
	DomainID uint64 `db:"domain_id" json:"domain_id" form:"domain_id" validation:"required"`
	Password string `db:"password" json:"-" form:"password" validation:"required"`
	Email    string `db:"email" json:"email" form:"email" validation:"email,required"`
}

func GetAllUsers(domain_id uint64) (result []Users) {

	sql := `select * from virtual_users where domain_id = ? order by id DESC`
	helpers.MyDB.Unsafe().Select(&result, sql, domain_id)
	return
}

func ChangePassword(data Users) {
	sql := `update virtual_users set password = ? where id = ? limit 1`
	helpers.MyDB.Unsafe().Exec(sql, data.GenEncryptedPassword(), data.Id)
}

func AddNewUser(data Users) (result bool, err error) {

	var count uint64
	sql := `select count(id) from virtual_users where email = ?`
	helpers.MyDB.Unsafe().Get(&count, sql, data.Email)

	if count != 0 {
		err = errors.New("User already exist!")
		return
	}

	sql = `insert into virtual_users (domain_id,password,email) values(?,?,?)`
	helpers.MyDB.Unsafe().Exec(sql, data.DomainID, data.GenEncryptedPassword(), data.Email)

	result = true

	return
}

func (data Users) GenEncryptedPassword() string {
	cmd := exec.Command("openssl", "passwd", "-1", data.Password)
	output, _ := cmd.Output()

	out := string(output)

	out = strings.Replace(out, "\n", "", -1)

	return out
}
