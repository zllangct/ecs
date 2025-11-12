package ecs

type SystemStandard interface {
	InitReceiver
}

type SystemInitContext struct {
	b    SystemContext
	opts []SystemOption
}

func (s *SystemInitContext) SetOption(opts ...SystemOption) {
	if !s.b.constraint.isValid() {
		return
	}
	s.opts = append(s.opts, opts...)
}
