package controllers

import (
	"booking/app/models"
	"booking/app/routes"
	"fmt"
	"log"
	"strings"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type Hotels struct {
	App
}

func (c Hotels) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.App.Index())
	}
	return nil
}

func (c Hotels) Index() revel.Result {
	c.Log.Info("Fetching Index...")
	var bookings []*models.Booking
	models.Dbmap.AddTableWithName(models.Booking{}, "Booking")
	_, err := models.Dbmap.Select(&bookings, "select * from Booking where UserId=?;", c.connected().UserId)
	var hotel *models.Hotel
	models.Dbmap.AddTableWithName(models.Hotel{}, "Hotel").SetKeys(true, "HotelId")
	for i, _ := range bookings {
		obj, _ := models.Dbmap.Get(&hotel, bookings[i].HotelId)
		bookings[i].Hotel = obj.(*models.Hotel)
	}
	if err != nil {
		panic(err)
	}
	return c.Render(bookings)
}

func (c Hotels) List(search string, size, page uint32) revel.Result {
	if page == 0 {
		page = 1
	}
	nextPage := page + 1
	search = strings.TrimSpace(search)

	var hotels []*models.Hotel
	models.Dbmap.AddTableWithName(models.Hotel{}, "Hotel")
	models.Dbmap.Select(&hotels, "SELECT * FROM Hotel WHERE Name LIKE '%"+search+"%' or City Like '%"+search+"%' or State Like '%"+search+"%' or Country Like '%"+search+"%' limit ?, ?;", (page-1)*size, size)
	fmt.Println(hotels, search, size, page, nextPage)
	return c.Render(hotels, search, size, page, nextPage)
}

func (c Hotels) loadHotelById(id int) *models.Hotel {
	models.Dbmap.AddTableWithName(models.Hotel{}, "Hotel").SetKeys(true, "HotelId")
	h, err := models.Dbmap.Get(&models.Hotel{}, id)
	if err != nil {
		panic(err)
	}
	if h == nil {
		return nil
	}
	return h.(*models.Hotel)
}

func (c Hotels) Show(id int) revel.Result {
	hotel := c.loadHotelById(id)
	if hotel == nil {
		c.NotFound("Hotel does not exist.")
	}
	title := hotel.Name
	return c.Render(title, hotel)
}

func (c Hotels) Settings() revel.Result {
	return c.Render()
}

func (c Hotels) SaveSettings(password, verifyPassword string) revel.Result {
	models.ValidatePassword(c.Validation, password)
	c.Validation.Required(verifyPassword).Message("Please verify your password.")
	c.Validation.Required(password == verifyPassword).Message("Your password doesn't match.")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		return c.Redirect(routes.Hotels.Settings())
	}

	user := c.connected()
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.HashedPassword = bcryptPassword
	user.Password = password

	models.Dbmap.AddTableWithName(models.User{}, "User").SetKeys(true, "UserId")
	models.Dbmap.Query("update User set Password=?, HashedPassword=? where UserId=?;", user.Password, user.HashedPassword, user.UserId)
	return c.Redirect(routes.Hotels.Index())
}

func (c Hotels) CancelBooking(id int) revel.Result {
	models.Dbmap.AddTableWithName(models.Booking{}, "Booking").SetKeys(true, "BookingId")
	_, err := models.Dbmap.Exec("delete from Booking where BookingId=?;", id)
	if err != nil {
		c.Flash.Error("Error occurrured while cancelling your booking. Please try after some time.")
		return c.Redirect(routes.Hotels.Index())
	}
	c.Flash.Success("Booking cancelled.")
	return c.Redirect(routes.Hotels.Index())
}

func (c Hotels) Book(id int) revel.Result {
	hotel := c.loadHotelById(id)
	if hotel == nil {
		return c.NotFound("Hotel %d does not exist", id)
	}

	title := "Book " + hotel.Name
	return c.Render(title, hotel)
}

func (c Hotels) ConfirmBooking(id int, booking models.Booking) revel.Result {
	hotel := c.loadHotelById(id)
	if hotel == nil {
		return c.NotFound("Hotel %d does not exist", id)
	}

	//title := fmt.Sprintf("Confirm %s booking", hotel.Name)
	booking.Hotel = hotel
	booking.User = c.connected()
	booking.Validate(c.Validation)

	// if c.Validation.HasErrors() {
	// 	c.Validation.Keep()
	// 	c.FlashParams()
	// 	return c.Redirect(routes.Hotels.Book(id))
	// }

	models.Dbmap.AddTableWithName(models.Booking{}, "Booking").SetKeys(true, "BookingId")
	_, err := models.Dbmap.Query("insert into Booking values(?,?,?,?,?,?,?,?,?,?,?);", 0, booking.User.UserId, booking.Hotel.HotelId, booking.CheckInStr, booking.CheckOutStr, booking.CardNumber, booking.NameOnCard, booking.CardExpMonth, booking.CardExpYear, booking.Smoking, booking.Beds)
	if err != nil {
		panic(err)
	}
	c.Flash.Success("Thank you, %s, your booking in %s is confirmed.",
		booking.User.Name, booking.Hotel.Name)
	return c.Redirect(routes.Hotels.Index())
}

func (c Hotels) Dashboard() revel.Result {
	user := c.connected()
	return c.Render(user)
}

func (c Hotels) AddHotel(hotel models.Hotel) revel.Result {
	c.Validation.Required(hotel.Name)
	c.Validation.Required(hotel.Address)
	c.Validation.Required(hotel.City)
	c.Validation.Required(hotel.State)
	c.Validation.Required(hotel.Country)
	c.Validation.Required(hotel.Zip)
	c.Validation.Required(hotel.Price)

	if c.Validation.HasErrors() {
		c.Flash.Error("Please provide all the details.")
		return c.Redirect(routes.Hotels.Dashboard())
	}

	models.Dbmap.AddTableWithName(models.Hotel{}, "Hotel")
	err := models.Dbmap.Insert(&hotel)
	if err != nil {
		c.Flash.Error("Unable to add hotel to system. Our team looking into this.")
		log.Print(err)
		return c.Redirect(routes.Hotels.Dashboard())
	}
	c.Flash.Success("Hotel saved to the system.")
	return c.Redirect(routes.Hotels.Dashboard())
}
