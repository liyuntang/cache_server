package HTTP

import (
	"bytes"
	"fmt"
	"net/http"
)

type rebalanceHandler struct {
	*Server
}

func (h *rebalanceHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	go h.rebalance()
}

func (h *rebalanceHandler)rebalance()  {
	s := h.NewScanner()
	defer s.Close()
	c := &http.Client{}
	for s.Scan() {
		k := s.Key()
		n, ok := h.ShouldProcess(k)
		if !ok {
			url := fmt.Sprintf("http://%s:12345/cache/%s", n, k)
			r, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(s.Value()))
			c.Do(r)
			h.Del(k)
		}
	}
}

func (s *Server)rebalanceHandler() http.Handler {
	return &rebalanceHandler{s}
}







