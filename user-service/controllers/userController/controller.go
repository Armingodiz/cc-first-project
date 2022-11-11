package userController

import (
	"cc-first-project/user-service/models"
	"cc-first-project/user-service/services/brokerService"
	"cc-first-project/user-service/services/userService"
	"net/http"
	"github.com/google/uuid"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService      userService.UserService
	BrokerService    brokerService.BrokerService
}

func (u *UserController) AddAdvertisement() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ad models.Advertisement
		if err := c.ShouldBindJSON(&ad); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ad.State = models.AdvertisementStatePending
		adid := uuid.New()
		ad.Id = adid.String()
		url, err := u.UserService.UploadImageFile(ad.Id, ad.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(url)
		//ad.Image = url
		id, err := u.UserService.CreateAdvertisement(&ad)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = u.BrokerService.Publish(ad)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "your advertisement created successfully with id: " + id})
	}
}

func (u *UserController) GetAdvertisement() gin.HandlerFunc {
	return func(c *gin.Context) {
		adId := c.Param("ad_id")
		ad, err := u.UserService.GetAdvertisement(adId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var res gin.H
		if ad.State == models.AdvertisementStateRejected {
			res = gin.H{"Status": "Your advertisement is rejected"}
		}else if ad.State == models.AdvertisementStateAccepted {
			res = gin.H{"Advertisment": ad}
		}else{
			res = gin.H{"Status": "Your advertisement is not approved yet"}
		}
		c.JSON(http.StatusOK, res)
	}
}