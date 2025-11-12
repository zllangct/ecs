package ecs

type LightSystem func(*SystemContext, Event) error

func (l LightSystem) Update(ctx *SystemContext, event Event) error {
	return TryAndReport(func() error {
		return l(ctx, event)
	})
}
