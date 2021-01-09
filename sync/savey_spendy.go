package main

import (
	"sync"
	"time"
)

var (
	money         = 100
	moneyWithLock = 100
	lock          = sync.Mutex{}
)

func spendy() {
	for i := 1; i < 1000; i++ {
		money -= 10
		time.Sleep(1 * time.Millisecond)
	}
	println("Spendy done...")
}

func savey() {
	for i := 1; i < 1000; i++ {
		money += 10
		time.Sleep(1 * time.Millisecond)
	}
	println("Savey done...")
}

func action() { // at times, race condition will happen
	go spendy()
	go savey()
	time.Sleep(2000 * time.Millisecond)
	println(money)
}

func spendyWithLock() {
	for i := 1; i < 1000; i++ {
		lock.Lock()
		moneyWithLock -= 10
		lock.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
	println("Spendy done...")
}

func saveyWithLock() {
	for i := 1; i < 1000; i++ {
		lock.Lock()
		moneyWithLock += 10
		lock.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
	println("Savey done...")
}

func actionWithLock() {
	go spendyWithLock()
	go saveyWithLock()
	time.Sleep(3000 * time.Millisecond)
	println(moneyWithLock)
}

func main() {
	action()
	actionWithLock()
}
