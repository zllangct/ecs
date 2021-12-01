package network

type TcpConn struct {
	r chan interface{}
	w chan interface{}
}

func NewTcpConn() *TcpConn {
	return &TcpConn{
		r: make(chan interface{}, 10),
		w: make(chan interface{}, 10),
	}
}

func (t *TcpConn) Write(in interface{}) {
	//ecs.Log.Info("Tcp Send message:", in)
	t.w <- in
}

func (t *TcpConn) Read() interface{} {
	read := <-t.r
	//ecs.Log.Info("Tcp Read message:", read)
	return read
}

var ch chan *TcpConn = make(chan *TcpConn, 10)

type FakeTcpServer struct{}

func Listen() (*FakeTcpServer, error) {
	return &FakeTcpServer{}, nil
}

func Dial(addr string) *TcpConn {
	connSrc := NewTcpConn()
	connDst := NewTcpConn()
	connDst.r, connDst.w = connSrc.w, connSrc.r
	ch <- connDst
	return connSrc
}

func (f *FakeTcpServer) Accept() *TcpConn {
	return <-ch
}
