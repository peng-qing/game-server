package main

import (
	"GameServer/utils"
	"fmt"
)

func main() {
	number := 0
	bs := utils.EncodeZigzag(number)

	res, err := utils.DecodeZigzag(bs)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
