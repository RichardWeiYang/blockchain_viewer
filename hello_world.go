package main

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/richardweiyang/blockchain_go"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	chain := bc.NewBlockchain("3000")
	defer chain.DB.Close()

	this.Ctx.WriteString("<html>")
	this.Ctx.WriteString(chain.PrintHTML())
	this.Ctx.WriteString("</html>")
}

type BlockController struct {
	beego.Controller
}

func (this *BlockController) Get() {
	var block *bc.Block
	chain := bc.NewBlockchain("3000")
	defer chain.DB.Close()

	bci := chain.Iterator()

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
	wallets, err := bc.NewWallets("3000")
	addresses := wallets.GetAddresses()
	chain := bc.NewBlockchain("3000")
	defer chain.DB.Close()

	if err != nil {
		this.Ctx.WriteString("<html><body>Invalid Wallet</body></html>")
	}

	this.Ctx.WriteString("<html><body><table border=\"1\"><tr> <td>Address</td> <td>Balance</td> </tr>")
	for _, address := range addresses {
		var lines []string
		lines = append(lines, fmt.Sprintf("<tr><td><a href=\"/wallet/%s\">%s</a></td>", address, address))
		lines = append(lines, fmt.Sprintf("<td>%d</td></tr>", chain.GetBalance(address)))
		this.Ctx.WriteString(strings.Join(lines, ""))
	}
	this.Ctx.WriteString("</table></body></html>")
}

type WalletController struct {
	beego.Controller
}

func (this *WalletController) Get() {
	var lines []string
	chain := bc.NewBlockchain("3000")
	defer chain.DB.Close()

	address := this.Ctx.Input.Param(":id")

	this.Ctx.WriteString("<html><body>")
	if bc.ValidateAddress(address) {
		lines = append(lines, fmt.Sprintf("Balance: %d", chain.GetBalance(address)))
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

	if !bc.ValidateAddress(from) || !bc.ValidateAddress(to) {
		this.Ctx.WriteString("Check your Address")
		return
	}

	chain := bc.NewBlockchain("3000")
	UTXOSet := bc.UTXOSet{chain}
	defer chain.DB.Close()

	if amount > chain.GetBalance(from) {
		this.Ctx.WriteString("You don't have enough money")
		return
	}

	wallets, err := bc.NewWallets("3000")
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

	tx := bc.NewUTXOTransaction(wallet, to, amount, &UTXOSet)
	cbTx := bc.NewCoinbaseTX(from, "")
	txs := []*bc.Transaction{cbTx, tx}

	newBlock := chain.MineBlock(txs)
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
