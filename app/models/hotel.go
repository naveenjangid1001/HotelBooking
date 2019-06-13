package models

import (
	"github.com/revel/revel"
)

type Hotel struct {
	HotelId          int
	Name, Address    string
	City, State, Zip string
	Country          string
	Rating           float32
	Price            int
}

func GetHotels() ([]Hotel, error) {
	Dbmap.AddTableWithName(Hotel{}, "Hotel")
	var hotels []Hotel
	_, err := Dbmap.Select(&hotels, "select * from Hotel;")
	return hotels, err
}

func (hotel *Hotel) Validate(v *revel.Validation) {
	v.Check(hotel.Name,
		revel.Required{},
		revel.MaxSize{35},
	)

	v.MaxSize(hotel.Address, 50)

	v.Check(hotel.City,
		revel.Required{},
		revel.MaxSize{20},
	)

	v.Check(hotel.State,
		revel.Required{},
		revel.MaxSize{15},
	)

	v.Check(hotel.Zip,
		revel.Required{},
		revel.MaxSize{6},
		revel.MinSize{5},
	)

	v.Check(hotel.Country,
		revel.Required{},
		revel.MaxSize{25},
		revel.MinSize{2},
	)
}
