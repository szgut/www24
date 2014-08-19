package backend

import "time"
import "github.com/szgut/www24/srv/core"

type Backend interface {
	Command(team core.Team, cmd core.Command) core.CommandResult
	Wait()
}

type backend struct {
	game   core.Game
	cmdCh  chan commandMessage
	tickCh notifier
	wait   func()
}

type commandMessage struct {
	team core.Team
	cmd  core.Command
	ch   chan<- core.CommandResult
}

func (b *backend) Command(team core.Team, cmd core.Command) core.CommandResult {
	ch := make(chan core.CommandResult)
	b.cmdCh <- commandMessage{team, cmd, ch}
	return <-ch
}

func (b *backend) Wait() {
	b.wait()
}

func (b *backend) Run() {
	for {
		select {
		case <-b.tickCh:
			b.game.Tick()
			b.tickCh.notify()
			continue
		default:
		}
		select {
		case <-b.tickCh:
			b.game.Tick()
			b.tickCh.notify()
		case msg := <-b.cmdCh:
			msg.ch <- b.game.Execute(msg.team, msg.cmd)
		}
	}
}

func StartNew(tickInterval int, game core.Game) Backend {
	tck := newTicker(tickInterval)
	bend := backend{cmdCh: make(chan commandMessage), game: game, wait: tck.Wait, tickCh: tck.backendCh}
	tck.Start()
	go bend.Run()
	return &bend
}

type notifier chan int

func (ch notifier) notify() {
	ch <- 0
}
func (ch notifier) wait() {
	<-ch
}

type ticker struct {
	interval  int
	backendCh notifier
	listenCh  chan notifier
}

func (self *ticker) Start() {
	tickCh := make(notifier)
	go func() {
		for {
			tickCh.notify()
			tickCh.wait()
			time.Sleep(time.Duration(self.interval) * time.Second)
		}
	}()
	go func() {
		queue := []notifier{}
		for {
			select {
			case <-tickCh:
				self.backendCh.notify()
				self.backendCh.wait()
				for _, listener := range queue {
					listener.notify()
				}
				queue = nil
				tickCh.notify()
			case listener := <-self.listenCh:
				queue = append(queue, listener)
			}
		}
	}()
}

func (self *ticker) Wait() {
	ch := make(notifier)
	self.listenCh <- ch
	ch.wait()
}

func newTicker(interval int) ticker {
	return ticker{interval: interval, backendCh: make(notifier), listenCh: make(chan notifier)}
}
