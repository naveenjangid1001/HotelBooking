package models

import (
	"time"

	"github.com/revel/revel"
)

type Booking struct {
	BookingId    int
	UserId       int
	HotelId      int
	CheckInStr   time.Time
	CheckOutStr  time.Time
	CardNumber   string
	NameOnCard   string
	CardExpMonth int
	CardExpYear  int
	Smoking      bool
	Beds         int

	//Transient
	//CheckInDate  time.Time
	//	CheckOutDate time.Time
	User  *User
	Hotel *Hotel
}

func (b Booking) Validate(v *revel.Validation) {
	v.Required(b.User)
	v.Required(b.Hotel)
	//v.Required(b.CheckInDate)
	//v.Required(b.CheckOutDate)
	//v.Match(b.CardNumber, regexp.MustCompile(`\d{16}`)).Message("Credit card number must be numeric and 16 digits.")
	v.Check(b.NameOnCard, revel.Required{},
		revel.MinSize{3},
		revel.MaxSize{20},
	)
}
