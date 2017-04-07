package modules

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/exec"
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

func letsEncryptInit(Hostname string) {

	// create a /ssl dir
	os.MkdirAll("/ssl", os.ModePerm)

	// mkdir -p /var/www/challenges/
	os.MkdirAll("/var/www/challenges/", os.ModePerm)

	// generate ssl certf
	commands := []cmds{}

	commands = append(commands, cmds{Program: "openssl", Command: []string{"genrsa", "4096", ">", "/ssl/account.key"}})
	//commands = append(commands, cmds{Program: "openssl", Command: []string{"genrsa", "4096", ">", "/ssl/domain.key"}})
	//commands = append(commands, cmds{Program: "openssl", Command: []string{"req", "-new", "-sha256", "-key", "/ssl/domain.key", "-subj", "/CN=" + Hostname, ">", "/ssl/domain.csr"}})
	//commands = append(commands, cmds{Program: "service", Command: []string{"nginx", "reload"}})
	//commands = append(commands, cmds{Program: "python", Command: []string{"/ssl/acme_tiny.py", "--account-key", "/ssl/account.key", "--csr", "/ssl/domain.csr", "--acme-dir", "/var/www/challenges/", ">", "/ssl/signed.crt"}})
	//commands = append(commands, cmds{Program: "wget", Command: []string{"-O", "-", "https://letsencrypt.org/certs/lets-encrypt-x3-cross-signed.pem", ">", "/ssl/intermediate.pem"}})
	//commands = append(commands, cmds{Program: "cat", Command: []string{"signed.crt", "/ssl/intermediate.pem", ">", "/ssl/chained.pem"}})

	for el := range commands {
		log.Println( commands[el].Program )
		log.Println( commands[el].Command )
		out, er := exec.Command(commands[el].Program, commands[el].Command...).Output()
		if er != nil {
			log.Println(er)

		}
		log.Println(string(out))
	}

}
