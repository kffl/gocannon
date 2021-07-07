package reqlog

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	header = "code;start;end;\n"
)

func TestToCSV(t *testing.T) {
	cases := []request{
		{200, 123001, 123007},
		{404, 220000, 220001},
		{-1, 123000, 124000},
	}

	results := []string{
		"200;1;7;\n",
		"404;97000;97001;\n",
		"-1;0;1000;\n",
	}

	for i := range cases {
		assert.Equal(t, results[i], cases[i].toCSV(123000))
	}
}

func TestWriteRawReqDataEmpty(t *testing.T) {
	reqs := newRequests(20, 1000)
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	flattened := reqs.flatten()

	flattened.writeRawReqData(w, 0)

	assert.Equal(
		t,
		header,
		b.String(),
		"only the header should be written if there are no requests",
	)
}

func TestWriteRawReqDataPopulated(t *testing.T) {
	reqs := newRequests(2, 1000)
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	reqs[0] = append(reqs[0], request{200, 123001, 123007})
	reqs[1] = append(reqs[1], request{-1, 10000, 11000})

	flattened := reqs.flatten()

	flattened.writeRawReqData(w, 0)

	assert.Equal(
		t,
		header+"200;123001;123007;\n"+"-1;10000;11000;\n",
		b.String(),
		"full CSV output should be written",
	)
}
