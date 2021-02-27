package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/api/sessions", Auth)
	r.POST("/api/graphql", Query)
	if err := r.Run("0.0.0.0:3000"); err != nil {
		log.Fatal(err)
	}
}

func Auth(ctx *gin.Context) {
	ctx.String(200, "ok")
}

type ProblemBody struct {
	Title string `json:"title"`
}
type Problem struct {
	Code string      `json:"code"`
	Body ProblemBody `json:"body"`
}

type Team struct {
	Number int `json:"number"`
}

type Penalty struct {
	ID        string    `json:"id"`
	Problem   Problem   `json:"problem"`
	Team      Team      `json:"team"`
	CreatedAt time.Time `json:"createdAt"`
}

type ResponseData struct {
	Penalties []Penalty `json:"penalties"`
}
type Response struct {
	Data ResponseData `json:"data"`
}

func Query(ctx *gin.Context) {
	t1, _ := time.Parse(time.RFC3339Nano, "2021-02-20T11:11:12+09:00")
	t2, _ := time.Parse(time.RFC3339Nano, "2021-02-24T11:11:12+09:00")
	ctx.JSON(200, Response{
		Data: ResponseData{
			Penalties: []Penalty{
				{
					ID: "test",
					Problem: Problem{
						Code: "test",
						Body: ProblemBody{
							Title: "test problem",
						},
					},
					Team: Team{
						Number: 2,
					},
					CreatedAt: t1,
				},
				{
					ID: "test2",
					Problem: Problem{
						Code: "hoge",
						Body: ProblemBody{
							Title: "hoge problem",
						},
					},
					Team: Team{
						Number: 10,
					},
					CreatedAt: t2,
				},
			},
		},
	})
}
