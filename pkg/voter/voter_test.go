package voter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vsychov/go-rating-stars/tests/mocks"
	"testing"
	"time"
)

func TestBasicDataFetch(t *testing.T) {
	storage := new(mocks.StorageInterface)
	storage.On("Get", "test_total").Return(float64(10), nil)
	storage.On("Get", "test_rating").Return(float64(3.232), nil)
	storage.On("Get", "f8e20ff64b6919aacf96b94c5b0d1b2e925d99b541bf90a568b876ec").Return(float64(0), fmt.Errorf("not found"))

	voter := Create(storage)
	results := voter.Results("test", "::1")
	assert.Equal(t, results.TotalVotes, 10)
	assert.Equal(t, results.Rating, 3.232)
	assert.True(t, results.AllowedToVote)
	storage.AssertExpectations(t)
}

func TestLockedVoting(t *testing.T) {
	storage := new(mocks.StorageInterface)
	storage.On("Get", "test_total").Return(float64(10), nil)
	storage.On("Get", "test_rating").Return(float64(3.232), nil)
	storage.On("Get", "f8e20ff64b6919aacf96b94c5b0d1b2e925d99b541bf90a568b876ec").Return(float64(1), nil)

	voter := Create(storage)
	results := voter.Results("test", "::1")
	assert.Equal(t, results.TotalVotes, 10)
	assert.Equal(t, results.Rating, 3.232)
	assert.False(t, results.AllowedToVote)
	storage.AssertExpectations(t)
}

func TestFirstVote(t *testing.T) {
	storage := new(mocks.StorageInterface)
	storage.On("Get", "test_total").Return(float64(0), nil)
	storage.On("Get", "test_rating").Return(float64(0), nil)
	storage.On("SetNX", "f8e20ff64b6919aacf96b94c5b0d1b2e925d99b541bf90a568b876ec", float64(1), time.Duration(86400000000000)).Return(nil)
	storage.On("Incr", "test_rating", float64(5)).Return(float64(5), nil)
	storage.On("Incr", "test_total", float64(1)).Return(float64(1), nil)

	voter := Create(storage)
	results := voter.Vote("test", "::1", 5)
	assert.Equal(t, results.TotalVotes, 1)
	assert.Equal(t, results.Rating, float64(5))
	assert.False(t, results.AllowedToVote)
	storage.AssertExpectations(t)
}

func TestKeepRatingVote(t *testing.T) {
	storage := new(mocks.StorageInterface)
	storage.On("Get", "test_total").Return(float64(1), nil)
	storage.On("Get", "test_rating").Return(float64(5), nil)
	storage.On("SetNX", "f8e20ff64b6919aacf96b94c5b0d1b2e925d99b541bf90a568b876ec", float64(1), time.Duration(86400000000000)).Return(nil)
	storage.On("Incr", "test_rating", float64(0)).Return(float64(5), nil)
	storage.On("Incr", "test_total", float64(1)).Return(float64(2), nil)

	voter := Create(storage)
	results := voter.Vote("test", "::1", 5)
	assert.Equal(t, results.TotalVotes, 2)
	assert.Equal(t, results.Rating, float64(5))
	assert.False(t, results.AllowedToVote)
	storage.AssertExpectations(t)
}

func TestDecrRatingVote(t *testing.T) {
	storage := new(mocks.StorageInterface)
	storage.On("Get", "test_total").Return(float64(1), nil)
	storage.On("Get", "test_rating").Return(float64(5), nil)
	storage.On("SetNX", "f8e20ff64b6919aacf96b94c5b0d1b2e925d99b541bf90a568b876ec", float64(1), time.Duration(86400000000000)).Return(nil)
	storage.On("Incr", "test_rating", float64(-2)).Return(float64(3), nil)
	storage.On("Incr", "test_total", float64(1)).Return(float64(2), nil)

	voter := Create(storage)
	results := voter.Vote("test", "::1", 1)
	assert.Equal(t, results.TotalVotes, 2)
	assert.Equal(t, results.Rating, float64(3))
	assert.False(t, results.AllowedToVote)
	storage.AssertExpectations(t)
}

func TestIncrRatingVote(t *testing.T) {
	storage := new(mocks.StorageInterface)
	storage.On("Get", "test_total").Return(float64(1), nil)
	storage.On("Get", "test_rating").Return(float64(1), nil)
	storage.On("SetNX", "f8e20ff64b6919aacf96b94c5b0d1b2e925d99b541bf90a568b876ec", float64(1), time.Duration(86400000000000)).Return(nil)
	storage.On("Incr", "test_rating", float64(2)).Return(float64(3), nil)
	storage.On("Incr", "test_total", float64(1)).Return(float64(2), nil)

	voter := Create(storage)
	results := voter.Vote("test", "::1", 5)
	assert.Equal(t, results.TotalVotes, 2)
	assert.Equal(t, results.Rating, float64(3))
	assert.False(t, results.AllowedToVote)
	storage.AssertExpectations(t)
}

func TestLockedVote(t *testing.T) {
	storage := new(mocks.StorageInterface)
	storage.On("Get", "test_total").Return(float64(10), nil)
	storage.On("Get", "test_rating").Return(float64(3.232), nil)
	storage.On("SetNX", "f8e20ff64b6919aacf96b94c5b0d1b2e925d99b541bf90a568b876ec", float64(1), time.Duration(86400000000000)).Return(fmt.Errorf("error"))

	voter := Create(storage)
	results := voter.Vote("test", "::1", 5)
	assert.Equal(t, results.TotalVotes, 10)
	assert.Equal(t, results.Rating, 3.232)
	assert.False(t, results.AllowedToVote)
	storage.AssertExpectations(t)
}

func TestResultsNoRating(t *testing.T) {
	storage := new(mocks.StorageInterface)
	storage.On("Get", "test_total").Return(float64(0), fmt.Errorf("not found"))
	storage.On("Get", "test_rating").Return(float64(0), fmt.Errorf("not found"))
	storage.On("Get", "f8e20ff64b6919aacf96b94c5b0d1b2e925d99b541bf90a568b876ec").Return(float64(0), fmt.Errorf("not found"))

	voter := Create(storage)
	results := voter.Results("test", "::1")
	assert.Equal(t, results.TotalVotes, 0)
	assert.Equal(t, results.Rating, float64(0))
	assert.True(t, results.AllowedToVote)
	storage.AssertExpectations(t)
}

// TestStorageIncrFailed if storage failed - return 0 values to client
func TestStorageIncrFailed(t *testing.T) {
	storage := new(mocks.StorageInterface)
	storage.On("Get", "test_total").Return(float64(0), nil)
	storage.On("Get", "test_rating").Return(float64(0), nil)
	storage.On("SetNX", "f8e20ff64b6919aacf96b94c5b0d1b2e925d99b541bf90a568b876ec", float64(1), time.Duration(86400000000000)).Return(nil)
	storage.On("Incr", "test_rating", float64(5)).Return(float64(0), fmt.Errorf("error incr"))
	storage.On("Incr", "test_total", float64(1)).Return(float64(0), fmt.Errorf("error incr"))

	voter := Create(storage)
	results := voter.Vote("test", "::1", 5)
	assert.Equal(t, results.TotalVotes, 0)
	assert.Equal(t, results.Rating, float64(0))
	assert.False(t, results.AllowedToVote)
	storage.AssertExpectations(t)
}
