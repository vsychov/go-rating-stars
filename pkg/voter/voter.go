package voter

import (
	"fmt"
	"github.com/vsychov/go-rating-stars/pkg/storage"
	"golang.org/x/crypto/sha3"
	"log"
	"time"
)

// Voter structure
type Voter struct {
	storage storage.StorageInterface
}

// VoteResults return results
type VoteResults struct {
	ResourceId    string
	Rating        float64
	TotalVotes    int
	AllowedToVote bool
}

// Create create voter instance
func Create(storage storage.StorageInterface) Voter {
	return Voter{
		storage: storage,
	}
}

// Vote for some resource
//500 votes, avg - 3.9
//+1 vote, 4 score
//500 * 3.9 = (1950 + 4) / 501 = 3.90019
//1 vote, avg - 5
//+1 vote, 2 score
//1 * 5 = (5 + 1) / 2 = 3
func (v Voter) Vote(resourceId string, sourceId string, vote float64) VoteResults {
	ratingKey := getRatingKey(resourceId)
	totalVotesKey := getTotalVotesKey(resourceId)
	totalRating := v.readOrGetDefault(ratingKey, 0)
	totalVotes := v.readOrGetDefault(totalVotesKey, 0)

	voteResult := ((totalVotes * totalRating) + vote) / (totalVotes + 1)
	voteCost := voteResult - totalRating

	err := v.lockVotingSource(computeLockKey(resourceId, sourceId))
	if err != nil {
		return VoteResults{
			TotalVotes:    int(v.readOrGetDefault(totalVotesKey, 0)),
			Rating:        v.readOrGetDefault(ratingKey, 0),
			AllowedToVote: false,
		}
	}

	return VoteResults{
		Rating:        v.incrKey(ratingKey, voteCost),
		TotalVotes:    int(v.incrKey(totalVotesKey, 1)),
		AllowedToVote: false,
	}
}

// Results response
func (v Voter) Results(resourceId string, sourceId string) VoteResults {
	totalVotesKey := getTotalVotesKey(resourceId)
	ratingKey := getRatingKey(resourceId)

	return VoteResults{
		TotalVotes:    int(v.readOrGetDefault(totalVotesKey, 0)),
		Rating:        v.readOrGetDefault(ratingKey, 0),
		AllowedToVote: !v.isVoteLocked(computeLockKey(resourceId, sourceId)),
	}
}

func (v Voter) lockVotingSource(lockKey string) error {
	return v.storage.SetNX(lockKey, 1, time.Second*86400)
}

func (v Voter) isVoteLocked(lockKey string) bool {
	_, err := v.storage.Get(lockKey)

	return err == nil
}

func (v Voter) incrKey(key string, val float64) float64 {
	val, err := v.storage.Incr(key, val)
	if err != nil {
		log.Printf("Error updating key: %s, error: %s", key, err)
		val = 0
	}

	return val
}

func (v Voter) readOrGetDefault(key string, def float64) float64 {
	val, err := v.storage.Get(key)
	if err != nil {
		log.Printf("Error read %s key, err: %s", key, err)
		val = def
	}

	return val
}

func computeLockKey(resourceId string, sourceId string) string {
	return fmt.Sprintf("%x", sha3.Sum224([]byte(resourceId+sourceId)))
}

func getTotalVotesKey(resourceId string) string {
	return resourceId + "_total"
}

func getRatingKey(resourceId string) string {
	return resourceId + "_rating"
}
