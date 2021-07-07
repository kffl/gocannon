package hist

import (
	"bufio"
	"strconv"
)

func (h *requestHist) writeHistData(w *bufio.Writer) error {

	var lastIdx int64

	for i := h.bins - 1; i >= 0; i-- {
		if h.data[i] > 0 {
			lastIdx = i
			break
		}
	}

	for i := int64(0); i <= lastIdx; i++ {
		_, err := w.WriteString(strconv.FormatInt(h.data[i], 10) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
