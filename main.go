package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	numGoroutines = 5  // количество горутин
	numIterations = 10 // количество итераций на горутину
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// примеры
	useMutex()
	useSemaphore()
	useBarrier()
	useSpinLock()
}

// Mutex
func useMutex() {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	wg.Add(numGoroutines)
	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				mutex.Lock()
				randomChar := rune(rand.Intn(95) + 32) // Случайный символ ASCII
				fmt.Printf("Mutex: Goroutine %d: %c\n", id, randomChar)
				mutex.Unlock()
				time.Sleep(time.Millisecond * 100)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("Mutex время: %v\n", time.Since(start))
}

// Semaphore (реализация через канал)
func useSemaphore() {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 1) // Семафор с емкостью 1

	wg.Add(numGoroutines)
	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				sem <- struct{}{} // Захват семафора
				randomChar := rune(rand.Intn(95) + 32)
				fmt.Printf("Semaphore: Goroutine %d: %c\n", id, randomChar)
				<-sem // Освобождение семафора
				time.Sleep(time.Millisecond * 100)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("Semaphore время: %v\n", time.Since(start))
}

// Barrier (реализация через WaitGroup)
func useBarrier() {
	var wg sync.WaitGroup
	var barrier sync.WaitGroup

	wg.Add(numGoroutines)
	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				randomChar := rune(rand.Intn(95) + 32)
				fmt.Printf("Barrier: Goroutine %d: %c\n", id, randomChar)
				time.Sleep(time.Millisecond * 100)
				barrier.Wait() // Ожидание других горутин
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("Barrier время: %v\n", time.Since(start))
}

// SpinLock (реализация через atomic)
func useSpinLock() {
	var wg sync.WaitGroup
	var spinLock int32

	wg.Add(numGoroutines)
	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				for atomic.CompareAndSwapInt32(&spinLock, 0, 1) == false {
					// Spin
				}
				randomChar := rune(rand.Intn(95) + 32)
				fmt.Printf("SpinLock: Goroutine %d: %c\n", id, randomChar)
				atomic.StoreInt32(&spinLock, 0)
				time.Sleep(time.Millisecond * 100)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("SpinLock время: %v\n", time.Since(start))
}
