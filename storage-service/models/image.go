package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

type ImageMetadata struct {
	ID          string
	Filename    string
	Size        int64
	ContentType string
	Signature   string
	Uploaded    bool
}

// In-memory storage for demo
var Images = make(map[string]*ImageMetadata)
var ImageIDCounter = 1

func VerifySignature(filename string, size int64, contentType, signature string) bool {
	secret := os.Getenv("JWT_SECRET")
	data := fmt.Sprintf("%s:%d:%s", filename, size, contentType)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	expected := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}
