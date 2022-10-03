package cache

import (
	"context"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type RateLimitTestSuite struct {
	suite.Suite
	mr *miniredis.Miniredis
}

func (s *RateLimitTestSuite) SetupSuite() {
	mr, err := miniredis.Run()
	s.Require().NoError(err)
	s.mr = mr
}

func (s *RateLimitTestSuite) TearDownSuite() {
	s.mr.Close()
}

func (s *RateLimitTestSuite) Test5TPSForCoreBanking() {
	c := redis.NewClient(&redis.Options{
		Addr: s.mr.Addr(),
	})
	rl := New(c)
	ctx, cfn := context.WithTimeout(context.Background(), 40*time.Second)
	defer cfn()
	retryAfterEquals1ns := func(d time.Duration) bool {
		return -1 == d
	}
	retryAfterNotEquals1ns := func(d time.Duration) bool {
		return -1 != d
	}
	fixtures := []struct {
		wantAllowed    bool
		wantRetryAfter func(time.Duration) bool
	}{
		{true, retryAfterEquals1ns},
		{true, retryAfterEquals1ns},
		{true, retryAfterEquals1ns},
		{true, retryAfterEquals1ns},
		{true, retryAfterEquals1ns},
		{false, retryAfterNotEquals1ns},
		{false, retryAfterNotEquals1ns},
		{false, retryAfterNotEquals1ns},
	}
	for idx, f := range fixtures {
		res, err := rl.GetRateLimit(ctx, "idx", 5, time.Second)
		s.Require().NoError(err)
		s.EqualValuesf(f.wantAllowed, res.Allowed, "idx: %d", idx)
		s.Truef(f.wantRetryAfter(res.RetryAfter), "idx: %d, retry after: %s", idx, res.RetryAfter)
	}
}

func TestRateLimitTestSuite(t *testing.T) {
	suite.Run(t, new(RateLimitTestSuite))
}
