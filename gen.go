package id

import (
	"errors"
	"log"
	"math/rand/v2"
	"time"
)

// KeyGenerator ensures valid keys for records
type KeyGenerator interface {
	// Validate ensures that the provided key
	// is comprised of a valid date and random
	// random number so that it is globally unique
	IsKeyValid(string) (string, bool)

	// Generate provides a new globally unique
	// URL safe key for a record
	Generate() (string, error)
}

// Key is a unique identifier
// for a record. The Time and Num create
// the id's uniqueness. The random number is
// is needed because records will be created
// quickly enough that they will have the same
// timestamp. The key is a single string value
// which represents the Time and Num in a url
// safe way for external use.
type Key struct {
	// Time is the created time using Unix time
	Time uint64
	// Num is a random number generated from the Time
	Num uint64
	// Value is a URL safe string which
	// represents the Time and Num
	Value string
}

// keyCoder creates a unique string id
// from an array of numbers. It's good
// for link shortening, fast & URL-safe
// ID generation and decoding back into
// numbers for quicker database lookups.
type keyCoder interface {
	// Decode parses out a list of numbers
	// from a Key
	Decode(string) []uint64
	// Encode creates a Key from an array
	// of numbers
	Encode([]uint64) (string, error)
}

// KeyGen creates and validates keys
// to be used as globally unique, URL
// safe identifiers
type KeyGen struct {
	maker keyCoder
}

// NewKeyGen creates an id generator
func NewKeyGen(maker keyCoder) KeyGen {
	return KeyGen{
		maker: maker,
	}
}

func (g KeyGen) IsKeyValid(val string) (string, bool) {
	if val == "" {
		return "", false
	}

	key, err := g.parseKey(val)
	if err != nil {
		log.Printf("[WARN] Error parsing key: %v", err)
		return "", false
	}

	if !isValidUnixDate(key.Time) || key.Num == 0 {
		return "", false
	}

	return key.Value, true
}

// NewId creates an identifier for a record
func (g KeyGen) Generate() (string, error) {
	now := uint64(time.Now().UnixNano())

	// Having a random id along with the time stamp
	// helps ensure that two records recorded at
	// the same nanosecond can still have a unique
	// composite key
	randomId := rand.Uint64()
	source := rand.NewPCG(now, randomId)
	r := rand.New(source)

	// create a key from the time and random id
	key, err := g.maker.Encode([]uint64{now, r.Uint64()})
	if err != nil {
		return "", err
	}

	return key, nil
}

// ParseKey parses out the numbers uses to create the key
func (g KeyGen) parseKey(key string) (Key, error) {
	ids := g.maker.Decode(key)
	if len(ids) != 2 {
		return Key{}, errors.New("incorrect number of ids for key")
	}

	return Key{
		Time:  ids[0],
		Num:   ids[1],
		Value: key,
	}, nil
}

// isValidUnixDate checks if a uint64 value is a valid Unix timestamp.
func isValidUnixDate(timestamp uint64) bool {
	// Define the minimum and maximum valid Unix timestamps in nanoseconds.
	const (
		minUnixNano uint64 = 0                          // Unix epoch start in nanoseconds
		maxUnixNano uint64 = 7258118400 * 1_000_000_000 // Year 2200 in nanoseconds
	)

	// Check if the timestamp is within the valid range.
	if timestamp < minUnixNano || timestamp > maxUnixNano {
		return false
	}

	// Validate by converting to a time.Time object.
	_, err := time.Unix(0, int64(timestamp)).MarshalText()
	return err == nil
}
