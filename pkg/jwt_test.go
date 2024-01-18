package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var payload = NewPayload(1, "Admin")

func TestNewPayload(t *testing.T) {

}

func TestGenerateToken(t *testing.T) {
	_, err := payload.GenerateToken()
	assert.NoError(t, err)
}

func TestVerifyToken(t *testing.T) {}
