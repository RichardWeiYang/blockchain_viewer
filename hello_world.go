package main

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	bc := NewBlockchain("3000")
	defer bc.db.Close()

	this.Ctx.WriteString("<html>")
	this.Ctx.WriteString(bc.PrintHTML())
	this.Ctx.WriteString("</html>")
}

func main() {
	beego.Router("/", &MainController{})
	beego.Router("/block/:id", &BlockController{})
	//beego.Router("/blockchain", &BlockChain{})
	beego.Run()
}
