package app

import (
	"log"
	"cc-first-project/advertisement-service/services/brokerService"
	"cc-first-project/advertisement-service/services/emailService"
	"cc-first-project/advertisement-service/services/imageService"
	"cc-first-project/advertisement-service/store"
)

type App struct {
	BrokerService brokerService.BrokerService
	MailService   emailService.MailService
	AdStore store.Store
	ImageService imageService.ImageService
}

func NewApp() *App {
	Store := store.NewStore(getPostgresConnection())
	if Store == nil {
		log.Fatal("Failed to setup database")
	}
	return &App{
		BrokerService: brokerService.NewBrokerService(),
		MailService:   emailService.NewMailService(config.Configs.App.SenderEmail),
		AdStore: Store,
		ImageService: imageService.NewImageService(),
	}
}

func (a *App) Start() error {
	adChannel, errChann, err := a.BrokerService.StartConsuming()
	if err != nil {
		return err
	}
	defer a.BrokerService.CloseConnection()
	defer a.BrokerService.CloseChannel()
	go func() {
		for ad := range adChannel {
			log.Println("Received Ad: ", ad)
			category, err := a.ImageService.GetTag(ad.Image)
			var state string 
			if err == nil{
				state = models.Advertisement.AdvertisementStateAccepted
			}else{
				state = models.Advertisement.AdvertisementStateRejected
				log.Println("Error getting ad image category: ", err)
				errChann <- err
			}
			err = a.AdStore.SetCategory(ad.Id, category, state)
			if err != nil {
				log.Println("Error setting as state: ", err)
				errChann <- err
			}
			err := a.MailService.SendEmail(ad.Email, "your Ad was accepted")
			if err != nil {
				log.Println("Error sending email: ", err)
				errChann <- err
			}
		}
	}()
	for err := range errChann {
		log.Printf("Received error: %s", err.Error())
	}
	return nil
}

func getPostgresConnection() *sql.DB {
	serviceURI := "postgres://avnadmin:AVNS_8OSIiwW4r_EQh5LHb-u@pg-9240e5b-armingodarzi1380-0360.aivencloud.com:12488/defaultdb?sslmode=require"
	conn, _ := url.Parse(serviceURI)
	conn.RawQuery = "sslmode=verify-ca;sslrootcert=ca.pem"

	db, err := sql.Open("postgres", conn.String())

	if err != nil {
		log.Fatal(err)
	}
	return db
}