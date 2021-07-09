package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/dvsekhvalnov/jose2go/base64url"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// PkceParams PKCE用のcode_challengeなど
type PkceParams struct {
	CodeChallenge       string `json:"code_challenge,omitempty"`
	CodeChallengeMethod string `json:"code_challenge_method,omitempty"`
	CodeVerifier        string `json:"code_verifier,omitempty"`
	ClientID            string `json:"client_id,omitempty"`
	ResponseType        string `json:"response_type,omitempty"`
}

func GenerateCode() (*PkceParams, error) {
	var pkce PkceParams
	pkce.ResponseType = "code"
	bcodeVerifier, err := randBytes(43)
	if err != nil {
		return nil, err
	}
	pkce.CodeVerifier = string(bcodeVerifier)
	b := sha256.Sum256(bcodeVerifier)
	pkce.CodeChallenge = base64url.Encode(b[:])
	pkce.CodeChallengeMethod = "S256"
	pkce.ClientID = clientID
	return &pkce, nil
}

func randBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	max := new(big.Int)

	max.SetInt64(int64(len(letterBytes)))
	for i := range buf {
		r, err := rand.Int(rand.Reader, max)
		if err != nil {
			return nil, fmt.Errorf("failed to generate random integer: %w", err)
		}

		buf[i] = letterBytes[r.Int64()]
	}

	return buf, nil
}
