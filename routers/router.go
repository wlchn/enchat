package routers

import (
	"enchat/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/join", &controllers.MainController{}, "post:Join")

	beego.Router("/socket", &controllers.WebSocketController{})
	beego.Router("/socket/join", &controllers.WebSocketController{}, "get:Join")
}
