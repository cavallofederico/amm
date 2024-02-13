package main

import "fmt"

func main() {
	price := getPriceFromOrderbook("0x0611780ba69656949525013d947713300f56c37b6175e02f26bffa495c3208fe")
	fmt.Println(price)
}
