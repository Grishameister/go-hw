package main

import (
	"sort"
	"strconv"
	"sync"
	//"sync/atomic"
)

func ExecutePipeline(jobs ...job) {
	input := make(chan interface{}, 100)
	output := make(chan interface{}, 100)
	wg := sync.WaitGroup{}
	for _, task := range jobs {
		wg.Add(1)
		go worker(task, input, output, &wg)
		input = output
		output = make(chan interface{}, 100)
	}
	wg.Wait()
}

func worker(task job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)
	task(in, out)
}

func SingleHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for val := range in {
		wg.Add(1)

		data := strconv.Itoa(val.(int))
		go func(out chan interface{}, wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()
			chF := make(chan string)
			chS := make(chan string)
			go func(ch chan string) {
				val := <-ch
				ch <- DataSignerCrc32(val)
			}(chF)
			go func(ch chan string) {
				val := <-ch
				mu.Lock()
				data := DataSignerMd5(val)
				mu.Unlock()
				ch <- DataSignerCrc32(data)
			}(chS)

			chF <- data
			chS <- data
			first := <-chF
			second := <-chS
			out <- first + "~" + second

		}(out, &wg, &mu)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}

	for val := range in {
		wg.Add(1)
		data := val.(string)
		go func(out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			ch1 := make(chan string)
			ch2 := make(chan string)
			ch3 := make(chan string)
			ch4 := make(chan string)
			ch5 := make(chan string)
			ch6 := make(chan string)
			channels := []chan string{
				ch1,
				ch2,
				ch3,
				ch4,
				ch5,
				ch6,
			}

			for i, channel := range channels {
				go func(ch chan string, i int) {
					val := <-ch
					ch <- DataSignerCrc32(strconv.Itoa(i) + val)
				}(channel, i)
			}

			for _, channel := range channels {
				channel <- data
			}

			var result string
			for _, ch := range channels {
				result += <-ch
			}

			out <- result
		}(out, &wg)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	hashes := make([]string, 0)
	for val := range in {
		hashes = append(hashes, val.(string))
	}
	sort.Strings(hashes)

	var result string
	first := true
	for _, hash := range hashes {
		if first {
			result += hash
			first = false
		} else {
			result += "_" + hash
		}
	}
	out <- result
}
