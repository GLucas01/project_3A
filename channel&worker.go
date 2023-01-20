package main

import (
	"fmt"
	"sync"
)

type Matrix struct {
	values [][]int
}
type Task struct {
	id   int
	a, b Matrix
}
type Result struct {
	id  int
	res Matrix
}

var tasks = make(chan Task)
var results = make(chan Result)

func worker() {
	for {
		task, more := <-tasks
		if more {
			fmt.Printf("worker : commencé la tâche %d\n", task.id)

			res := Matrix{values: make([][]int, len(task.a.values))}
			for i := range res.values {
				res.values[i] = make([]int, len(task.b.values[0]))
			}
			for i := 0; i < len(task.a.values); i++ {
				for j := 0; j < len(task.b.values[0]); j++ {
					for k := 0; k < len(task.b.values); k++ {
						res.values[i][j] += task.a.values[i][k] * task.b.values[k][j]
					}
				}
			}
			fmt.Printf("worker : fini la tâche %d\n", task.id)
			results <- Result{id: task.id, res: res}
		} else {
			fmt.Println("worker : tâches terminées")
			return
		}
	}
}

func main() {

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			worker()
			wg.Done()
		}()
	}

	for i := 0; i < 10; i++ {
		tasks <- Task{id: i, a: Matrix{{values: [][]int{{1, 2, 3}, {4, 5, 6}}}}, b: Matrix{{values: [][]int{{1, 2}, {3, 4}, {5, 6}}}}}
	}
	close(tasks)
	wg.Wait()

	for i := 0; i < 10; i++ {
		result := <-results
		fmt.Printf("Resultat de la tache %d: %v\n", result.id, result.res)
	}
}
