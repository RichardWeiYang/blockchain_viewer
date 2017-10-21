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

	this.Ctx.WriteString(bc.Print())
}

func main() {
	beego.Router("/", &MainController{})
	//beego.Router("/blockchain", &BlockChain{})
	beego.Run()
}
