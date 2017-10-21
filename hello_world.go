package main

import (
	"fmt"
	"strings"

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

type WalletsController struct {
	beego.Controller
}

func (this *WalletsController) Get() {
	wallets, err := NewWallets("3000")
	addresses := wallets.GetAddresses()
	bc := NewBlockchain("3000")
	defer bc.db.Close()

	if err != nil {
		this.Ctx.WriteString("<html><body>Invalid Wallet</body></html>")
	}

	this.Ctx.WriteString("<html><body><table border=\"1\"><tr> <td>Address</td> <td>Balance</td> </tr>")
	for _, address := range addresses {
		var lines []string
		lines = append(lines, fmt.Sprintf("<tr><td><a href=\"/wallet/%s\">%s</a></td>", address, address))
		lines = append(lines, fmt.Sprintf("<td>%d</td></tr>", bc.getBalance(address)))
		this.Ctx.WriteString(strings.Join(lines, ""))
	}
	this.Ctx.WriteString("</table></body></html>")
}

func main() {
	beego.Router("/", &MainController{})
	beego.Router("/block/:id", &BlockController{})
	beego.Router("/wallets", &WalletsController{})
	beego.Run()
}
