package modules

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/runeasymail/ManagementAPI/models"
	"strconv"
)

func HandlerUserLists(c *gin.Context) {

	domain_id_string := c.Param("domain_id")

	domain_id, _ := strconv.ParseUint(domain_id_string, 10, 64)

	data := models.GetAllUsers(domain_id)

	if len(data) == 0 {
		data = []models.Users{}
	}

	c.JSON(200, gin.H{"users": data})
}

func HandlerUserPasswordChange(c *gin.Context) {

	data := models.Users{}
	c.Bind(&data)

	if data.DomainID != 0 && data.Id != 0 {
		go models.ChangePassword(data)
	}

	c.JSON(200,gin.H{"result":true})
}



func HandleUserAdd(c *gin.Context) {

	data := models.Users{}
	c.Bind(&data)

	is_valid, err := govalidator.ValidateStruct(data)

	if !is_valid {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}

	res, err := models.AddNewUser(data)

	if !res {
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"result": true})

}
