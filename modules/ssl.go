package modules

import (
	"github.com/gin-gonic/gin"
	"os"
	"os/exec"
	"io/ioutil"
	"log"
	"net/http"
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
	req, _ := http.NewRequest("POST", url, nil)
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

	// python acme_tiny.py --account-key ./account.key --csr ./domain.csr --acme-dir /var/www/challenges/ > ./signed.crt


	out, er = exec.Command("python", []string{"/ssl/acme_tiny.py", "--account-key", "/ssl/account.key", "--csr", "/ssl/domain.csr", "--acme-dir", "/var/www/challenges/" }...).Output()

	if er !=nil {
		log.Println(er)
	}
	ioutil.WriteFile("/ssl/signed.crt", out, os.ModePerm)


	get_external("https://letsencrypt.org/certs/lets-encrypt-x3-cross-signed.pem","/ssl/lets-encrypt-x3-cross-signed.pem")


	


	//// generate ssl certf
	//commands := []cmds{}
	//
	////commands = append(commands, cmds{Program: "openssl", Command: []string{"req", "-new", "-sha256", "-key", "/ssl/domain.key", "-subj", "/CN=" + Hostname, ">", "/ssl/domain.csr"}})
	////commands = append(commands, cmds{Program: "service", Command: []string{"nginx", "reload"}})
	////commands = append(commands, cmds{Program: "python", Command: []string{"/ssl/acme_tiny.py", "--account-key", "/ssl/account.key", "--csr", "/ssl/domain.csr", "--acme-dir", "/var/www/challenges/", ">", "/ssl/signed.crt"}})
	////commands = append(commands, cmds{Program: "wget", Command: []string{"-O", "-", "https://letsencrypt.org/certs/lets-encrypt-x3-cross-signed.pem", ">", "/ssl/intermediate.pem"}})
	////commands = append(commands, cmds{Program: "cat", Command: []string{"signed.crt", "/ssl/intermediate.pem", ">", "/ssl/chained.pem"}})
	//
	//for el := range commands {
	//	log.Println( commands[el].Program )
	//	log.Println( commands[el].Command )
	//	out, er := exec.Command(commands[el].Program, commands[el].Command...).Output()
	//	if er != nil {
	//		log.Println(er)
	//
	//	}
	//	log.Println(string(out))
	//}

}
