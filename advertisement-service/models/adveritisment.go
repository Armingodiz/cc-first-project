package models

type Advertisement struct {
	Id             	 string `json:"id" db:"id"`
	Description            string `json:"description" db:"description"`
	Email            string `json:"email" bson:"email"`
	Category         string `json:"category" db:"category"`
	State          string `json:"state" db:"state"`
	Image 		string `json:"image"`
}


const(
	AdvertisementStatePending = "pending"
	AdvertisementStateAccepted = "accepted"
	AdvertisementStateRejected = "rejected"
)