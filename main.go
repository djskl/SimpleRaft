package main

import (
	"fmt"
	"SimpleRaft/settings"
)

func main()  {
	fmt.Println(settings.LEADER)
	fmt.Println(settings.FOLLOWER)
	fmt.Println(settings.CANDIDATE)
}
