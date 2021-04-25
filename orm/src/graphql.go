package orm

import (
	"errors"
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
				Type: graphql.Int,
			},
			"memoryLimit": &graphql.Field{
				Type: graphql.Int,
			},
			"inputs": &graphql.Field{
				Type: graphql.NewList(testFile),
			},
			"outputs": &graphql.Field{
				Type: graphql.NewList(testFile),
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
		},
	},
)

var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queries,
	},
)
