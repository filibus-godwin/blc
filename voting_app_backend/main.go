package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {

	sa := option.WithCredentialsFile("./service_account.json")

	app, err := firebase.NewApp(context.Background(), nil, sa)

	if err != nil {
		fmt.Println(err)
	}

	client, err := app.Firestore(context.Background())

	if err != nil {
		fmt.Println(err)
	}

	defer client.Close()

	ids := make([]uuid.UUID, 0, 10)

	router := gin.Default()

	bc := NewBlockchain()

	bc.Print()

	// for a new user to obtain his/her id
	router.GET("/id", func(ctx *gin.Context) {
		id := uuid.New()
		ids = append(ids, id)
		ctx.JSON(http.StatusAccepted, gin.H{
			"id": id.String(),
		})
	})

	router.GET("/candidates", func(ctx *gin.Context) {
		ctx.JSON(http.StatusAccepted, gin.H{
			"candidates": []any{
				map[string]string{
					"candidateName":     "Asiwaju Bola Tinubu",
					"partyId":           "1",
					"candidateImageUrl": "https://i0.wp.com/thedune.ng/wp-content/uploads/2022/06/FB_IMG_1617617998962.jpg?resize=696%2C870&ssl=1",
					"partyName":         "All Progressives Congress (APC)",
				},
				map[string]string{
					"candidateName":     "Atiku Abubakar",
					"partyId":           "2",
					"candidateImageUrl": "https://leadership.ng/wp-content/uploads/2022/11/Atiku-2.webp",
					"partyName":         "People's Democratic Party (PDP)",
				},
				map[string]string{
					"candidateName":     "Peter Obi",
					"partyId":           "3",
					"candidateImageUrl": "https://media.premiumtimesng.com/wp-content/files/2022/10/78f1dc4e-142f-44e4-a328-f15724fe63d4_peter-obi-1140x570.jpg",
					"partyName":         "Labour Party (Lp)",
				},
			},
		})
	})

	// to check if a user has voted
	router.GET("/voted/:id", func(ctx *gin.Context) {
		id, _ := ctx.Params.Get("id")
		col := client.Collection("votes")

		_, err := col.Doc(id).Get(ctx)

		if status.Code(err) == codes.NotFound {
			ctx.JSON(http.StatusAccepted, gin.H{
				"voted": false,
			})
			return
		}

		ctx.JSON(http.StatusAccepted, gin.H{
			"voted": true,
		})

	})

	// for a user to cast a vote
	router.GET("/vote/:voterId/:partyId", func(ctx *gin.Context) {
		voter_id, _ := ctx.Params.Get("voterId")
		party_id, _ := ctx.Params.Get("partyId")

		client.Collection("votes").Doc(voter_id).Set(context.Background(), map[string]string{"voterId": voter_id, "partyId": party_id})

		ctx.JSON(http.StatusAccepted, gin.H{
			"success": true,
		})
	})

	router.Run()
}

type Vote struct {
	From  uuid.UUID `json:"id"`
	Party uuid.UUID `json:"party"`
}

type Block struct {
	Hash      []byte `json:"hash"`
	PrevHash  []byte `json:"prev_hash"`
	Nonce     uint   `json:"nonce"`
	Votes     []Vote `json:"transactions"`
	Timestamp int64  `json:"timestamp"`
}

func (block *Block) HashBlock() {
	mBlock, _ := json.Marshal(block.Votes)
	data := bytes.Join([][]byte{block.PrevHash, []byte(fmt.Sprint(block.Timestamp)), []byte(mBlock)}, []byte{})
	hash := sha256.Sum256(data)
	block.Hash = hash[:]
}

func (block *Block) Print() {
	fmt.Printf("Block hash:         %x\n", block.Hash)
	fmt.Printf("Prev block hash:    %x\n", block.PrevHash)
	fmt.Printf("Nonce:              %d\n", block.Nonce)
	fmt.Printf("Votes:       %+v\n", block.Votes)
	fmt.Printf("Timestamp:          %d\n", block.Timestamp)
}

func NewBlock(nonce uint, prevHash []byte, votes []Vote) *Block {
	block := &Block{
		PrevHash:  prevHash,
		Nonce:     nonce,
		Votes:     votes,
		Timestamp: time.Now().UnixNano(),
	}
	block.HashBlock()
	return block
}

type Blockchain struct {
	votePool []string
	chain    []*Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		votePool: make([]string, 0),
		chain:    make([]*Block, 0),
	}
}

func (bc *Blockchain) Print() {
	for index, block := range bc.chain {
		fmt.Printf("\n%sblock %d%s\n", strings.Repeat("=", 25), index, strings.Repeat("=", 25))
		block.Print()
		fmt.Printf("%s\n", strings.Repeat("*", 55))
	}
}

func (bc *Blockchain) CreateBlock(nonce uint, prevHash []byte, txs []Vote) *Block {
	block := NewBlock(nonce, prevHash, txs)
	bc.chain = append(bc.chain, block)
	return block
}
