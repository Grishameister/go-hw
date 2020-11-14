package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

const capacity int = 100

func ExecutePipeline(jobs ...job) {
	input := make(chan interface{}, capacity)
	output := make(chan interface{}, capacity)
	wg := sync.WaitGroup{}
	for _, task := range jobs {
		wg.Add(1)
		go worker(task, input, output, &wg)
		input = output
		output = make(chan interface{}, capacity)
	}
	wg.Wait()
}

func worker(task job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)
	task(in, out)
}

func JobSingleHash(data string, out chan interface{}, wg *sync.WaitGroup, mu *sync.Mutex) {
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
}

func SingleHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for val := range in {
		wg.Add(1)
		data := strconv.Itoa(val.(int))
		go JobSingleHash(data, out, &wg, &mu)
	}
	wg.Wait()
}

func JobMultiHash(data string, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	var channels []chan string
	for i := 0; i < 6; i++ {
		ch := make(chan string)
		channels = append(channels, ch)
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
}

func MultiHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}

	for val := range in {
		wg.Add(1)
		go JobMultiHash(val.(string), out, &wg)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	hashes := make([]string, 0)
	for val := range in {
		hashes = append(hashes, val.(string))
	}
	sort.Strings(hashes)
	out <- strings.Join(hashes, "_")
}
