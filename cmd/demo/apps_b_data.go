// Data of the crypto dashboard, from Tabler's shared/data/crypto-*.json.

package main

import "github.com/mind-vm/tabler-templ/blocks/crypto"

var cryptoMarkets = []crypto.Market{
	{Coin: "ETH", ChangeText: "+1.11", Price: "0.03530991", Volume: "168,726.341", Change: 1.11},
	{Coin: "XRP", ChangeText: "-0.98", Price: "0.00002269", Volume: "22,843.185", Change: -0.98},
	{Coin: "ETC", ChangeText: "+2.35", Price: "0.00014852", Volume: "529.209", Change: 2.35},
	{Coin: "LTC", ChangeText: "+3.89", Price: "0.00089821", Volume: "5,928.536", Change: 3.89},
	{Coin: "XMR", ChangeText: "+4.42", Price: "0.00310306", Volume: "1,382.046", Change: 4.42},
	{Coin: "BSC", ChangeText: "+0.79", Price: "0.00989344", Volume: "12,531.037", Change: 0.79},
	{Coin: "ABC", ChangeText: "-3.20", Price: "0.02613303", Volume: "21.597", Change: -3.2},
	{Coin: "TRX", ChangeText: "+0.44", Price: "0.0000027", Volume: "4,250.916", Change: 0.44},
	{Coin: "SC", ChangeText: "+2.23", Price: "0.00000002", Volume: "48.092", Change: 2.23},
	{Coin: "OMG", ChangeText: "+0.30", Price: "0.00000102", Volume: "9.695", Change: 0.3},
}

var cryptoSellOrders = []crypto.Order{
	{Price: "0.00091279", Amount: "0.55002151", Sum: "0.55002151"},
	{Price: "0.00092146", Amount: "1.97294243", Sum: "2.52296394"},
	{Price: "0.00094672", Amount: "0.67387767", Sum: "3.19684161"},
	{Price: "0.00095269", Amount: "0.93727595", Sum: "4.13411756"},
	{Price: "0.00095926", Amount: "0.8976753", Sum: "5.03179286"},
	{Price: "0.00097736", Amount: "0.94088124", Sum: "5.9726741"},
	{Price: "0.00100026", Amount: "0.51299752", Sum: "6.48567162"},
	{Price: "0.00102267", Amount: "1.18050103", Sum: "7.66617265"},
	{Price: "0.00104506", Amount: "1.17318909", Sum: "8.83936174"},
	{Price: "0.00104974", Amount: "2.19498873", Sum: "11.03435047"},
}

var cryptoBuyOrders = []crypto.Order{
	{Price: "0.00088217", Amount: "2.59949934", Sum: "2.59949934"},
	{Price: "0.00085874", Amount: "0.9852387", Sum: "3.58473804"},
	{Price: "0.00084954", Amount: "2.39434245", Sum: "5.97908049"},
	{Price: "0.0008322", Amount: "2.31112332", Sum: "8.29020381"},
	{Price: "0.00081524", Amount: "2.66443099", Sum: "10.9546348"},
	{Price: "0.00080303", Amount: "1.79223619", Sum: "12.74687099"},
	{Price: "0.00077858", Amount: "1.61051825", Sum: "14.35738924"},
	{Price: "0.00075508", Amount: "0.71279599", Sum: "15.07018523"},
	{Price: "0.00073013", Amount: "0.89425132", Sum: "15.96443655"},
	{Price: "0.00071941", Amount: "1.85521536", Sum: "17.81965191"},
}

var cryptoSymbols = []string{"BTC", "ETH", "XRP", "BCH", "ADA", "LTC", "EOS", "XEM", "XLM", "NEO", "MIOTA", "DASH", "XMR", "TRX", "BTG", "ICX", "QTUM", "ETC", "LSK", "VET"}
