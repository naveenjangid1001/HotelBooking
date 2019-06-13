package controllers

import (
	"booking/app/models"
	"booking/app/routes"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.ViewArgs["user"] = user
	}
	return nil
}

func (c App) connected() *models.User {
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username.(string))
	}

	return nil
}

func (c App) getUser(username string) (user *models.User) {
	user = &models.User{}
	c.Session.GetInto("fulluser", user, false)

	if user.Username == username {
		return user
	}

	models.Dbmap.AddTableWithName(models.User{}, "User").SetKeys(false, "Username")
	obj, err := models.Dbmap.Get(models.User{}, username)
	var dbUser *models.User
	if err == nil && obj != nil {
		dbUser = obj.(*models.User)
	}
	c.Session["fulluser"] = dbUser
	return dbUser
}

func (c App) Index() revel.Result {
	if c.connected() != nil {
		return c.Redirect(routes.Hotels.Index())
	}
	c.Flash.Error("Please log in first")
	return c.Render()
}

func (c App) Register() revel.Result {
	return c.Render()
}

func (c App) Login(username, password string, remember bool) revel.Result {
	user := c.getUser(username)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = username
			if remember {
				c.Session.SetDefaultExpiration()
			} else {
				c.Session.SetNoExpiration()
			}
			c.Flash.Success("Welcome, " + username)
			return c.Redirect(routes.Hotels.Index())
		}
	}

	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed")
	return c.Redirect(routes.App.Index())
}

func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Index())
}

func (c App) SaveUser(user models.User, confirmPass string) revel.Result {
	c.Validation.Required(confirmPass)
	c.Validation.Required(confirmPass == user.Password).Message("Passwords does not match.")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Register())
	}
	models.Dbmap.AddTableWithName(models.User{}, "User").SetKeys(false, "Username")
	user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	err := models.Dbmap.Insert(&user)

	if err != nil {
		log.Panic(err)
	}

	c.Session["user"] = user.Username
	c.Flash.Success("Welcome, " + user.Name)

	return c.Redirect(routes.Hotels.Index())
}
