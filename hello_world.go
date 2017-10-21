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

type WalletController struct {
	beego.Controller
}

func (this *WalletController) Get() {
	var lines []string
	bc := NewBlockchain("3000")
	defer bc.db.Close()

	address := this.Ctx.Input.Param(":id")

	this.Ctx.WriteString("<html><body>")
	if ValidateAddress(address) {
		lines = append(lines, fmt.Sprintf("Balance: %d", bc.getBalance(address)))
		lines = append(lines, fmt.Sprintf("<form action=\"/wallet/%s\" method=\"post\">", address))
		lines = append(lines, fmt.Sprintf("From:<br> <input type=\"text\" name=\"from\" value=\"%s\"> <br>", this.Ctx.Input.Param(":id")))
		lines = append(lines, fmt.Sprintf("To:  <br> <input type=\"text\" name=\"to\"> <br>"))
		lines = append(lines, fmt.Sprintf("Amount:  <br> <input type=\"text\" name=\"amount\"> <br><br>"))
		lines = append(lines, fmt.Sprintf("<input type=\"submit\" value=\"Submit\"> </form>"))
		this.Ctx.WriteString(strings.Join(lines, ""))
	} else {
		this.Ctx.WriteString("Invalid Address!")
	}
	this.Ctx.WriteString("</body></html>")
}

func (this *WalletController) Post() {
	from := this.GetString("from")
	to := this.GetString("to")
	amount, _ := this.GetInt("amount")

	this.Ctx.WriteString("<html><body>")

	if !ValidateAddress(from) || !ValidateAddress(to) {
		this.Ctx.WriteString("Check your Address")
		return
	}

	bc := NewBlockchain("3000")
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	if amount > bc.getBalance(from) {
		this.Ctx.WriteString("You don't have enough money")
		return
	}

	wallets, err := NewWallets("3000")
	if err != nil {
		this.Ctx.WriteString("Check your wallet")
		return
	}
	wallet := wallets.GetWallet(from)
	if wallet == nil {
		fmt.Println()
		this.Ctx.WriteString("The Address doesn't belongs to you!")
		return
	}

	tx := NewUTXOTransaction(wallet, to, amount, &UTXOSet)
	cbTx := NewCoinbaseTX(from, "")
	txs := []*Transaction{cbTx, tx}

	newBlock := bc.MineBlock(txs)
	UTXOSet.Update(newBlock)

	this.Ctx.WriteString("Succeed!</br></body></html>")
	this.Ctx.WriteString(fmt.Sprintf("New block created <a href=\"/block/%x\">%x</a> to see your block</body></html>", newBlock.Hash, newBlock.Hash))
}

func main() {
	beego.Router("/", &MainController{})
	beego.Router("/block/:id", &BlockController{})
	beego.Router("/wallets", &WalletsController{})
	beego.Router("/wallet/:id", &WalletController{})
	beego.Run()
}
