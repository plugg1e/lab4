package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Student struct {
	Name   string
	Course int
	Debts  int
}

func generateStudents(size int) []Student {
	students := make([]Student, size)
	for i := 0; i < size; i++ {
		students[i] = Student{
			Name:   fmt.Sprintf("Student-%d", i),
			Course: rand.Intn(4) + 1, // курсы 1 по 10
			Debts:  rand.Intn(10),    // долги 0 до 9 у меня кстати 3 долга не бейте меня молю
		}
	}
	return students
}

func sequentialProcessing(students []Student, K int) []string {
	var result []string
	for _, s := range students {
		if s.Debts > 3 && s.Course > K {
			result = append(result, s.Name)
		}
	}
	return result
}

func parallelProcessing(students []Student, K int, numThreads int) []string {
	var results []string
	var mutex sync.Mutex
	var wg sync.WaitGroup

	chunkSize := len(students) / numThreads
	for i := 0; i < numThreads; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numThreads-1 {
			end = len(students)
		}
		chunk := students[start:end]
		wg.Add(1)
		go func(chunk []Student) {
			defer wg.Done()
			var temp []string
			for _, s := range chunk {
				if s.Debts > 3 && s.Course > K {
					temp = append(temp, s.Name)
				}
			}
			mutex.Lock()
			results = append(results, temp...)
			mutex.Unlock()
		}(chunk)
	}
	wg.Wait()
	return results
}
