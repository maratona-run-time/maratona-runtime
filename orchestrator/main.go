package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/maratona-run-time/Maratona-Runtime/queue"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
)

var verdictResponseError = errors.New("Error on verdict response")

func callVerdict(id string) ([]byte, error) {
	res, err := utils.MakeSubmissionRequest("http://mart-verdict:8083", id)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, verdictResponseError
	}

	binary, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return binary, nil
}

func main() {
	logger, logFile := utils.InitLogger("orchestrator")
	defer logFile.Close()

	msgs, err := queue.GetSubmissionChannel()
	for err != nil {
		fmt.Println(err)
		msgs, err = queue.GetSubmissionChannel()
		time.Sleep(2 * time.Second)
	}
	for queueMessage := range msgs {
		id := string(queueMessage.Body)
		verdictResponse, err := callVerdict(id)
		if err != nil {
			logger.Error().Err(err).Msgf("An error occurred when calling verdict with submission id %v", id)
		}
		fmt.Println(verdictResponse)
	}
}
