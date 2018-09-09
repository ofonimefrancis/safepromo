package event

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

//Routes Events Routes
func Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Post("/", CreateEvent)
		r.Get("/all", AllPromoCodes)
		r.Get("/validity/{promoCode}", GetPromocodeValidity)
		r.Put("/activate/{promoCode}", ActivateCode)
		r.Put("/deactivate/{promoCode}", DeactivateCode)
	})

	return router
}

//Given a promoCode and the current date
//Say i want to order a ride today
//I want to check if the promocode is valid

/**
	Promocode is valid if
	- If the promocode is enabled
	- Used in between the event start and end date
	- If used within the event radius
**/
func GetPromocodeValidity(w http.ResponseWriter, r *http.Request) {
	promoCode := chi.URLParam(r, "promoCode")
	event, err := Event{}.FindByPromoCode(promoCode)
	if err != nil {
		log.Printf("Error: %v\n", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	//TODO: Change time.Now() to any time specified by the user
	if event.IsExpired(time.Now()) {
		log.Println("Error: Expired time")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{
			"success": false,
			"message": "Invalid Promo Code",
		})
		return
	}
	//TODO: Check if the promo code hasn't been used more than stipulated

	//TODO: The promo code is used within the stipulated event radius

	//If all goes well the promo code is valid
	log.Println("Success")
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, map[string]interface{}{
		"message": "Valid Promo Code",
	})
	return

}
