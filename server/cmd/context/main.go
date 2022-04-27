package main

import (
	"context"
	"fmt"
	"time"
)

type pramKey struct {
}

func main() {
	c := context.WithValue(context.Background(), pramKey{}, "abc") //key value形式
	c, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()
	go mainTask(c)
	var cmd string
	for {
		fmt.Scan(&cmd)
		if cmd == "c" {
			cancel()
		}
	}

}
func mainTask(c context.Context) {
	fmt.Printf("main task started with %q\n", c.Value(pramKey{}))
	go func() {
		c1, cancel := context.WithTimeout(c, 21*time.Second) //c1 是子任务
		defer cancel()
		smallTask(c1, "task1", 10*time.Second)
	}()
	smallTask(c, "task2", 5*time.Second)
}
func smallTask(c context.Context, name string, d time.Duration) {
	fmt.Printf("%s:main task started with %q\n", name, c.Value(pramKey{}))
	select {
	case <-time.After(d):
		fmt.Printf("%s done\n", name)
	case <-c.Done():
		fmt.Printf("%s cancel\n", name)
	}
}
