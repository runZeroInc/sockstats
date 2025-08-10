/**
 * Copyright (c) 2022, Xerra Earth Observation Institute.
 * Copyright (c) 2025, Simeon Miteff.
 *
 * See LICENSE.TXT in the root directory of this source tree.
 */

package exporter

import (
	"fmt"
	"github.com/higebu/netfd"
	"net"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simeonmiteff/go-tcpinfo/pkg/linux"
)

type info struct {
	description *prometheus.Desc
	supplier    func(tcpInfo *linux.TCPInfo, labelValues []string) prometheus.Metric
}

type connEntry struct {
	fd     int
	labels []string
}

type TCPInfoCollector struct {
	conns  map[net.Conn]connEntry
	mu     sync.Mutex
	logger func(error)
	infos  []info
}

func (t *TCPInfoCollector) Describe(descs chan<- *prometheus.Desc) {
	for _, info := range t.infos {
		descs <- info.description
	}
}

func (t *TCPInfoCollector) Collect(metrics chan<- prometheus.Metric) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for conn, entry := range t.conns {
		tcpInfo, err := linux.GetTCPInfo(entry.fd)
		if err != nil {
			t.logger(fmt.Errorf("error getting connection tcpinfo (removing conn %v -> %v): %w", conn.LocalAddr(), conn.RemoteAddr(), err))

			delete(t.conns, conn)
			continue
		}

		for _, info := range t.infos {
			metrics <- info.supplier(tcpInfo, entry.labels)
		}
	}
}

func (t *TCPInfoCollector) Add(conn net.Conn, labels []string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.conns[conn] = connEntry{
		fd:     netfd.GetFdFromConn(conn),
		labels: labels,
	}
}

func (t *TCPInfoCollector) Remove(conn net.Conn) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.conns, conn)
}

func NewTCPInfoCollector(
	prefix string,
	connectionLabels []string, // connectionLabels are known up front for the collector and values are provided when adding a connection.
	constLabels prometheus.Labels, // constLabels is meant for labels with values that are constant for the whole process.
	errorLoggingCallback func(error),
) *TCPInfoCollector {
	t := TCPInfoCollector{ //nolint:exhaustivestruct
		conns:  make(map[net.Conn]connEntry),
		logger: errorLoggingCallback,
	}
	t.addMetrics(prefix, connectionLabels, constLabels)
	return &t
}
