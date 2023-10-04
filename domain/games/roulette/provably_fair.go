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

type ProvablyFair struct {
	ServerSeed        []byte
	ClientSeed        []byte
	BlindedServerSeed []byte
	Nonce             uint64
}

func NewProvablyFair(clientSeed []byte, serverSeed []byte) (*ProvablyFair, error) {
	if len(serverSeed) == 0 {
		var err error
		serverSeed, err = newServerSeed(32)
		if err != nil {
			return nil, err
		}
	}

	blindedSeed := sha256.Sum256(serverSeed)

	return &ProvablyFair{
		Nonce:             0,
		ClientSeed:        clientSeed,
		ServerSeed:        serverSeed,
		BlindedServerSeed: blindedSeed[:],
	}, nil
}

func (f *ProvablyFair) SetClientSeed(clientSeed []byte) {
	f.ClientSeed = clientSeed
}

func (f *ProvablyFair) Calculate() (uint64, error) {
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

func (f *ProvablyFair) CalculateHMAC() ([]byte, error) {
	if len(f.ClientSeed) == 0 {
		return nil, ErrClientSeedBlank
	}

	h := hmac.New(sha512.New, f.ServerSeed)
	h.Write(append(append(f.ClientSeed, '-'), []byte(strconv.FormatUint(f.Nonce, 10))...))

	ourHMAC := make([]byte, 128)
	hex.Encode(ourHMAC, h.Sum(nil))

	return ourHMAC, nil
}

func (f *ProvablyFair) StringServerSeed() string {
	return hex.EncodeToString(f.ServerSeed)
}

func (f *ProvablyFair) StringClientSeed() string {
	return hex.EncodeToString(f.ClientSeed)
}

func (f *ProvablyFair) StringBlindedServerSeed() string {
	return hex.EncodeToString(f.BlindedServerSeed)
}

func newServerSeed(byteCount int) ([]byte, error) {
	seed := make([]byte, byteCount)

	_, err := rand.Read(seed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate server seed")
	}

	return seed, err
}
