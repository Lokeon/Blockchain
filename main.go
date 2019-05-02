package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Block container
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

// Message container
type Message struct {
	BPM int
}

// Blockchain itself
var Blockchain []Block

func calculateHash(block Block) string {
	registro := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	hash := sha256.New()
	hash.Write([]byte(registro))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(antBlock Block, BPM int) (Block, error) {
	var newBlock Block
	tiem := time.Now()

	newBlock.Index = antBlock.Index + 1
	newBlock.Timestamp = tiem.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = antBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func isBlockValid(newBlock, antBlock Block) bool {
	if antBlock.Index+1 != newBlock.Index {
		return false
	}

	if antBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func cambiarChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, cod int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(cod)
	w.Write(response)
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", " ")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var men Message

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&men); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)

		return
	}

	defer r.Body.Close()
}
func crearMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")

	return muxRouter
}

func run() error {
	mux := crearMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Mirando ...", os.Getenv("ADDR"))
	conf := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := conf.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Print("HOOOOLA")
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		tiem := time.Now()
		genesisBlock := Block{0, tiem.String(), 0, "", ""}
		spew.Dump(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)
	}()

	log.Fatal(run())
}
