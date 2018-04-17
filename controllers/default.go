package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.TplName = "welcome.html"
}

func (c *MainController) Join() {
	email := c.GetString("email")

	if len(email) == 0 {
		c.Redirect("/", 302)
		return
	}

	c.Redirect("/socket?email="+email, 302)

	return
}
