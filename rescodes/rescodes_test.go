package rescodes

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordRequest(t *testing.T) {
	r := NewRescodes()

	r.RecordRequest(200)
	r.RecordRequest(200)
	r.RecordRequest(404)
	r.RecordRequest(200)
	r.RecordRequest(429)

	assert.Equal(t, int64(3), r[200])
	assert.Equal(t, int64(1), r[404])
	assert.Equal(t, int64(1), r[429])
	assert.Equal(t, int64(0), r[100])
}

func TestRecordRequestThreadSafe(t *testing.T) {
	r := NewRescodes()

	var wg sync.WaitGroup

	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 20000; j++ {
				r.RecordRequestThreadSafe(200)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(200000), r[200])
}
