package main

import (
	"fmt"

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

type BlockController struct {
	beego.Controller
}

func (this *BlockController) Get() {
	var block *Block
	bc := NewBlockchain("3000")
	defer bc.db.Close()

	bci := bc.Iterator()

	for {
		block = bci.Next()

		hash := fmt.Sprintf("%x", block.Hash)
		if hash == this.Ctx.Input.Param(":id") {
			break
		}

		if len(block.PrevBlockHash) == 0 {
			block = nil
			break
		}
	}

	this.Ctx.WriteString("<html>")
	if block != nil {
		this.Ctx.WriteString(block.PrintHTML(true))
	} else {
		this.Ctx.WriteString("No such Block")
	}
	this.Ctx.WriteString("</html>")
}

func main() {
	beego.Router("/", &MainController{})
	beego.Router("/block/:id", &BlockController{})
	//beego.Router("/blockchain", &BlockChain{})
	beego.Run()
}
