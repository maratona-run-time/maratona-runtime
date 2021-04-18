package main

import (
	"github.com/graphql-go/graphql"

	"github.com/maratona-run-time/Maratona-Runtime/orm/src"
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
					return orm.FindChallenge(submission.ChallengeID)
				},
			},
		},
	},
)
