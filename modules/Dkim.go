package modules

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"github.com/op/go-logging"
)

var (
	trusted_host string = "/etc/opendkim/TrustedHosts"
	keytable     string = "/etc/opendkim/KeyTable"
	signinTable  string = "/etc/opendkim/SigningTable"
)

func HandlerNewDkimDomain(c *gin.Context) {

	// validate hostname

	add("mail3.yuks.me")
}

func add(Domain string) {

	var log = logging.MustGetLogger("mail")

	log.Debug("Adding ", Domain, " in dkim registers")

	f, _ := ioutil.ReadFile(trusted_host)
	content := string(f)
	content = fmt.Sprintf("%s \n %s \n", content, "*"+Domain )

	log.Debug("Original ", trusted_host, "content: ", content)

	ioutil.WriteFile(trusted_host, []byte(content), os.ModePerm)


	f, _ = ioutil.ReadFile(keytable)
	content = string(f)
	log.Debug("Original ", keytable, "content: ", content)
	content = fmt.Sprintf("%s \n %s \n", content, fmt.Sprintf("mail._domainkey.%s %s:mail:/etc/opendkim/keys/%s/mail.private", Domain, Domain, Domain) )
	ioutil.WriteFile(keytable, []byte(content), os.ModePerm)


	f, _ = ioutil.ReadFile(signinTable)
	content = string(f)
	log.Debug("Original ", signinTable, "content: ", content)
	content = fmt.Sprintf("%s \n %s \n", content, fmt.Sprintf("*@%s mail._domainkey.%s", Domain, Domain) )
	ioutil.WriteFile(signinTable, []byte(content), os.ModePerm)


	
	// create dir
	os.MkdirAll("/etc/opendkim/keys/"+Domain, os.ModePerm)

}
