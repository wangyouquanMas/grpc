package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

/*
   Q1 实现code接口的GobCoder struct 应该包含哪些内容？

	看源码；它实现的功能包含 从连接远端读取到的数据和类型接受 ；并发安全 ；等等
   type Decoder struct {
	mutex        sync.Mutex                              // each item must be received atomically
	r            io.Reader                               // source of the data
	buf          decBuffer                               // buffer for more efficient i/o from r
	wireType     map[typeId]*wireType                    // map from remote ID to local description
	decoderCache map[reflect.Type]map[typeId]**decEngine // cache of compiled engines
	ignorerCache map[typeId]**decEngine                  // ditto for ignored objects
	freeList     *decoderState                           // list of free decoderStates; avoids reallocation
	countBuf     []byte                                  // used for decoding integers while parsing messages
	err          error
}
   其中buffer 源码如下
// decBuffer is an extremely simple, fast implementation of a read-only byte buffer.
// It is initialized by calling Size and then copying the data into the slice returned by Bytes().
type decBuffer struct {
	data   []byte
	offset int // Read offset.
}

    在本项目中使用bufio.Writer 作为buffer， encode，deCode都用到
// Writer implements buffering for an io.Writer object.
// If an error occurs writing to a Writer, no more data will be
// accepted and all subsequent writes, and Flush, will return the error.
// After all data has been written, the client should call the
// Flush method to guarantee all data has been forwarded to
// the underlying io.Writer.
type Writer struct {
	err error
	buf []byte
	n   int
	wr  io.Writer
}

   day1 中参考了源码Decoder结构体，数据流接收
   io.Reader - > io.ReaderWriterCloser
   decBuffer - >   *bufio.Writer

    直接使用了gob 的 Decoder和 Encoder 结构体 【有点冗余？】
*/


type GobCoder struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec *gob.Decoder
	enc *gob.Encoder
}

//确保接口被实现常用的方式。即利用强制类型转换，确保 struct HTTPPool 实现了接口 PeerPicker。这样 IDE 和编译期间就可以检查，而不是等到使用的时候
var _ Codec = (*GobCoder)(nil)

//conn 是由构建函数传入，通常是通过 TCP 或者 Unix 建立 socket 时得到的链接实例
//buf 是为了防止阻塞而创建的带缓冲的 Writer，一般这么做能提升性能。

//先写入到 buffer 中, 然后我们再调用 buffer.Flush() 来将 buffer 中的全部内容写入到 conn 中, 从而优化效率.
// 原本是直接写入到io.writer中，现在是增加了一个 缓冲区， 【 信息 -> buffer -> io.writer 】
// Buffer的诞生可以较少实际物理磁盘的读取
// Buffer在创建的时候就被分配给内存, 这块内存可以被重用 所以减少了动态分配内存空间和回收的次数
//对于读则不需要这方面的考虑, 所以直接在 conn 中读内容即可.
func NewGobCodec(conn io.ReadWriteCloser)Codec{
	buf := bufio.NewWriter(conn)
	return &GobCoder{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}



//实现Codec接口函数
func(c *GobCoder) ReadHeader(h *Header) error{
	return c.dec.Decode(h)
}

func (c *GobCoder) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

// 这种 return , error参与到方法体的用法？
func (c *GobCoder) Write(h *Header, body interface{}) (err error) {

	// buf中内容flush到io.writer
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	if err = c.enc.Encode(h); err != nil {
		log.Println("rpc : gob error encoding header:", err)
		return
	}
	if err = c.enc.Encode(body); err != nil {
		log.Println("rpc: gob error encoding body:", err)
		return
	}
	return
}

func (c *GobCoder) Close() error {
	return c.conn.Close()
}



