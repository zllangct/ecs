package ecs

// TODO 世界的序列化、反序列化
type ICustomSerialize interface {
	Serialize() []byte
	DeSerialize(b []byte)
}
