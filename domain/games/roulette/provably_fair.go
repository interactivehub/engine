package roulette

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"strconv"

	"github.com/pkg/errors"
)

var (
	ErrClientSeedBlank = errors.New("empty client seed")
	ErrInvalidNonce    = errors.New("invalid nonce")
)

type provablyFair struct {
	serverSeed        []byte
	clientSeed        []byte
	blindedServerSeed []byte
	nonce             uint64
}

func NewProvablyFair(clientSeed []byte, serverSeed []byte) (*provablyFair, error) {
	if len(serverSeed) == 0 {
		var err error
		serverSeed, err = newServerSeed(32)
		if err != nil {
			return nil, err
		}
	}

	blindedSeed := sha256.Sum256(serverSeed)

	return &provablyFair{
		nonce:             0,
		clientSeed:        clientSeed,
		serverSeed:        serverSeed,
		blindedServerSeed: blindedSeed[:],
	}, nil
}

func (f *provablyFair) SetClientSeed(clientSeed []byte) {
	f.clientSeed = clientSeed
}

func (f *provablyFair) Calculate() (uint64, error) {
	hmac, err := f.CalculateHMAC()
	if err != nil {
		return 0, errors.Wrap(err, "failed to calculate outcome")
	}

	stringifiedHMAC := string(hmac)

	var randNum uint64
	for i := 0; i < len(stringifiedHMAC)-5; i++ {
		idx := i * 5
		if len(stringifiedHMAC) < (idx + 5) {
			break
		}

		randNum, err = strconv.ParseUint(stringifiedHMAC[idx:idx+5], 16, 0)
		if err != nil {
			return 0, err
		}

		if randNum <= 999999 {
			break
		}
	}

	if randNum > 999999 {
		return 0, ErrInvalidNonce
	}

	return randNum % 15, nil
}

func (f *provablyFair) CalculateHMAC() ([]byte, error) {
	if len(f.clientSeed) == 0 {
		return nil, ErrClientSeedBlank
	}

	h := hmac.New(sha512.New, f.serverSeed)
	h.Write(append(append(f.clientSeed, '-'), []byte(strconv.FormatUint(f.nonce, 10))...))

	ourHMAC := make([]byte, 128)
	hex.Encode(ourHMAC, h.Sum(nil))

	return ourHMAC, nil
}

func newServerSeed(byteCount int) ([]byte, error) {
	seed := make([]byte, byteCount)

	_, err := rand.Read(seed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate server seed")
	}

	return seed, err
}
