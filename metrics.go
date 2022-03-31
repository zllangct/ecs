package ecs

import "time"

type Metrics struct {
	enable  bool
	isPrint bool
	m       map[string]*MetricReporter
}

func (m *Metrics) NewReporter(name string) *MetricReporter {
	mr := &MetricReporter{
		name:          name,
		sampleElapsed: []reporterStep{},
	}
	mr.metrics = m
	return mr
}

func (m *Metrics) Print() {
	if !m.enable {
		return
	}
	for _, reporter := range m.m {
		reporter.Print()
	}
}

func NewMetrics(enable bool, print bool) *Metrics {
	return &Metrics{
		enable:  enable,
		isPrint: print,
		m:       make(map[string]*MetricReporter),
	}
}

type reporterStep struct {
	name    string
	elapsed time.Duration
}

type MetricReporter struct {
	name          string
	metrics       *Metrics
	start         time.Time
	last          time.Time
	sampleElapsed []reporterStep
	elapsedTotal  time.Duration
}

func (m *MetricReporter) Start() {
	if !m.metrics.enable {
		return
	}
	m.start = time.Now()
	m.last = m.start
}

func (m *MetricReporter) Sample(name string) {
	if !m.metrics.enable {
		return
	}
	now := time.Now()
	m.elapsedTotal = now.Sub(m.start)
	m.sampleElapsed = append(m.sampleElapsed, reporterStep{
		name:    name,
		elapsed: now.Sub(m.last),
	})
	m.last = now
}

func (m *MetricReporter) Stop() {
	if !m.metrics.enable {
		return
	}
	now := time.Now()
	m.elapsedTotal = now.Sub(m.start)
	if m.metrics != nil {
		m.metrics.m[m.name] = m
	}
}

func (m *MetricReporter) Print(force ...bool) {
	if !m.metrics.enable || !m.metrics.isPrint {
		if len(force) > 0 && force[0] {
		} else {
			return
		}
	}
	Log.Infof("%s: cost: %+v\n", m.name, m.elapsedTotal)
	for _, r := range m.sampleElapsed {
		Log.Infof("    ├─%20s: %+v\n", r.name, r.elapsed)
	}
}
