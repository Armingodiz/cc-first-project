package userService

import (
	"context"
	"github.com/google/uuid"
	"cc-first-project/user-service/models"
	"cc-first-project/user-service/store"
)

type UserService interface {
	CreateAdvertisement(ad *models.Advertisement) (string, error)
	GetAdvertisement(adId string) (*models.Advertisement, error)
}

var simpleContext = context.Background()

func NewUserService(store store.Store) UserService {
	return &userService{
		Store: store,
	}
}

type userService struct {
	Store store.Store
}

func (s *userService) CreateAdvertisement(ad *models.Advertisement) (string, error) {
	ad.State = models.AdvertisementStatePending
	id := uuid.New()
    ad.Id = id.String()
	return ad.Id, s.Store.CreateAdvertisement(simpleContext, ad)
}

func (s *userService) GetAdvertisement(adId string) (*models.Advertisement, error) {
	return s.Store.GetAdvertisement(simpleContext, adId)
}
