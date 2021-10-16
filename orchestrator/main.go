package main

import (
	"fmt"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/maratona-run-time/Maratona-Runtime/model"
	"github.com/maratona-run-time/Maratona-Runtime/queue"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
)

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
		client := graphql.NewClient("http://orm:8084/graphql", nil)
		err := utils.SaveSubmissionStatus(client, id, model.PENDING, "")
		if err != nil {
			msg := "An error occurred while trying to save submission '" + id + "' '" + model.PENDING + "' status"
			logger.Error().
				Err(err).
				Msg(msg)
		}
		logger.Info().Msgf("Creating mart Pod for submission %v", id)
		createPod(id)
	}
}
