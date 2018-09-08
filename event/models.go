package event

import (
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	"github.com/ofonimefrancis/safepromo/config"
)

//Event Houses the event and the promocode details associated with the event
type Event struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	PromoCode        string    `json:"promoCode"`
	PromoCodeAmount  string    `json:"promoCodeAmount"` //PromoCode Worth
	NumberOfRides    int       `json:"allowedRides"`
	CurrentRideCount int       `json:"ridesTaken"`
	PromoCodeRadius  string    `json:"promoCodeRadius"` //LatLng radius of the event
	IsEnabled        bool      `json:"isEnabled"`
	CreatedAt        time.Time `json:"createdAt"`
	StartDate        time.Time `json:"startDate"`
	EndDate          time.Time `json:"endDate"` //Expiration date of promo
}

//TODO: Functionalities
// ● Generation of new promo codes for events
// ● The promo code is worth a specific amount of ride
// ● The promo code can expire
// ● Only valid when user’s pickup or destination is within x radius of the event venue
// ● The promo code radius should be configurable

func (e Event) getPromoCodeAmount() string {
	return e.PromoCodeAmount
}

//IsExpired Validates if the promocode is expired based on current Time
//If the current Date is before the startDate
//Or the current Date is after the endDate
//PromoCode is expired
func (e Event) IsExpired(currentTime time.Time) bool {
	return currentTime.Before(e.StartDate) || currentTime.After(e.EndDate)
}

//AddEvent Adds an Event and Returns the PromoCode for the event
func (e Event) AddEvent() error {
	session := config.Get().Session.Copy()
	defer session.Close()

	//TODO: Generate PROMOCODE For Event
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
//TODO: We only activate a promocode if the current time is not passed the endDate of the event
func (e Event) Activate(promoCode string) error {
	session := config.Get().Session.Copy()
	defer session.Close()
	if e.IsExpired(time.Now()) {
		return errors.New("You cannot activate an expired promo code")
	}
	collection := session.DB(config.DATABASENAME).C(config.EVENTSCOLLECTION)
	return collection.Update(bson.M{"promocode": promoCode}, bson.M{"isenabled": true})
}
