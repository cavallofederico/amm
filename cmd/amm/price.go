package main

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/InjectiveLabs/sdk-go/client/common"
	exchangeclient "github.com/InjectiveLabs/sdk-go/client/exchange"
)

func getPriceFromOrderbook(market string) (price float64) {
	network := common.LoadNetwork("testnet", "lb")
	exchangeClient, err := exchangeclient.NewExchangeClient(network)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	marketIds := []string{market}
	stream, err := exchangeClient.StreamSpotOrderbookV2(ctx, marketIds)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			res, err := stream.Recv()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(res.MarketId, res.Orderbook.Sells, len(res.Orderbook.Sells), len(res.Orderbook.Buys))
			sells := make([]float64, len(res.Orderbook.Sells))
			for i, sell := range res.Orderbook.Sells {
				sellPrice, err := strconv.ParseFloat(sell.Price, 64)
				if err != nil {
					fmt.Println(err)
					return
				}
				sells[i] = sellPrice
			}
			minSell := slices.Min(sells)
			fmt.Println(minSell)

			buys := make([]float64, len(res.Orderbook.Buys))
			for i, buy := range res.Orderbook.Buys {
				buyPrice, err := strconv.ParseFloat(buy.Price, 64)
				if err != nil {
					fmt.Println(err)
					return
				}
				buys[i] = buyPrice
			}
			maxBuy := slices.Max(buys)
			fmt.Println(maxBuy)

			return (minSell-maxBuy)/2 + maxBuy
		}
	}
}
