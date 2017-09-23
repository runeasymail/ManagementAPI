package modules

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/runeasymail/ManagementAPI/helpers"
)

type cmds struct {
	Program string
	Command []string
}

func LetsEncryptHandler(c *gin.Context) {

	var log = logging.MustGetLogger("mail")

	hostname := helpers.Config.App.Hostname

	log.Debug("Generate SSL for ", hostname)

	acme_default_ca := c.PostForm("default_ca")

	status := letsEncryptInit(hostname, acme_default_ca)

	if status != nil {
		c.JSON(200, gin.H{"result": false, "error": fmt.Sprintf("%s", status)})
	} else {
		c.JSON(200, gin.H{"result": true})
	}

}

func UploadMySSLHandler(c *gin.Context) {

	//ssl_certificate, _ := c.FormFile("ssl_certificate")
	ssl_certificate, _, _ := c.Request.FormFile("ssl_certificate")
	ssl_certificate_key, _, _ := c.Request.FormFile("ssl_certificate_key")

	os.Remove("/etc/dovecot/dovecot.pem")
	os.Remove("/etc/dovecot/private/dovecot.pem")

	// create file
	out, _ := os.Create("/etc/dovecot/dovecot.pem")
	defer out.Close()
	io.Copy(out, ssl_certificate)
	// end of file creation

	// create file
	out, _ = os.Create("/etc/dovecot/private/dovecot.pem")
	defer out.Close()
	io.Copy(out, ssl_certificate_key)
	// end of file creation

	//restart services
	exec.Command("service", []string{"postfix", "restart"}...).Output()
	exec.Command("service", []string{"dovecot", "restart"}...).Output()
	exec.Command("service", []string{"nginx", "restart"}...).Output()

	c.JSON(200, gin.H{"result": true})
}

func CheckSSLisValidHandler(c *gin.Context) {
	hostname := c.PostForm("hostname")

	days_left := checkValid(hostname)

	c.JSON(200, gin.H{"days_left": days_left})
}

func checkValid(url string) int64 {
	fullURL := url + ":443"
	conn, er := tls.Dial("tcp", fullURL, &tls.Config{})
	if er != nil {
		cert := conn.ConnectionState().PeerCertificates[0]
		end := cert.NotAfter
		diff := end.Sub(time.Now())
		days_left := (diff / (time.Hour * 24)).Nanoseconds()
		return days_left
	}
	return 0
}

func get_external(url string, filename string) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Get External error", err.Error())
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	ioutil.WriteFile(filename, body, os.ModePerm)

}

func letsEncryptInit(Hostname string, acme_default_ca string) (err error) {

	var log = logging.MustGetLogger("mail")

	// create a /ssl dir
	os.RemoveAll("/ssl")
	os.MkdirAll("/ssl", os.ModePerm)

	// mkdir -p /var/www/challenges/
	os.MkdirAll("/var/www/challenges/", os.ModePerm)

	get_external("https://github.com/hlandau/acme/releases/download/v0.0.61/acmetool-v0.0.61-linux_amd64.tar.gz", "/ssl/acmetool-v0.0.61-linux_amd64.tar.gz")

	cmd := exec.Command("tar", []string{"-zxvf", "/ssl/acmetool-v0.0.61-linux_amd64.tar.gz"}...)
	cmd.Dir = "/ssl"
	cmd.Output()

	exec.Command("cp", []string{"/ssl/acmetool-v0.0.61-linux_amd64/bin/acmetool", "/ssl/acmetool"}...).Output()

	os.RemoveAll("/ssl/acmetool-v0.0.61-linux_amd64")

	os.MkdirAll("/var/lib/acme/conf", os.ModePerm)

	acme_config := `
request:
  provider: https://acme-v01.api.letsencrypt.org/directory
  key:
    type: rsa
  challenge:
    webroot-paths:
    - /var/www/challenges
`
	ioutil.WriteFile("/var/lib/acme/conf/target", []byte(acme_config), os.ModePerm)

	go func(Hostname string) {
		out, er := exec.Command("/ssl/acmetool", []string{"--xlog.stderr", "want", Hostname}...).Output()

		if er != nil {
			fmt.Println(er.Error())
		} else {
			log.Debug(out)
		}

		dir := `/var/lib/acme/live/` + Hostname

		exec.Command("cp", []string{dir + "/fullchain", "/etc/dovecot/dovecot.pem"}...).Output()
		exec.Command("cp", []string{dir + "/privkey", "/etc/dovecot/private/dovecot.pem"}...).Output()

		time.Sleep(5 * time.Second)
		exec.Command("service", []string{"nginx", "reload"}...).Output()

	}(Hostname)

	//
	//// ACME STG ... testing
	//if acme_default_ca == "staging" {
	//	readed , _ := ioutil.ReadFile("/ssl/acme_tiny.py")
	//	replaced_ := strings.Replace(string(readed), `DEFAULT_CA = "https://acme-v01.api.letsencrypt.org"`, `DEFAULT_CA = "https://acme-staging.api.letsencrypt.org"`, -1)
	//	ioutil.WriteFile("/ssl/acme_tiny.py", []byte(replaced_), os.ModePerm)
	//}
	//
	//out, er = exec.Command("python", []string{"/ssl/acme_tiny.py", "--account-key", "/ssl/account.key", "--csr", "/ssl/domain.csr", "--acme-dir", "/var/www/challenges/"}...).Output()
	//
	//if er != nil {
	//	log.Debug("Python acme Tiny Error", er)
	//}
	//
	//ioutil.WriteFile("/ssl/signed.crt", out, os.ModePerm)
	//
	//if string(out) == "" {
	//	log.Debug("LetsEncrypt output is empty")
	//	err = errors.New("Let's encrypt error")
	//	return
	//}
	//
	//get_external("https://letsencrypt.org/certs/lets-encrypt-x3-cross-signed.pem", "/ssl/intermediate.pem")
	//
	//out, er = exec.Command("cat", []string{"/ssl/signed.crt", "/ssl/intermediate.pem"}...).Output()
	//if er != nil {
	//	log.Debug("cat signed with intermediate error", er)
	//}
	//ioutil.WriteFile("/ssl/chained.pem", out, os.ModePerm)
	//
	//exec.Command("cp", []string{"/ssl/chained.pem", "/etc/dovecot/dovecot.pem"}...).Output()
	//exec.Command("cp", []string{"/ssl/domain.key", "/etc/dovecot/private/dovecot.pem"}...).Output()
	//
	//go func(){
	//	time.Sleep(5 * time.Second)
	//	out, er := exec.Command("service", []string{"nginx", "reload"}...).Output()
	//	if er != nil {
	//		log.Debug("NGINX Service Reload error ", er)
	//	} else {
	//		log.Debug("NGINX reload output", string(out))
	//	}
	//}()
	//
	//out, er = exec.Command("service", []string{"postfix", "reload"}...).Output()
	//if er != nil {
	//	log.Debug("Postfix Service Reload error ", er)
	//} else {
	//	log.Debug("Postfix reload output", string(out))
	//}
	//
	//out, er = exec.Command("service", []string{"dovecot", "restart"}...).Output()
	//if er != nil {
	//	log.Debug("Dovecot Service Restart error ", er)
	//} else {
	//	log.Debug("Dovecot reload output", string(out))
	//}

	return
}
