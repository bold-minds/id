package id_test

import (
	"errors"
	"testing"

	"github.com/bold-minds/id"
	"github.com/sqids/sqids-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Generate(t *testing.T) {
	keyMaker, err := sqids.New()
	require.NoError(t, err)

	gen := id.NewKeyGen(keyMaker)

	// Act
	key, err := gen.Generate()
	require.NoError(t, err)
	t.Logf("Generated key : %+v", key)

	// Assert
	_, valid := gen.IsKeyValid(key)
	assert.True(t, valid)
}

func Test_Generate_Error(t *testing.T) {
	// Arrange
	fakeErr := errors.New("fake error")

	keys := &KeyMakerMock{}
	keys.On("Encode", mock.Anything).Return("", fakeErr)

	gen := id.NewKeyGen(keys)

	// Act
	key, err := gen.Generate()

	// Assert
	assert.ErrorIs(t, err, fakeErr)
	assert.Empty(t, key)
	keys.AssertExpectations(t)
}

func Test_Generate_NoDups(t *testing.T) {
	keyMaker, err := sqids.New()
	require.NoError(t, err)

	gen := id.NewKeyGen(keyMaker)

	// Act
	keys := map[string]bool{}

	for i := 0; i < 10000; i++ {
		key, err := gen.Generate()
		require.NoError(t, err)

		// Assert
		require.NotContains(t, keys, key)
		_, valid := gen.IsKeyValid(key)
		require.True(t, valid)

		keys[key] = false
	}
}

func Test_IsKeyValid(t *testing.T) {
	// Arrange
	keyMaker, err := sqids.New()
	require.NoError(t, err)

	gen := id.NewKeyGen(keyMaker)
	expected, err := gen.Generate()
	require.NoError(t, err)

	// Act
	actual, valid := gen.IsKeyValid(expected)

	// Assert
	assert.True(t, valid)
	assert.Equal(t, expected, actual)
}

func Test_IsKeyValid_DecodeError(t *testing.T) {
	// Arrange
	key := "abc"

	keys := &KeyMakerMock{}
	keys.On("Decode", key).Return([]uint64{})
	gen := id.NewKeyGen(keys)

	// Act
	_, valid := gen.IsKeyValid(key)

	// Assert
	assert.False(t, valid)
	keys.AssertExpectations(t)
}

func Test_IsKeyValid_EmptyKey(t *testing.T) {
	// Arrange
	keys := &KeyMakerMock{}
	gen := id.NewKeyGen(keys)

	// Act
	key, valid := gen.IsKeyValid("")

	// Assert
	assert.False(t, valid)
	assert.Empty(t, key)
	keys.AssertExpectations(t)
}

func Test_IsKeyValid_InvalidTime(t *testing.T) {
	// Arrange
	key := "abc"

	keys := &KeyMakerMock{}
	keys.On("Decode", key).Return([]uint64{9223372036854775808, 456})
	gen := id.NewKeyGen(keys)

	// Act
	_, valid := gen.IsKeyValid(key)

	// Assert
	assert.False(t, valid)
	keys.AssertExpectations(t)
}

func Test_IsKeyValid_InvalidRandomNumber(t *testing.T) {
	// Arrange
	key := "abc"

	keys := &KeyMakerMock{}
	keys.On("Decode", key).Return([]uint64{1609459200, 0})
	gen := id.NewKeyGen(keys)

	// Act
	_, valid := gen.IsKeyValid(key)

	// Assert
	assert.False(t, valid)
	keys.AssertExpectations(t)
}

type KeyMakerMock struct {
	mock.Mock
}

func (m *KeyMakerMock) Decode(key string) []uint64 {
	args := m.Called(key)
	return args.Get(0).([]uint64)
}

func (m *KeyMakerMock) Encode(ids []uint64) (string, error) {
	args := m.Called(ids)
	return args.String(0), args.Error(1)
}
