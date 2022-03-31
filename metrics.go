package ecs

import "time"

var metrics = NewMetrics()

var M = metrics

type Metrics struct {
	m map[string]*MetricReporter
}

func (m *Metrics) Print() {
	for s, reporter := range m.m {
		Log.Infof("%s: \n", s)
		for s2, durations := range reporter.sampleElapsed {
			Log.Infof("    %s: %+v\n", s2, durations)
		}
	}
}

func NewMetrics() *Metrics {
	return &Metrics{
		m: make(map[string]*MetricReporter),
	}
}

type MetricReporter struct {
	name          string
	Metrics       *Metrics
	start         time.Time
	last          time.Time
	sampleElapsed map[string][]time.Duration
	elapsedTotal  time.Duration
}

func NewMetricsReporter(name string, metrics ...*Metrics) *MetricReporter {
	mr := &MetricReporter{
		name:          name,
		sampleElapsed: map[string][]time.Duration{},
	}
	if len(metrics) > 0 {
		mr.Metrics = metrics[0]
	}
	return mr
}

func (m *MetricReporter) Start() {
	m.start = time.Now()
	m.last = m.start
}

func (m *MetricReporter) Sample(name string) {
	now := time.Now()
	m.elapsedTotal = now.Sub(m.start)
	elapsed := now.Sub(m.last)
	if list, ok := m.sampleElapsed[name]; ok {
		list = append(list, elapsed)
		m.sampleElapsed[name] = list
	} else {
		m.sampleElapsed[name] = []time.Duration{elapsed}
	}
	m.last = now
}

func (m *MetricReporter) Stop() {
	now := time.Now()
	m.elapsedTotal = now.Sub(m.start)
	if m.Metrics != nil {
		m.Metrics.m[m.name] = m
	}
}

func (m *MetricReporter) Print() {
	Log.Infof("%s: cost: %+v\n", m.name, m.elapsedTotal)
	for s, durations := range m.sampleElapsed {
		Log.Infof("    ├─%s: %+v\n", s, durations)
	}
}
