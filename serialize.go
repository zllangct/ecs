package ecs

type ICustomSerialize interface {
	Serialize() []byte
	DeSerialize(b []byte)
}