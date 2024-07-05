package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/marco-m/clim"
)

func main() {
	os.Exit(mainInt())
}

func mainInt() int {
	err := mainErr(os.Args[1:])
	if err == nil {
		return 0
	}
	fmt.Println(err)
	if errors.Is(err, clim.ErrHelp) {
		return 0
	}
	if errors.Is(err, clim.ErrParse) {
		return 2
	}
	return 1
}

type Application struct {
	timeout time.Duration
	sleep   time.Duration
}

func mainErr(args []string) error {
	app := Application{
		sleep: 20 * time.Millisecond,
	}
	cli := clim.New("timeout", "uses context", app.run)

	cli.AddFlag(&clim.Flag{Value: clim.Duration(&app.timeout, 100*time.Millisecond),
		Long: "timeout", Desc: "Context timeout (eg: 1h34m20s4ms)"})

	action, err := cli.Parse(args)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), app.timeout)
	defer cancel()

	return action(user{ctx})
}

type user struct {
	ctx context.Context
}

func (app *Application) run(user user) error {
	count := 1
	for {
		select {
		case <-user.ctx.Done():
			return user.ctx.Err()
		default:
			count = doSomething(app.sleep, count)
		}
	}
}

func doSomething(sleep time.Duration, count int) int {
	time.Sleep(sleep)
	fmt.Println(count)
	count++
	return count
}
