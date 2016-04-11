package main

import (
	game "github.com/scp1513/san/game/master"
	gate "github.com/scp1513/san/gate/master"
	portal "github.com/scp1513/san/portal/master"
)

func main() {
	c := make(chan struct{}, 3)
	go func() {
		portal.Run()
		c <- struct{}{}
	}()
	go func() {
		gate.Run()
		c <- struct{}{}
	}()
	go func() {
		game.Run()
		c <- struct{}{}
	}()
	for i := 0; i < 3; i++ {
		<-c
	}
}
