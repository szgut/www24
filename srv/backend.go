package main

import "time"
import "fmt"

type CommandMessage struct {
	Team Team
	Cmd  Command
	Ch   chan<- ResultMessage
}

type ResultMessage struct {
	Err    *CommandError
	Params []interface{}
}

func StartBackend(game Game) (ch chan<- CommandMessage, wait func()) {
	cmdCh := make(chan CommandMessage)
	tickWait, tickCh := newTicker(1)
	go func() {
		for {
			select {
			case <-tickCh:
				game.Tick()
				tickCh.notify()
				continue
			default:
			}
			select {
			case <-tickCh:
				game.Tick()
				tickCh.notify()
			case msg := <-cmdCh:
				params, err := game.Execute(msg.Team, msg.Cmd)
				msg.Ch <- ResultMessage{Err: err, Params: params}
			}
		}
	}()
	return cmdCh, tickWait
}

type notifier chan int

func (ch notifier) notify() {
	ch <- 0
}
func (ch notifier) wait() {
	<-ch
}

func newTicker(interval int) (func(), notifier) {
	tickCh := make(notifier)
	
	go func() {
		for {
			time.Sleep(time.Duration(interval) * time.Second)
			tickCh.notify()
			tickCh.wait()
		}
	}()

	listenCh := make(chan notifier)
	backendTickCh := make(notifier)

	go func() {
		queue := []notifier{}
		for {
			if len(queue) > 0 {
				fmt.Println(queue)
			}
			select {
			case <-tickCh:
				backendTickCh.notify()
				backendTickCh.wait()
				for _, listener := range queue {
					listener.notify()
				}
				queue = nil
				tickCh <- 0

			case listener := <-listenCh:
				queue = append(queue, listener)
			}
		}
	}()

	wait := func() {
		ch := make(notifier)
		listenCh <- ch
		ch.wait()
	}

	return wait, backendTickCh
}
