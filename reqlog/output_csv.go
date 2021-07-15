package reqlog

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func (r request) toCSV(start int64, connection int) string {
	return fmt.Sprintf("%v;%v;%v;%v;\n", r.code, r.start-start, r.end-start, connection)
}

func (reqs *flatRequestLog) writeRawReqData(w *bufio.Writer, start int64, connection int) error {

	for i := 0; i < len(*reqs); i++ {
		_, err := fmt.Fprintf(w, "%v", (*reqs)[i].toCSV(start, connection))
		if err != nil {
			return err
		}
	}

	err := w.Flush()

	return err
}

func (reqs *requestLog) saveCSV(start int64, outputFile string) error {
	f, err := os.Create(outputFile)
	if err != nil {
		return errors.New("error creating output file")
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	_, err = fmt.Fprintf(w, "code;start;end;connection;\n")
	if err != nil {
		return err
	}

	for i, req := range *reqs {
		flatReq := flatRequestLog(req)
		err = (&flatReq).writeRawReqData(w, start, i)
		if err != nil {
			return err
		}
	}

	err = w.Flush()

	return err
}
