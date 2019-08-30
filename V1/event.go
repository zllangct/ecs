package main

type IEventStart interface {
	Start()
}

type IEventUpdate interface {
	Update()
}

type IEventDestroy interface {
	Destroy()
}

