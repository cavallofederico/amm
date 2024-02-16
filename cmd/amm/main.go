package main

import (
	"fmt"
	"time"

	exchangetypes "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
)

func main() {
	marketId := "0x0611780ba69656949525013d947713300f56c37b6175e02f26bffa495c3208fe"
	price := getPriceFromOrderbook(marketId)
	print(price)
	conectionData, err := connect()

	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Second * 5)
	fmt.Println(price * 0.999)
	createOrder(marketId, price*0.995, &conectionData, exchangetypes.OrderType_BUY, 10)
	time.Sleep(time.Second * 5)
	fmt.Println(price * 1.001)
	createOrder(marketId, price*1.005, &conectionData, exchangetypes.OrderType_SELL, 10)

}
