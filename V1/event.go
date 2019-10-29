package main

type IEventInit interface {
	Initialise()
}

type IEventStart interface {
	Start()
}

type IEventUpdate interface {
	Update()
}

type IEventDestroy interface {
	Destroy()
}

//TODO 事件系统，事件需要分发到其他系统，实现系统间互通