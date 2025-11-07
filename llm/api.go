package llm

import (
	"encoding/json"
	"log"
	"musaic/util"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type LLMSingleQuery struct {
	Question string `json:"question" bson:"bson"`
}

var (
	InputChan    chan string
	ResponseChan chan LLMResp
	wg           sync.WaitGroup
)

func Init() {
	InputChan = make(chan string)
	ResponseChan = make(chan LLMResp)
	go RunAgent(InputChan, ResponseChan)
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down...")
		close(InputChan)
		close(ResponseChan)
		wg.Wait()

	}()
}

func HandleLLMSingleQuery(w http.ResponseWriter, r *http.Request) {
	log.Println("LLM api endpoint called")
	var singleQuery LLMSingleQuery
	if err := json.NewDecoder(r.Body).Decode(&singleQuery); err != nil {
		util.JsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Music API endpoint called with album: %v", singleQuery)
	InputChan <- singleQuery.Question
	log.Printf("Job added: %s", singleQuery)
	resp := <-ResponseChan
	if resp.Err != nil {
		util.JsonError(w, "LLM Failed", http.StatusBadRequest)
		return
	}

}
