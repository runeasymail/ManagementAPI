package modules

import (
	"github.com/gin-gonic/gin"
	"os"
	"os/exec"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type cmds struct {
	Program string
	Command []string
}

func LetsEncryptHandler(c *gin.Context) {

	hostname := c.PostForm("hostname")

	letsEncryptInit(hostname)

	c.JSON(200, gin.H{"result": true})

}


func get_external(url string, filename string) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	ioutil.WriteFile(filename, body, os.ModePerm)

}

func letsEncryptInit(Hostname string) {

	// create a /ssl dir
	os.MkdirAll("/ssl", os.ModePerm)

	// mkdir -p /var/www/challenges/
	os.MkdirAll("/var/www/challenges/", os.ModePerm)



	out, _ := exec.Command("openssl", []string{"genrsa", "4096"}...).Output()
	ioutil.WriteFile("/ssl/account.key", out, os.ModePerm)

	out, _ = exec.Command("openssl", []string{"genrsa", "4096"}...).Output()
	ioutil.WriteFile("/ssl/domain.key", out, os.ModePerm)

	out, er := exec.Command("openssl", []string{"req", "-new", "-sha256", "-key", "/ssl/domain.key", "-subj", "/CN=" + Hostname }...).Output()
	if er != nil {
		log.Println(er)
	}
	ioutil.WriteFile("/ssl/domain.csr", out, os.ModePerm)

	get_external("https://raw.githubusercontent.com/diafygi/acme-tiny/master/acme_tiny.py","/ssl/acme_tiny.py")

	// ACME STG ... testing
	readed , _ := ioutil.ReadFile("/ssl/acme_tiny.py")
	replaced_ := strings.Replace(string(readed), `DEFAULT_CA = "https://acme-v01.api.letsencrypt.org"`, `DEFAULT_CA = "https://acme-staging.api.letsencrypt.org"`, -1)
	ioutil.WriteFile("/ssl/acme_tiny.py", []byte(replaced_), os.ModePerm)



	out, er = exec.Command("python", []string{"/ssl/acme_tiny.py", "--account-key", "/ssl/account.key", "--csr", "/ssl/domain.csr", "--acme-dir", "/var/www/challenges/" }...).Output()

	if er !=nil {
		log.Println(er)
	}
	ioutil.WriteFile("/ssl/signed.crt", out, os.ModePerm)


	get_external("https://letsencrypt.org/certs/lets-encrypt-x3-cross-signed.pem","/ssl/lets-encrypt-x3-cross-signed.pem")

	out, er = exec.Command("cat", []string{"/ssl/signed.crt","/ssl/intermediate.pem" }...).Output()
	if er !=nil {
		log.Println(er)
	}
	ioutil.WriteFile("/ssl/chained.pem", out, os.ModePerm)


	exec.Command("cp", []string{"/ssl/chained.pem","/etc/dovecot/dovecot.pem" }...).Output()
	exec.Command("cp", []string{"/ssl/domain.key","/etc/dovecot/private/dovecot.pem" }...).Output()

	// replace nginx config
	lets_encry_nginx_config := `

	location /.well-known/acme-challenge/ {
		alias /var/www/challenges/;
		try_files $uri =404;
	}

	`
	nginx_config, _ := ioutil.ReadFile("/etc/nginx/sites-enabled/roundcube")
	nginx_new := strings.Replace( string(nginx_config), "# __EASY_MAIL_INCLUDE_LETSENCRYPT__", lets_encry_nginx_config, -1 )
	ioutil.WriteFile("/etc/nginx/sites-enabled/roundcube", []byte(nginx_new), os.ModePerm)



	out, er = exec.Command("service", []string{"nginx","reload" }...).Output()
	if er !=nil {
		log.Println(er)
	} else {
		log.Println(string(out))
	}

}
