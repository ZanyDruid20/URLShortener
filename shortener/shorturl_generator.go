package shortener
import (
	"crypto/sha256"
	"fmt"
	"github.com/itchyny/base58-go"
	"math/big"
	"os"
)

// sha256Of returns the SHA-256 hash of the input string as a byte slice.
func sha256Of(input string) []byte {
	// Create a new SHA-256 hash instance.
	algo := sha256.New()
	// Write the input string as bytes to the hash instance.
	algo.Write([]byte(input))
	// Compute and return the final hash as a byte slice.
	return algo.Sum(nil)
}

// base58Encoded encodes the given byte slice using Base58 encoding and returns the result as a string.
func base58Encoded(bytes []byte) string {
	// Use the Bitcoin Base58 encoding scheme.
	encoding := base58.BitcoinEncoding
	// Encode the byte slice using Base58.
	encoded, err := encoding.Encode(bytes)
	// If there is an error during encoding, print the error and exit the program.
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// Return the encoded result as a string.
	return string(encoded)
}

// GenerateShortLink generates a short link for the given URL.
// GenerateShortLink generates a short link for the given URL and user UUID.
func GenerateShortLink(initialLink string, userId string) string {
	// Concatenate the original URL and the user UUID, then compute the SHA-256 hash.
	urlHashBytes := sha256Of(initialLink + userId)
	// Convert the hash bytes to a big integer, then get its uint64 representation.
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	// Format the uint64 number as a decimal string, then encode it using Base58.
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	// Return the first 8 characters of the Base58-encoded string as the short link.
	return finalString[:8]
}
