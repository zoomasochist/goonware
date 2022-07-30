package daemon

import (
	types "goonware/types"

	"math/rand"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	StatePaused = iota
	StateRunning
)

func Tick(c types.Config, pkg types.EdgewarePackage) {
	glfw.Init()
	// We take gooning seriously
	rand.Seed(time.Now().UnixNano())

	ws := make(chan int32)
	if c.DriveFiller {
		go WorkerManager(ws, &c, &pkg, DoDriveFiller)
	}
	if c.Annoyances {
		go WorkerManager(ws, &c, &pkg, DoAnnoyances)
	}

	// Hibernation mode
	if c.Mode == types.ModeHibernate {
		for {
			hibernationCountdown :=
				rand.Intn(int(c.HibernateMaxWaitMinutes)-int(c.HibernateMinWaitMinutes)+1) +
					int(c.HibernateMinWaitMinutes)

			time.Sleep(time.Duration(hibernationCountdown) * time.Minute)
			ws <- StateRunning
			time.Sleep(time.Duration(c.HibernateActivityLength) * time.Second)
			ws <- StatePaused
		}
		// Normal mode
	} else {
		ws <- StateRunning
		// Recieve from a nil channel; block main thread execution indefinitely.
		<-(chan int)(nil)
	}
}

func WorkerManager(ws <-chan int32, c *types.Config, pkg *types.EdgewarePackage,
	workerFunc func(*types.Config, *types.EdgewarePackage)) {
	state := StatePaused

	for {
		select {
		case c := <-ws:
			state = int(c)
		default:
			if state == StateRunning {
				workerFunc(c, pkg)
			}
		}
	}
}
