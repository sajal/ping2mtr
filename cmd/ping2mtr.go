package main

import (
	"fmt"
	"github.com/sajal/mtrparser"
	"github.com/sajal/ping2mtr"
)

func main() {
	raw := ping2mtr.Ping2MTR("205.251.242.160")
	fmt.Println(raw)
	res, err := mtrparser.NewMTROutPut(raw, "205.251.242.160", 10)
	fmt.Println(err)
	fmt.Println(res)

}
