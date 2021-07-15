package reqlog

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCSV(t *testing.T) {
	cases := []request{
		{200, 123001, 123007},
		{404, 220000, 220001},
		{0, 123000, 124000},
	}

	results := []string{
		"200;1;7;0;\n",
		"404;97000;97001;1;\n",
		"0;0;1000;2;\n",
	}

	for i := range cases {
		assert.Equal(t, results[i], cases[i].toCSV(123000, i))
	}
}

func TestWriteRawReqDataEmpty(t *testing.T) {
	reqs := make(flatRequestLog, 0)
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	reqs.writeRawReqData(w, 0, 0)

	assert.Equal(
		t,
		"",
		b.String(),
		"no output should be written if the request log for a given connection is empty",
	)
}

func TestWriteRawReqDataPopulated(t *testing.T) {
	reqs := make(flatRequestLog, 0)
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	reqs = append(reqs, request{200, 123001, 123007})
	reqs = append(reqs, request{0, 10000, 11000})

	reqs.writeRawReqData(w, 0, 1)

	assert.Equal(
		t,
		"200;123001;123007;1;\n"+"0;10000;11000;1;\n",
		b.String(),
		"full CSV output of a given connection should be written",
	)
}
