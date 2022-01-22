package app_test

import (
	"hash/adler32"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashFunction(t *testing.T) {

	key := "random key"
	hashFunc := adler32.New()
	_, err := hashFunc.Write(([]byte)(key))
	assert.NoError(t, err)

	assert.Equal(t, hashFunc.Sum32(), hashFunc.Sum32())

	key2 := "random key2"
	hashFunc2 := adler32.New()
	_, err = hashFunc2.Write(([]byte)(key2))
	assert.NoError(t, err)

	key3 := "random key3"
	hashFunc3 := adler32.New()
	_, err = hashFunc3.Write(([]byte)(key3))
	assert.NoError(t, err)

	assert.NotEqual(t, hashFunc.Sum32(), hashFunc2.Sum32())
	assert.NotEqual(t, hashFunc2.Sum32(), hashFunc3.Sum32())
}
