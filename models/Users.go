package models

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/runeasymail/ManagementAPI/helpers"
	"log"
	"os"
)

type Users struct {
	Id                  uint64 `db:"id" json:"id" form:"id"`
	DomainID            uint64 `db:"domain_id" json:"domain_id" form:"domain_id" validation:"required"`
	Password            string `db:"password" json:"-" form:"password" validation:"required"`
	PasswordIsEncrypted bool   `form:is_encrypted`
	Email               string `db:"email" json:"email" form:"email" validation:"email,required"`
}

func GetAllUsers(domain_id uint64) (result []Users) {

	sql := `select * from virtual_users where domain_id = ? order by id DESC`
	helpers.MyDB.Unsafe().Select(&result, sql, domain_id)
	return
}

func ChangePassword(data Users) {
	sql := `update virtual_users set password = ? where id = ? and domain_id = ?  limit 1`
	helpers.MyDB.Unsafe().Exec(sql, data.GenEncryptedPassword(), data.Id, data.DomainID)
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

	password := data.GenEncryptedPassword()
	if data.PasswordIsEncrypted {
		password = data.Password
		log.Println("Password is already encrypted")
	}

	helpers.MyDB.Unsafe().Exec(sql, data.DomainID, password, data.Email)

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

func DeleteUser(userName string) {

	if userName == "" {
		return
	}

	// domain name
	var domain_name string
	sql := `select name from virtual_domains where id IN (select domain_id from virtual_users where email = ?) limit 1`
	helpers.MyDB.Unsafe().Get(&domain_name, sql, userName)

	if domain_name == "" {
		return
	}

	cmp := strings.Split(userName, "@")

	// delete dir
	os.RemoveAll("/var/mail/vhosts/" + domain_name + "/" + cmp[0] + "/")

	sql = `delete from virtual_users where email = ? limit 1`
	helpers.MyDB.Unsafe().Exec(sql, userName)

}
