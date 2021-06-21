package data

import "github.com/clarke94/roulette-service/storage/bet"

var BetData = []bet.Bet{
	{
		ID:       "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		TableID:  "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		Bet:      "foo",
		Type:     "bar",
		Amount:   10,
		Currency: "GBP",
	},
}
