package stats

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func (r request) toCSV(start int64) string {
	return fmt.Sprintf("%v;%v;%v;\n", r.code, r.start-start, r.end-start)
}

func (reqs flattenedRequests) writeRawReqData(w *bufio.Writer, start int64) error {
	_, err := fmt.Fprintf(w, "code;start;end;\n")
	if err != nil {
		return err
	}

	for i := 0; i < len(reqs); i++ {
		_, err := fmt.Fprintf(w, "%v", reqs[i].toCSV(start))
		if err != nil {
			return err
		}
	}

	w.Flush()

	return nil
}

func (reqs flattenedRequests) saveCSV(start int64, outputFile string) error {
	f, err := os.Create(outputFile)
	if err != nil {
		return errors.New("error creating output file")
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	err = reqs.writeRawReqData(w, start)

	return err
}
