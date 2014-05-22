package main

import (
	"fmt"
	"os"
	"sort"

	"code.google.com/p/gowut/gwu"
)

func startWebServer() {
	tables := gwu.NewHorizontalPanel()
	tables.SetVAlign(gwu.VA_TOP)
	updateTables(tables)

	win := gwu.NewWindow("scrolls", "ClockworkAgent")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HA_CENTER)
	win.Style().SetBackground(`
		url(images/background.jpg) no-repeat center center fixed;
		background-size: cover;
	`)
	win.Add(tables)
	win.AddEHandlerFunc(func(e gwu.Event) {
		updateTables(tables)
	}, gwu.ETYPE_WIN_LOAD)

	server := gwu.NewServer("", "localhost:8080")
	server.AddWin(win)
	wd, err := os.Getwd()
	deny(err)
	server.AddStaticDir("images", wd+"/images")
	server.Start("scrolls")
}

type StockedItem struct {
	card Card
	buy  int
	sell int
}
type SortBySellPrice []StockedItem

func (a SortBySellPrice) Len() int           { return len(a) }
func (a SortBySellPrice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortBySellPrice) Less(i, j int) bool { return a[i].sell > a[j].sell }

func updateTables(tables gwu.Panel) {
	var sorted SortBySellPrice
	for card, _ := range PriceHistory {
		sorted = append(sorted, StockedItem{
			card,
			currentState.DeterminePrice(card, 1, true),
			currentState.DeterminePrice(card, 1, false),
		})
	}
	sort.Sort(sorted)

	tables.Clear()
	tables.Add(createFAQ())

	tables.AddHSpace(100)
	for rarity, rarityStr := range []string{"Commons", "Uncommons", "Rares"} {
		layout := gwu.NewPanel()
		tables.Add(layout)
		tables.AddHSpace(30)

		layout.SetHAlign(gwu.HA_CENTER)
		layout.Style().SetBackground("rgba(0,0,0,0.5)")
		layout.Style().SetWidthPx(320)

		header := gwu.NewLabel(rarityStr)
		header.Style().SetColor("rgb(255,255,255)")
		header.Style().SetFontWeight("bold")
		header.Style().SetFontSize("x-large")
		layout.Add(header)

		table := gwu.NewTable()
		layout.Add(table)

		table.SetCellPadding(5)

		row := 0
		for _, sortedItem := range sorted {
			if CardRarities[sortedItem.card] != rarity {
				continue
			}

			table.Add(gwu.NewImage("", fmt.Sprintf("images/%s.png", CardResources[sortedItem.card])), row, 0)
			table.Add(gwu.NewLabel(string(sortedItem.card)), row, 1)
			table.Add(gwu.NewLabel(fmt.Sprintf("%d", sortedItem.buy)), row, 2)
			table.Add(gwu.NewLabel(fmt.Sprintf("%d", sortedItem.sell)), row, 3)
			if row%2 == 0 {
				table.RowFmt(row).Style().SetBackground("rgba(255,255,255,0.75)")
			} else {
				table.RowFmt(row).Style().SetBackground("rgba(150,150,150,0.75)")
			}
			row++
		}
	}
}

func createFAQ() gwu.Panel {
	panel := gwu.NewPanel()
	panel.Style().SetBackground("rgba(255,255,255,0.75)")

	addHeader := func(text string) {
		label := gwu.NewLabel(text)
		label.Style().SetColor("rgb(0,0,0)")
		label.Style().SetFontWeight("bold")
		label.Style().SetFontSize("large")
		panel.Add(label)
	}

	addText := func(text string) {
		label := gwu.NewLabel(text)
		label.Style().SetColor("rgb(0,0,0)")
		panel.Add(label)
	}

	addHeader("What is this?")
	addText("This site displays the prices for which the ClockworkAgent, a trading bot for the game Scrolls, will buy and sell cards for.")

	addHeader("How do I engage in trading with the bot?")
	addText("Just join the ingame channel 'clockwork' and say '!trade'. You will then be queued up for interaction with the bot in a trade.")

	addHeader("Why are these prices so different from Scrollsguide prices?")
	addText("Scrollsguide prices are determined from WTB and WTS messages in the Trading-channel. Thus they reflect what people expect to " +
		"pay/get for a card, not necessarily what the card is actually traded for. Since most people adjust their expectations to what " +
		"the current Scrollsguide price is, this can lead to a self-fulfilling prophecy. Also, it is pretty easy to manipulate the prices " +
		"for cards that are traded less often, enabling a way to scam the bot if it would use these prices.")

	addHeader("How then are these prices calculated?")
	addText("The price starts at 1500 for rares, 750 for uncommons and 187 for commons. Each time a card is sold to the bot, it will assume " +
		"that the card is less valuable, reducing the price by 100 / 50 / 12.5 depending on rarity. Each time a card is bought from the bot, " +
		"the price will	go up again.")

	addHeader("Why does the buy price fluctuate, when the sell price remains constant?")
	addText("If the bot has less than 2000 gold, the buy prices will be lowered down to a minimum of 50%. This way the bot can aquire more cards, " +
		"balancing out the fact that it had to overpay for most cards in order to determine the price, as well as new additions and general price deflation.")

	addHeader("Who created the bot, and where can I find the source code?")
	addText("Ingame: redefiance")
	addText("Reddit: lando-garner")
	panel.Add(gwu.NewLink("Source", "https://github.com/redefiance/ScrollsTradeBot"))

	return panel
}
