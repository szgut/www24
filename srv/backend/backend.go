package backend

import "time"
import "log"
import "github.com/szgut/www24/srv/core"

type Backend interface {
	Command(team core.Team, cmd core.Command) core.CommandResult
	Wait()
}

type backend struct {
	ch   chan commandMessage
	wait func()
}

type commandMessage struct {
	team core.Team
	cmd  core.Command
	ch   chan<- core.CommandResult
}

func (b *backend) Command(team core.Team, cmd core.Command) core.CommandResult {
	ch := make(chan core.CommandResult)
	b.ch <- commandMessage{team, cmd, ch}
	return <-ch
}

func (b *backend) Wait() {
	b.wait()
}

func StartNew(tickInterval int, game core.Game) Backend {
	cmdCh := make(chan commandMessage)
	tickWait, tickCh := newTicker(tickInterval)
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
				msg.ch <- game.Execute(msg.team, msg.cmd)
			}
		}
	}()
	return &backend{ch: cmdCh, wait: tickWait}
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
				log.Println("waiting:", queue)
			}
			select {
			case <-tickCh:
				backendTickCh.notify()
				backendTickCh.wait()
				for _, listener := range queue {
					listener.notify()
				}
				queue = nil
				tickCh.notify()

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
