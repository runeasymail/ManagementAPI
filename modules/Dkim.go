package modules

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
	"os/exec"
)

var (
	trusted_host string = "/etc/opendkim/TrustedHosts"
	keytable     string = "/etc/opendkim/KeyTable"
	signinTable  string = "/etc/opendkim/SigningTable"
)

func HandlerNewDkimDomain(c *gin.Context) {

	// validate hostname

	res := add("mail3.yuks.me")

	c.JSON(200, gin.H{"result": res})
}

func add(Domain string) string {

	var log = logging.MustGetLogger("mail")

	log.Debug("Adding ", Domain, " in dkim registers")

	f, _ := ioutil.ReadFile(trusted_host)
	content := string(f)
	content = fmt.Sprintf("%s \n %s \n", content, "*"+Domain)

	log.Debug("Original ", trusted_host, "content: ", content)

	ioutil.WriteFile(trusted_host, []byte(content), os.ModePerm)

	f, _ = ioutil.ReadFile(keytable)
	content = string(f)
	log.Debug("Original ", keytable, "content: ", content)
	content = fmt.Sprintf("%s \n %s \n", content, fmt.Sprintf("mail._domainkey.%s %s:mail:/etc/opendkim/keys/%s/mail.private", Domain, Domain, Domain))
	ioutil.WriteFile(keytable, []byte(content), os.ModePerm)

	f, _ = ioutil.ReadFile(signinTable)
	content = string(f)
	log.Debug("Original ", signinTable, "content: ", content)
	content = fmt.Sprintf("%s \n %s \n", content, fmt.Sprintf("*@%s mail._domainkey.%s", Domain, Domain))
	ioutil.WriteFile(signinTable, []byte(content), os.ModePerm)

	// create dir
	os.MkdirAll("/etc/opendkim/keys/"+Domain, os.ModePerm)

	os.Chdir("/etc/opendkim/keys/" + Domain)

	_, er := exec.Command("opendkim-genkey", []string{"-s", "mail", "-d", Domain}...).Output()
	if er != nil {
		log.Debug("opendkim-genkey error", er)
	}

	_, er = exec.Command("chown", []string{"opendkim:opendkim", "mail.private"}...).Output()
	if er != nil {
		log.Debug("chown error", er)
	}

	_, er = exec.Command("service", []string{"postfix", "restart"}...).Output()
	if er != nil {
		log.Debug("service postfixx restart error", er)
	}

	_, er = exec.Command("service", []string{"opendkim", "restart"}...).Output()
	if er != nil {
		log.Debug("service opendkim restart error", er)
	}


	content_of_dkim, _ := ioutil.ReadFile("/etc/opendkim/keys/" + Domain + "/mail.txt" )

	return string(content_of_dkim)

}
