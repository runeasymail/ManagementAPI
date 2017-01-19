package modules

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/runeasymail/ManagementAPI/models"
)

func HandlerGetAllDomains(c *gin.Context) {

	data := models.GetDomains()

	c.JSON(200, gin.H{"domains": data})
}

func HandlerAddNewDomain(c *gin.Context) {

	type formData struct {
		DomainName string `form:"domain" valid:"host,required"`
		Username   string `form:"username" valid:"email,required"`
		Password   string `form:"password" valid:"required"`
	}

	data := formData{}
	c.Bind(&data)

	is_valid, err := govalidator.ValidateStruct(data)

	if !is_valid {
		c.JSON(200, gin.H{"result": false, "error_msg": err.Error()})
		return
	}

	res, err := models.AddNewDomain(data.DomainName, data.Username, data.Password)

	if res == false {
		c.JSON(200, gin.H{"result": false, "error_msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"result": true})

}
