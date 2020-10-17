package ecs

type IEventInit interface {
	Initialize()
}

type IEventPreStart interface {
	PreStart()
}

type IEventStart interface {
	Start()
}

type IEventPostStart interface {
	PostStart()
}

type IEventPreUpdate interface {
	PreUpdate()
}

type IEventUpdate interface {
	Update()
}

type IEventPostUpdate interface {
	PostUpdate()
}

type IEventPreDestroy interface {
	PreDestroy()
}

type IEventDestroy interface {
	Destroy()
}

type IEventPostDestroy interface {
	PostDestroy()
}

//TODO 事件系统，事件需要分发到其他系统，实现系统间互通
