package main

import (
	"io/ioutil"
	"os"

	"gordiff/delta"
	"gordiff/signature"
)

const defaultWindow = 3
const defaultSigFilename = "signature.yml"
const defaultDeltaFilename = "delta.yml"

func main() {
	// Making signature for the first file
	// Representing the array of checksums
	// For each chunk. Chunk size is configurable.
	f1, err := os.Open("old.txt")
	defer f1.Close()
	sb := signature.NewSignatureBuilder(f1, defaultWindow)
	signature := sb.BuildSignature()

	// Try marshal and save
	mrhs, _ := signature.MarshalYAML()
	err = ioutil.WriteFile(defaultSigFilename, mrhs, 0644)
	if err != nil {
		panic(err)
	}

	// Building the delta between 2nd file and the signature.
	// Ideally sig should be unpacked. Window size for the rolling hash
	// Listed in signature of the old file
	// REquires sig validation in future
	f2, err := os.Open("new.txt")
	defer f1.Close()
	db := delta.NewDeltaBuilder(signature, f2)
	d := db.BuildDelta()

	mrhs, _ = d.MarshalYAML()
	err = ioutil.WriteFile(defaultDeltaFilename, mrhs, 0644)
	if err != nil {
		panic(err)
	}
}
