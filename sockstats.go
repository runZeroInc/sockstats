package sockstats

import (
	"encoding/json"
	"net"
	"time"

	"github.com/runZeroInc/sockstats/simeonmiteff/tcpinfo"
)

const (
	SockStatsOpen  = 0
	SockStatsClose = 1
)

var StateMap = map[int]string{
	SockStatsOpen:  "open",
	SockStatsClose: "close",
}

type ReportStatsFn func(tic *Conn, state int)

type Conn struct {
	net.Conn
	reportStats  func(*Conn, int)
	OpenedAt     int64
	ClosedAt     int64
	FirstReadAt  int64
	FirstWriteAt int64
	SentBytes    int64
	RecvBytes    int64
	RecvErr      error
	SentErr      error
	Attempts     int
	Details      map[string]any
}

func WrapConn(ncon net.Conn, reportStatsFn ReportStatsFn) net.Conn {
	w := &Conn{
		Conn:        ncon,
		reportStats: reportStatsFn,
		OpenedAt:    time.Now().UnixNano(),
		Details:     make(map[string]any),
	}
	w.gatherAndReport(SockStatsOpen)
	return w
}

func (w *Conn) gatherAndReport(state int) {
	if w.reportStats == nil {
		return
	}

	tcpConn, ok := w.Conn.(*net.TCPConn)
	if !ok {
		return
	}

	rawConn, err := tcpConn.SyscallConn()
	if err != nil {
		return
	}

	var tcpInfo *tcpinfo.SysInfo
	if err := rawConn.Control(func(fd uintptr) {
		tcpInfo, err = tcpinfo.GetTCPInfo(int(fd))
	}); err != nil {
		return
	}
	b, err := json.Marshal(tcpInfo)
	if err != nil {
		return
	}
	if err := json.Unmarshal(b, &w.Details); err != nil {
		return
	}
	w.reportStats(w, state)
}

// SetConnectionAttempts stores the number of attempts that were needed to open this connection.
func (w *Conn) SetConnectionAttempts(attempts int) {
	w.Attempts = attempts
}

// Close invokes the reportWrapper with a close event before closing the connection.
func (w *Conn) Close() error {
	w.ClosedAt = time.Now().UnixNano()
	w.reportStats(w, SockStatsClose)
	return w.Conn.Close()
}

// Read wraps the underlying Read method and tracks the data
func (w *Conn) Read(b []byte) (int, error) {
	n, err := w.Conn.Read(b)
	if err == nil && w.RecvBytes == 0 && n > 0 {
		// Track the timestamp of the first successful read
		w.FirstReadAt = time.Now().UnixNano()
	}
	w.RecvBytes += int64(n)
	if err, ok := err.(net.Error); ok && !err.Timeout() {
		w.RecvErr = err
	}
	return n, err
}

// Write wraps the underlying Write method and tracks the data
func (w *Conn) Write(b []byte) (int, error) {
	n, err := w.Conn.Write(b)
	if err == nil && w.SentBytes == 0 && n > 0 {
		// Track the timestamp of the first successful write
		w.FirstWriteAt = time.Now().UnixNano()
	}
	w.SentBytes += int64(n)
	w.SentErr = err
	if err, ok := err.(net.Error); ok && !err.Timeout() {
		w.SentErr = err
	}
	return n, err
}
