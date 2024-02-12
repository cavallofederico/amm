package main

import (
	"context"
	"fmt"

	"github.com/InjectiveLabs/sdk-go/client/common"
	exchangeclient "github.com/InjectiveLabs/sdk-go/client/exchange"

	spotExchangePB "github.com/InjectiveLabs/sdk-go/exchange/spot_exchange_rpc/pb"
)

func main() {
	//network := common.LoadNetwork("mainnet", "k8s")
	network := common.LoadNetwork("testnet", "lb")
	exchangeClient, err := exchangeclient.NewExchangeClient(network)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	marketStatus := "active"
	quoteDenom := "peggy0x87aB3B4C8661e07D6372361211B96ed4Dc36B1B5"

	// req := derivativeExchangePB.MarketsRequest{
	req := spotExchangePB.MarketsRequest{
		MarketStatus: marketStatus,
		QuoteDenom:   quoteDenom,
	}

	// res, err := exchangeClient.GetDerivativeMarkets(ctx, &req)
	res, err := exchangeClient.GetSpotMarkets(ctx, &req)
	if err != nil {
		panic(err)
	}

	// str, _ := json.MarshalIndent(res, "", " ")
	// fmt.Print(string(str))

	for _, market := range res.Markets {
		fmt.Println(market.MarketId, market.MarketStatus, market.Ticker)
	}

	res2, err2 := exchangeClient.GetSpotMarket(ctx, "0x0611780ba69656949525013d947713300f56c37b6175e02f26bffa495c3208fe")

	if err2 != nil {
		panic(err)
	}

	// res2.Market.MarketStatus

	fmt.Println(res2)
}
