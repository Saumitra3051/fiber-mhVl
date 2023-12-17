package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	return port
}
package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type SortRequest struct {
	ToSort [][]int `json:"to_sort"`
}

type SortResponse struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNs       int64   `json:"time_ns"`
}

func sortSequential(toSort [][]int) [][]int {
	sortedArrays := make([][]int, len(toSort))
	for i, arr := range toSort {
		sortedArr := make([]int, len(arr))
		copy(sortedArr, arr)
		sort.Ints(sortedArr)
		sortedArrays[i] = sortedArr
	}
	return sortedArrays
}

func sortConcurrent(toSort [][]int) ([][]int, time.Duration) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	sortedArrays := make([][]int, len(toSort))
	startTime := time.Now()

	for i, arr := range toSort {
		wg.Add(1)
		go func(i int, arr []int) {
			defer wg.Done()
			sortedArr := make([]int, len(arr))
			copy(sortedArr, arr)
			sort.Ints(sortedArr)

			mu.Lock()
			sortedArrays[i] = sortedArr
			mu.Unlock()
		}(i, arr)
	}

	wg.Wait()
	elapsedTime := time.Since(startTime)
	return sortedArrays, elapsedTime
}

func processSingleHandler(w http.ResponseWriter, r *http.Request) {
	var req SortRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := sortSequential(req.ToSort)
	elapsedTime := time.Since(startTime)

	response := SortResponse{
		SortedArrays: sortedArrays,
		TimeNs:       elapsedTime.Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processConcurrentHandler(w http.ResponseWriter, r *http.Request) {
	var req SortRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	sortedArrays, elapsedTime := sortConcurrent(req.ToSort)

	response := SortResponse{
		SortedArrays: sortedArrays,
		TimeNs:       elapsedTime.Nanoseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}




func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			r := mux.NewRouter()
	r.HandleFunc("/process-single", processSingleHandler).Methods("POST")
	r.HandleFunc("/process-concurrent", processConcurrentHandler).Methods("POST")

	http.Handle("/", r)

	http.ListenAndServe(":8000", nil)
		})
	})
	
	

	app.Listen(getPort())
}
