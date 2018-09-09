package event

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/go-chi/render"
)

//CreateEvent Creates a New Event and Generates a Promocode
//The creator sees a success and not the promo code created.
func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]interface{}
	//var event Event
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.Println("Error: Invalid Payload")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"message": "Invalid Payload",
			"error":   err,
		})
		return
	}

	fmt.Printf("%v\n", requestData)

	event := Event{
		Name:             requestData["name"].(string),
		NumberOfRides:    requestData["ridesTaken"].(float64),
		PromoCodeAmount:  requestData["promoCodeAmount"].(string), //TODO: Change to Integer
		CurrentRideCount: requestData["allowedRides"].(float64),
		StartDate:        convertStringDateToTime(requestData["startDate"].(string)),
		EndDate:          convertStringDateToTime(requestData["endDate"].(string)),
	}

	fmt.Printf("%v\n", event)

	err = event.AddEvent()
	if err != nil {
		log.Println("Error Adding Event ", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"success": true,
	})
}

func convertStringDateToTime(dateTime string) time.Time {
	t, _ := time.Parse(time.RFC1123, dateTime)
	return t
}

//AllPromoCodes Return an array of promocodes from all events
func AllPromoCodes(w http.ResponseWriter, r *http.Request) {
	events, err := Event{}.GetAllPromoCodes()
	if err != nil {
		log.Println(err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err)
		return
	}
	var promoCodes []string

	for i := 0; i < len(events); i++ {
		promoCodes = append(promoCodes, events[i].PromoCode)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, promoCodes)
	return
}

//DeactivateCode Deactivates a promo code
func DeactivateCode(w http.ResponseWriter, r *http.Request) {
	promoCode := chi.URLParam(r, "promoCode")
	if len(promoCode) == 0 {
		log.Println("Invalid Promo Code length")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "Invalid Promo Code length")
		return
	}
	err := Event{}.Deactivate(promoCode)
	if err != nil {
		log.Println("Error: ", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "Error Deactivating Code")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"success": true,
	})

}

//ActivateCode Activates a promo code
func ActivateCode(w http.ResponseWriter, r *http.Request) {
	promoCode := chi.URLParam(r, "promoCode")
	if len(promoCode) == 0 {
		log.Println("Invalid Promo Code length")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "Invalid Promo Code length")
		return
	}
	err := Event{}.Activate(promoCode)
	if err != nil {
		log.Println("Error: ", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "Error Activating Code")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"success": true,
	})

}
