package main

import "fmt"

func main() {
	marketId := "0x0611780ba69656949525013d947713300f56c37b6175e02f26bffa495c3208fe"
	price := getPriceFromOrderbook(marketId)
	print(price)
	conectionData, err := connect()

	if err != nil {
		fmt.Println(err)
		return
	}
	createBuyOrder(marketId, price*0.99, &conectionData)

}
