package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "wanglei.io"
	c.Data["Email"] = "geekwanglei@gmail.com"
	c.TplName = "index.tpl"
}
