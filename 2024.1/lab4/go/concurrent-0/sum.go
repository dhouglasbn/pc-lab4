package main

import (
	"fmt"
	"io/ioutil"
	"os"
    "sync"
)

type result struct {
	path string
	sum int
	err error
}

// read a file from a filepath and return a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v", filePath, err)
		return nil, err
	}
	return data, nil
}

// sum all bytes of a file
func sum(filePath string) (int, error) {
	data, err := readFile(filePath)
	if err != nil {
		return 0, err
	}

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}

	return _sum, nil
}


// print the totalSum for all files and the files with equal sum
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	var totalSum int64
	sums := make(map[int][]string)
	
    var wg sync.WaitGroup
    results := make(chan result, len(os.Args) - 1)

	for _, path := range os.Args[1:] {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			_sum, err := sum(filePath)
			results <- result{path: filePath, sum: _sum, err: err}
		}(path)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

    for res := range results {
		if res.err != nil {
			continue
		}
		totalSum += int64(res.sum)
		sums[res.sum] = append(sums[res.sum], res.path)
	}

	fmt.Println(totalSum)
	for sum, files := range sums {
		if len(files) > 1 {
			fmt.Printf("Sum %d: %v\n", sum, files)
		}
	}
}
