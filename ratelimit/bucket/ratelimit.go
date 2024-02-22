package ratelimit // The algorithm that this implementation uses does computational work
import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Bucket struct {
	Options Options

	RedisKey string

	clock Clock

	//每次申请的数量
	applyNum int64

	// startTime holds the moment when the bucket was
	// first created and ticks began.
	startTime time.Time

	// capacity holds the overall capacity of the bucket.
	capacity int64

	// quantum holds how many tokens are added on
	// each tick.
	quantum int64

	// fillInterval holds the interval between each tick.
	fillInterval time.Duration

	// mu guards the fields below it.
	mu sync.Mutex

	// availableTokens holds the number of available
	// tokens as of the associated latestTick.
	// It will be negative when there are consumers
	// waiting for tokens.
	availableTokens int64

	// latestTick holds the latest tick for which
	// we know the number of tokens in the bucket.
	latestTick int64
}

type Buckets func(*Bucket)

func SetApplyNum(applyNum int64) Buckets {
	return func(bucket *Bucket) {
		bucket.applyNum = applyNum
	}
}

func SetCapacity(capacity int64) Buckets {
	return func(bucket *Bucket) {
		bucket.capacity = capacity
	}
}

func SetAvailableTokens(availableTokens int64) Buckets {
	return func(bucket *Bucket) {
		bucket.availableTokens = availableTokens
	}
}

func SetOptions(o Options) Buckets {
	return func(bucket *Bucket) {
		bucket.Options = o
	}
}

// NewBucketWithQuantumAndClock is like NewBucketWithQuantum, but
// also has a clock argument that allows clients to fake the passing
// of time. If clock is nil, the system clock will be used.
func NewBucket(opt ...Buckets) *Bucket {
	clock := realClock{}

	bucket := &Bucket{
		Options:      DefaultOptions,
		clock:        clock,
		startTime:    clock.Now(),
		latestTick:   0,
		fillInterval: time.Second,
	}

	for _, v := range opt {
		v(bucket)
	}

	return bucket
}

// Take takes count tokens from the bucket without blocking. It returns
// the time that the caller should wait until the tokens are actually
// available.
//
// Note that if the request is irrevocable - there is no way to return
// tokens to the bucket once this method commits us to taking them.
func (tb *Bucket) Take(ctx context.Context, count int64) (bool, error) {
	ok, err := tb.take(ctx, tb.clock.Now(), count)
	return ok, err
}

// Capacity returns the capacity that the bucket was created with.
func (tb *Bucket) Capacity() int64 {
	return tb.capacity
}

// take is the internal version of Take - it takes the current time as
// an argument to enable easy testing.
func (tb *Bucket) take(ctx context.Context, now time.Time, count int64) (bool, error) {
	if count <= 0 {
		return true, nil
	}

	avail := tb.availableTokens - count
	if avail >= 0 {
		atomic.AddInt64(&tb.availableTokens, -count)
		return true, nil
	}

	tick := tb.currentTick(now)

	ok, err := tb.Options.Apply(ctx, tb.applyNum, tick)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	atomic.AddInt64(&tb.availableTokens, -count+tb.applyNum)

	return true, nil
}

//返回到现在的秒数
func (tb *Bucket) currentTick(now time.Time) int64 {
	return int64(now.Sub(tb.startTime) / 1e9)
}

// Clock represents the passage of time in a way that
// can be faked out for tests.
type Clock interface {
	// Now returns the current time.
	Now() time.Time
	// Sleep sleeps for at least the given duration.
	Sleep(d time.Duration)
}

// realClock implements Clock in terms of standard time functions.
type realClock struct{}

// Now implements Clock.Now by calling time.Now.
func (realClock) Now() time.Time {
	return time.Now()
}

// Now implements Clock.Sleep by calling time.Sleep.
func (realClock) Sleep(d time.Duration) {
	time.Sleep(d)
}
