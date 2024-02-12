package main

import (
	"context"
	"fmt"

	"github.com/InjectiveLabs/sdk-go/client/common"
	exchangeclient "github.com/InjectiveLabs/sdk-go/client/exchange"
	derivativeExchangePB "github.com/InjectiveLabs/sdk-go/exchange/derivative_exchange_rpc/pb"
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

	req := derivativeExchangePB.MarketsRequest{
		MarketStatus: marketStatus,
		QuoteDenom:   quoteDenom,
	}

	res, err := exchangeClient.GetDerivativeMarkets(ctx, &req)
	if err != nil {
		panic(err)
	}

	// str, _ := json.MarshalIndent(res, "", " ")
	// fmt.Print(string(str))

	for _, market := range res.Markets {
		fmt.Println(market.MarketId, market.MarketStatus, market.Ticker)
	}
}
