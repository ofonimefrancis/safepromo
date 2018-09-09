package event

import (
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	"github.com/ofonimefrancis/safepromo/config"
)

//GeoLocation Geofenced coordinates for the event
type GeoLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

//Event Houses the event and the promocode details associated with the event
type Event struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	PromoCode        string        `json:"promoCode"`
	PromoCodeAmount  string        `json:"promoCodeAmount"`
	NumberOfRides    float64       `json:"allowedRides"`
	CurrentRideCount float64       `json:"ridesTaken"`
	EventLocation    []GeoLocation `json:"geofence"`
	IsEnabled        bool          `json:"isEnabled"`
	CreatedAt        time.Time     `json:"createdAt"`
	StartDate        time.Time     `json:"startDate"`
	EndDate          time.Time     `json:"endDate"`
}

func (e Event) getPromoCodeAmount() string {
	return e.PromoCodeAmount
}

//IsExpired Checks if the Promocode is Expired based on the current time
func (e Event) IsExpired(currentTime time.Time) bool {
	return currentTime.Before(e.StartDate) || currentTime.After(e.EndDate)
}

//AddEvent Adds an Event and Returns the PromoCode for the event
func (e Event) AddEvent() error {
	session := config.Get().Session.Copy()
	defer session.Close()

	e.PromoCode = uuid.Must(uuid.NewRandom()).String() //Using UUID (Generate Promo Code)
	e.CreatedAt = time.Now()
	e.IsEnabled = false //By default promo code is inactive
	e.CurrentRideCount = 0
	collection := session.DB(config.DATABASENAME).C(config.EVENTSCOLLECTION)
	return collection.Insert(e)
}

//GetActivePromoCodes Return all active promo codes
func (e Event) GetActivePromoCodes() ([]Event, error) {
	session := config.Get().Session.Copy()
	defer session.Close()
	var activePromos []Event
	collection := session.DB(config.DATABASENAME).C(config.EVENTSCOLLECTION)
	err := collection.Find(bson.M{"isenabled": true}).All(&activePromos)
	return activePromos, err
}

//FindByPromoCode Find an event by PromoCode
func (e Event) FindByPromoCode(promoCode string) (Event, error) {
	session := config.Get().Session.Copy()
	defer session.Close()
	var event Event
	collection := session.DB(config.DATABASENAME).C(config.EVENTSCOLLECTION)
	err := collection.Find(bson.M{"promocode": promoCode}).One(&event)
	return event, err
}

//GetAllPromoCodes Returns all Events
func (e Event) GetAllPromoCodes() ([]Event, error) {
	session := config.Get().Session.Copy()
	defer session.Close()
	var allEvents []Event
	collection := session.DB(config.DATABASENAME).C(config.EVENTSCOLLECTION)
	err := collection.Find(bson.M{}).All(&allEvents)
	return allEvents, err
}

//Deactivate Promo Code Can be deactivated
func (e Event) Deactivate(promoCode string) error {
	session := config.Get().Session.Copy()
	defer session.Close()
	collection := session.DB(config.DATABASENAME).C(config.EVENTSCOLLECTION)
	return collection.Update(bson.M{"promocode": promoCode}, bson.M{"isenabled": false})
}

//Activate Activates a promo code
func (e Event) Activate(promoCode string) error {
	session := config.Get().Session.Copy()
	defer session.Close()
	if e.IsExpired(time.Now()) {
		return errors.New("You cannot activate an expired promo code")
	}
	collection := session.DB(config.DATABASENAME).C(config.EVENTSCOLLECTION)
	return collection.Update(bson.M{"promocode": promoCode}, bson.M{"isenabled": true})
}
