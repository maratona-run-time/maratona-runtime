package orm

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/graphql-go/graphql"

	"github.com/maratona-run-time/Maratona-Runtime/model"
)

var testFile = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TestFile",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.ID,
			},
			"fileName": &graphql.Field{
				Type: graphql.String,
			},
			"content": &graphql.Field{
				Type: graphql.NewList(graphql.Int),
			},
		},
	},
)

var challenge = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Challenge",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.ID,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"timeLimit": &graphql.Field{
				Type: graphql.Float,
			},
			"memoryLimit": &graphql.Field{
				Type: graphql.Int,
			},
			"inputs": &graphql.Field{
				Type: graphql.NewList(testFile),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					challenge := p.Source.(model.Challenge)
					return model.InputsArray(challenge.Inputs).TestFiles(), nil
				},
			},
			"outputs": &graphql.Field{
				Type: graphql.NewList(testFile),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					challenge := p.Source.(model.Challenge)
					return model.OutputsArray(challenge.Outputs).TestFiles(), nil
				},
			},
		},
	},
)

var submission = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Submission",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.ID,
			},
			"language": &graphql.Field{
				Type: graphql.String,
			},
			"source": &graphql.Field{
				Type: graphql.NewList(graphql.Int),
			},
			"challenge": &graphql.Field{
				Type: challenge,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					submission := p.Source.(model.Submission)
					return FindChallenge(submission.ChallengeID)
				},
			},
			"verdict": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					submission := p.Source.(model.Submission)
					return submission.Status.Verdict, nil
				},
			},
			"message": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					submission := p.Source.(model.Submission)
					return submission.Status.Message, nil
				},
			},
		},
	},
)

var queries = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"challenges": &graphql.Field{
				Type: graphql.NewList(challenge),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					challenges, err := FindAllChallenges()
					return challenges, err
				},
			},
			"challenge": &graphql.Field{
				Type: challenge,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					stringId := p.Args["id"].(string)
					id, err := strconv.ParseUint(stringId, 10, 64)
					if err != nil {
						return nil, errors.New("Could not convert id field to an uint")
					}
					return FindChallenge(uint(id))
				},
			},
			"submission": &graphql.Field{
				Type: submission,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.ID,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					stringId := p.Args["id"].(string)
					id, err := strconv.ParseUint(stringId, 10, 64)
					if err != nil {
						return nil, errors.New("Could not convert id field to an uint")
					}
					return FindSubmission(uint(id))
				},
			},
		},
	},
)

var mutations = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"judge": &graphql.Field{
				Type: submission,
				Args: graphql.FieldConfigArgument{
					"submissionID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
					"verdict": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"message": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					stringSubmissionID := p.Args["submissionID"].(string)
					submissionID, err := strconv.ParseUint(stringSubmissionID, 10, 64)
					if err != nil {
						return nil, errors.New("Could not convert id field to an uint")
					}
					verdict := p.Args["verdict"].(string)
					message := p.Args["message"].(string)
					submission, err := FindSubmission(uint(submissionID))
					if err != nil {
						return nil, errors.New(fmt.Sprintf("Submission with id %v not found ", submissionID))
					}
					submission.Status.Verdict = verdict
					submission.Status.Message = message
					if err = UpdateSubmission(submission); err != nil {
						return nil, err
					}
					return submission, nil
				},
			},
		},
	},
)

var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queries,
		Mutation: mutations,
	},
)
