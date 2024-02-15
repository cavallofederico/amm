package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/InjectiveLabs/sdk-go/client"
	"github.com/InjectiveLabs/sdk-go/client/common"
	exchangeclient "github.com/InjectiveLabs/sdk-go/client/exchange"
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	exchangetypes "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
	chainclient "github.com/InjectiveLabs/sdk-go/client/chain"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
)

type ConectionData struct {
	chainClient         chainclient.ChainClient
	defaultSubaccountID eth.Hash
	senderAddress       types.AccAddress
	marketsAssistant    chainclient.MarketsAssistant
	clientCtx           cosmosclient.Context
}

func connect() (ConectionData, error) {
	network := common.LoadNetwork("testnet", "lb")
	tmClient, err := rpchttp.New(network.TmEndpoint, "/websocket")
	if err != nil {
		panic(err)
	}

	senderAddress, cosmosKeyring, err := chainclient.InitCosmosKeyring(
		os.Getenv("HOME")+"/.injectived",
		"injectived",
		"file",
		"inj1cguu2tgxkge5k3gu8q7s88qaquf3e6asyyjry5", // cosmos key from
		"12345678", // passphrase
		"87cd67b525dcb9b30ad4c23bd7ad39d589b15a014eab7123fa0aeac85030af17", // priv key
		false,
	)

	if err != nil {
		panic(err)
	}

	clientCtx, err := chainclient.NewClientContext(
		network.ChainId,
		senderAddress.String(),
		cosmosKeyring,
	)
	if err != nil {
		fmt.Println(err)
		return ConectionData{}, err
	}
	clientCtx = clientCtx.WithNodeURI(network.TmEndpoint).WithClient(tmClient)

	exchangeClient, err := exchangeclient.NewExchangeClient(network)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	marketsAssistant, err := chainclient.NewMarketsAssistantInitializedFromChain(ctx, exchangeClient)
	if err != nil {
		panic(err)
	}

	chainClient, err := chainclient.NewChainClient(
		clientCtx,
		network,
		common.OptionGasPrices(client.DefaultGasPriceWithDenom),
	)

	if err != nil {
		panic(err)
	}

	defaultSubaccountID := chainClient.DefaultSubaccount(senderAddress)
	data := ConectionData{
		defaultSubaccountID: defaultSubaccountID,
		chainClient:         chainClient,
		senderAddress:       senderAddress,
		marketsAssistant:    marketsAssistant,
		clientCtx:           clientCtx,
	}
	return data, nil
}

func createBuyOrder(marketId string, priceFloat float64, conectionData *ConectionData) {
	amount := decimal.NewFromFloat(1)
	price := decimal.NewFromFloat(priceFloat * 1e12)

	order := conectionData.chainClient.CreateSpotOrder(
		conectionData.defaultSubaccountID,
		&chainclient.SpotOrderData{
			OrderType:    exchangetypes.OrderType_SELL, //BUY SELL BUY_PO SELL_PO
			Quantity:     amount,
			Price:        price,
			FeeRecipient: conectionData.senderAddress.String(),
			MarketId:     marketId,
			Cid:          uuid.NewString(),
		},
		conectionData.marketsAssistant,
	)

	msg := new(exchangetypes.MsgCreateSpotLimitOrder)
	msg.Sender = conectionData.senderAddress.String()
	msg.Order = exchangetypes.SpotOrder(*order)

	simRes, err := conectionData.chainClient.SimulateMsg(conectionData.clientCtx, msg)

	if err != nil {
		fmt.Println(err)
		return
	}

	msgCreateSpotLimitOrderResponse := exchangetypes.MsgCreateSpotLimitOrderResponse{}
	err = msgCreateSpotLimitOrderResponse.Unmarshal(simRes.Result.MsgResponses[0].Value)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("simulated order hash: ", msgCreateSpotLimitOrderResponse.OrderHash)

	//AsyncBroadcastMsg, SyncBroadcastMsg, QueueBroadcastMsg
	err = conectionData.chainClient.QueueBroadcastMsg(msg)

	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second * 5)

	gasFee, err := conectionData.chainClient.GetGasFee()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("gas fee:", gasFee, "INJ")
}
