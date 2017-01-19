package modules

import (
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
