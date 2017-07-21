package models

import (
	"errors"
	"github.com/runeasymail/ManagementAPI/helpers"
	"os"
	"os/exec"
	"time"
	"io/ioutil"
)

type Domains struct {
	Id         uint64 `db:"id" json:"id" form:"id"`
	Name       string `db:"name" json:"name" form:"name" validation:"required"`
	UsersCount uint64 `db:"users_count" json:"users_count"`
}
type ExportUsers struct {
	Email string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type ExportDKIM struct {
	Public string `json:"public"`
	Private string `json:"private"`
}

type Export struct {
	Accounts []ExportUsers `json:"accounts"`
	Dkim ExportDKIM `json:"dkim"`
	DomainName string `json:"domainName"`
	GenTime string `json:"generatedAt"`
	Filename string `json:"export_filename"`
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

	if username != "" && password != "" {

		// add new User
		userData := Users{DomainID: uint64(id), Password: password, Email: username}
		_, err = AddNewUser(userData)

		if err != nil {
			return
		}

	}

	result = true
	return
}

func DeleteDomain(domain string)  {
	sql := `delete from virtual_domains where name = ? limit 1`
	helpers.MyDB.Unsafe().Exec(sql, domain)
	os.RemoveAll("/var/mail/vhosts/"+domain+"/")
}

func ExportToFile(domain string) (filename string, export_data Export, err error) {

	filename = "/var/mail/vhosts/" + domain + ".tar.gz"
	os.Chdir("/var/mail/vhosts/")

	cmd := []string{"-zcvf", domain+".tar.gz", domain}
	exec.Command("tar", cmd...).Output()


	base_path := "/opt/easymail/backupsByDomain/"+domain
	os.MkdirAll(base_path, os.ModePerm)

	current_time := time.Now().Local()

	old_file := filename
	filename = base_path + "/" + current_time.Format("2006-01-02") + ".tar.gz"

	cmd = []string{old_file, filename}
	_, err = exec.Command("mv", cmd...).Output()

	if err != nil {
		return
	}

	// db part
	var domain_id = ""
	sql := `select id from virtual_domains where name = ? limit 1`
	helpers.MyDB.Unsafe().Get(&domain_id,sql, domain)

	if domain_id == "" {
		return
	}

	sql = `select email,password from virtual_users where domain_id = ? order by id desc`
	helpers.MyDB.Unsafe().Select(&export_data.Accounts,sql, domain_id)


	export_data.GenTime = current_time.Format("2006-01-02 10:40:21")
	export_data.DomainName = domain
	export_data.Filename = filename


	//
	content_of_dkim, _ := ioutil.ReadFile("/etc/opendkim/keys/" + domain + "/mail.txt" )
	content_of_private_dkim, _ := ioutil.ReadFile("/etc/opendkim/keys/" + domain + "/mail.private" )

	export_data.Dkim = ExportDKIM{}
	export_data.Dkim.Public = string(content_of_dkim)
	export_data.Dkim.Private = string(content_of_private_dkim)


	return
}