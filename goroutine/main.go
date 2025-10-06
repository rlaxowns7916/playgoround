package main

import (
	"playgoround/goroutine/pipeline"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	//http.Execute(&wg)
	pipeline.Execute(&wg)
}
