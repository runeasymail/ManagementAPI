package modules

import (
	"os"
	"fmt"
)

var (
	trusted_host string = "/etc/opendkim/TrustedHosts"
	keytable string = "/etc/opendkim/KeyTable"
	signinTable string = "/etc/opendkim/SigningTable"

)



func AddNewDomain(Domain string) {

	f, _ := os.OpenFile(trusted_host, os.O_APPEND, 0666)
	f.WriteString("*"+Domain)
	f.Close()

	f, _ = os.OpenFile(keytable, os.O_APPEND, 0666)
	f.WriteString( fmt.Sprintf("mail._domainkey.%s %s:mail:/etc/opendkim/keys/%s/mail.private",Domain,Domain,Domain))
	f.Close()

	f, _ = os.OpenFile(signinTable, os.O_APPEND, 0666)
	f.WriteString( fmt.Sprintf("*@%s mail._domainkey.%s",Domain,Domain))
	f.Close()


	// create dir
	os.MkdirAll("/etc/opendkim/keys/" + Domain, os.ModePerm)



}
