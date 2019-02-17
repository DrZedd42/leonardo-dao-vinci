package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var (
	// This will be modified by multiple routines
	iteration = atomicInt{
		val: 0,
		mu:  sync.RWMutex{},
	}

	// maps image number to a set of user adresses that voted for it
	votes = make(map[int]map[string]struct{})
	port  = flag.String("port", "8080", "The port that the server will bind on")
)

type voteMessage struct {
	UserAddress string `json:"user_address"`
	Iteration   int    `json:"iteration"`
	Images      []int  `json:"images"`
}

type iterationMessage struct {
	Iteration string `json:"iteration"`
}

func main() {
	flag.Parse()

	content, err := ioutil.ReadFile("iteration")
	if err == nil {
		num, err := strconv.Atoi(string(content))
		if err != nil {
			log.Fatalf("expecting integer defined in iteration file, instead got %s", string(content))
		}
		iteration = atomicInt{
			val: num,
			mu:  sync.RWMutex{},
		}

	}

	go generateArt()

	imageHandler := http.StripPrefix("/images/", http.FileServer(http.Dir("images")))
	imageHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		imageHandler.ServeHTTP(w, r)
	}

	bindAddress := fmt.Sprintf("localhost:%s", *port)
	http.HandleFunc("/images/", addCors(imageHandlerFunc))
	http.HandleFunc("/iteration", addCors(handleIteration))
	http.HandleFunc("/vote", addCors(handleVotes))

	fmt.Printf("Current iteration: %d\n", iteration)
	fmt.Printf("Serving images at %s/images/<iteration>/<image>\n", bindAddress)
	fmt.Printf("Serving iteration at %s/iteration\n", bindAddress)
	fmt.Printf("Accepting votes at %s/vote\n", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, nil))
}

// Runs in perpetuity
func generateArt() {
	for {
		println()
		for im := 1; im < 11; im++ {
			fmt.Printf("Generating image %d\n", im)
		}

		iteration.Inc()
		i := iteration.Get()
		fmt.Printf("Iteration: %d\n", i)

		err := ioutil.WriteFile("iteration", []byte(strconv.Itoa(i)), os.ModePerm)
		if err != nil {
			fmt.Printf("could not persist iteration file: %v", err)
		}

		err = os.Chdir("images")
		if err != nil {
			fmt.Printf("could not cd to images: %v", err)
		}
		err = os.Symlink("test_images", fmt.Sprintf("%d", i))
		if err != nil {
			fmt.Printf("could not symlink test_images: %v", err)
		}
		err = os.Chdir("..")
		if err != nil {
			fmt.Printf("could not cd to .. : %v", err)
		}
		fmt.Printf("Press enter to move to the next iteration\n")
		fmt.Scanln()
	}

}

func addCors(handleFunc func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		handleFunc(w, r)
	}
}

func handleIteration(w http.ResponseWriter, r *http.Request) {
	i := iterationMessage{
		Iteration: strconv.Itoa(iteration.Get()),
	}
	data, err := json.Marshal(i)
	if err != nil {
		http.Error(w, "Could not marshal iteration", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func handleVotes(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	fmt.Printf("content: %s\n", string(content))
	v := &voteMessage{}
	err = json.Unmarshal(content, v)
	if err != nil {
		errmsg := "Could not unmarshal votes message"
		fmt.Printf("%s\n", errmsg)
		http.Error(w, errmsg, http.StatusBadRequest)
		fmt.Printf("offending content: %s\n", string(content))
		return
	}

	// Verify that this user has not voted so far in this iteration
	for _, voters := range votes {
		_, ok := voters[v.UserAddress]
		if ok {
			errmsg := fmt.Sprintf("User %q has already voted in this iteration", v.UserAddress)
			fmt.Printf("%s\n", errmsg)
			http.Error(w, errmsg, http.StatusBadRequest)
			return
		}
	}

	// Add the votes for this user
	for _, i := range v.Images {
		voters, ok := votes[i]
		if !ok {
			voters = make(map[string]struct{})
			votes[i] = voters
		}
		voters[v.UserAddress] = struct{}{}
	}
	w.WriteHeader(http.StatusOK)
}

type atomicInt struct {
	mu  sync.RWMutex
	val int
}

func (s *atomicInt) Get() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.val
}

func (s *atomicInt) Inc() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.val++
}
