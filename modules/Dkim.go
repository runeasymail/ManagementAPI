package modules

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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

	f, _ := ioutil.ReadFile(trusted_host)
	content := string(f)
	content = fmt.Sprintf(`%s \n %s`, content, "*"+Domain )

	ioutil.WriteFile(trusted_host, []byte(content), os.ModePerm)



	//f, er := os.OpenFile(trusted_host, os.O_APPEND, 0666)
	//if er != nil {
	//	log.Println(er)
	//}
	//
	//_, er = f.WriteString("*" + Domain)
	//if er != nil {
	//	log.Println(er)
	//}
	//er = f.Close()
	//if er != nil {
	//	log.Println(er)
	//}
	//
	//
	//
	//f, _ = os.OpenFile(keytable, os.O_APPEND, 0666)
	//f.WriteString(fmt.Sprintf("mail._domainkey.%s %s:mail:/etc/opendkim/keys/%s/mail.private", Domain, Domain, Domain))
	//f.Close()
	//
	//f, _ = os.OpenFile(signinTable, os.O_APPEND, 0666)
	//f.WriteString(fmt.Sprintf("*@%s mail._domainkey.%s", Domain, Domain))
	//f.Close()

	// create dir
	os.MkdirAll("/etc/opendkim/keys/"+Domain, os.ModePerm)

}
