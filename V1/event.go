package main

type IEventInit interface {
	Init()
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

