package main

type SystemOrder uint32
const(
	SYSTEM_ORDER_DEFAULT SystemOrder = iota
	SYSTEM_ORDER_START_PRE
	SYSTEM_ORDER_UPDATE_PRE
	SYSTEM_ORDER_DESTROY_PRE
	SYSTEM_ORDER_DESTROY_POST
)

type SystemGroup struct {

}