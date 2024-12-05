package main

import (
	"context"
	"fmt"
	pgembed "pg-vector-db/internal/pg_embed"
)

func main() {
	// pgembed.Connect()
	// pgembed.PgVectorDbEmbed()
	input := []string{
		"The dog is barking",
		"The cat is purring",
		"The bear is growling",
		"The lion is roaring",
		"The tiger is snarling",
		"The elephant is trumpeting",
		"Animals are amazing",
		"Sky is blue",
	}
	embedResults, err := pgembed.FetchEmbeddings(input, pgembed.EmbedderUrl, "nomic-embed-text")
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(embedResults)
	// find dimension of embeddings
	fmt.Printf("Embedding dimension: %d\n", len(embedResults[0]))

	ctx := context.Background()

	conn, err := pgembed.CreateConnection(ctx)
	if err != nil {
		fmt.Println(err)
	}

	table := "documents"
	err = pgembed.CreateVectorTable(ctx, conn, table)
	if err != nil {
		fmt.Println(err)
	}

	err = pgembed.LoadVectorData(ctx, input, embedResults, conn)
	if err != nil {
		fmt.Println(err)
	}

	pgembed.RunQuery(ctx, conn, []string{"which animal is roaring"}, table)
}
