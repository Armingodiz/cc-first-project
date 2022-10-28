package userService

import (
	"context"
    "io"
    "log"
    "net/http"
    "os"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"cc-first-project/user-service/models"
	"cc-first-project/user-service/store"
)

type UserService interface {
	CreateAdvertisement(ad *models.Advertisement) (string, error)
	GetAdvertisement(adId string) (*models.Advertisement, error)
	UploadImageFile(adId string, imageUrl string) error
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

func (s *userService) UploadImageFile(adId string, imageUrl string) error {
	fname := "image.jpg"
    f, err := os.Create(fname)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    res, err := http.Get(imageUrl)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()
    _, err = io.Copy(f, res.Body)

    if err != nil {
        log.Fatal(err)
    }
    file, err := os.Open(fname)
    if err != nil {
        return err
    }
    defer file.Close()

    sess, err := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials("0f4b4ae8-3d7f-4051-9c5a-6e4c36cba55b", "0f16304e084b0ce416c28a3a6780eac1485831bdfe034c039356767bee37623b", ""),
        Region:      aws.String("default"),
        Endpoint:    aws.String("https://s3.ir-thr-at1.arvanstorage.com"),
    })
    uploader := s3manager.NewUploader(sess)
    _, err = uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String("arminccproject"),
        Key: aws.String(adId),
        Body: file,
    })
	return err
}
