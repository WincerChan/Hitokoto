package rpcserver

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	DELIMITER = 18
	STRING    = 0
	ATOM      = 1
	INTEGER   = 2
	FLOAT     = 3
)

var once sync.Once

var rpcConn *RPCConn
var wg sync.WaitGroup

type RPCEncode struct {
	buffer bytes.Buffer
}

type RPCConn struct {
	conn   net.Conn
	signal chan error
	addr   *net.TCPAddr
	tmp    bytes.Buffer
	reader *bufio.Reader
}

func (r *RPCEncode) setLength(value uint32) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, value)
	r.buffer.Write(data)
}
func (r *RPCEncode) getLength(value uint32) []byte {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, value)
	return data
}

func (r *RPCEncode) encodeString(value string) {
	r.buffer.WriteByte(DELIMITER)
	r.setLength(uint32(len(value)))
	r.buffer.WriteByte(STRING)
	r.buffer.Write([]byte(value))
}

func (r *RPCEncode) encodeAtom(value string) {
	r.buffer.WriteByte(DELIMITER)
	r.setLength(uint32(len(value)))
	r.buffer.WriteByte(ATOM)
	r.buffer.Write([]byte(value))
}

func (r *RPCEncode) encodeInteger(value uint32) {
	r.buffer.WriteByte(DELIMITER)
	r.setLength(4)
	r.buffer.WriteByte(INTEGER)
	r.setLength(value)
}

func (r *RPCEncode) encodeFloat(value float64) {
	r.buffer.WriteByte(DELIMITER)
	bits := math.Float64bits(value)
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, bits)
	r.setLength(8)
	r.buffer.WriteByte(FLOAT)
	r.buffer.Write(data)
}

func (r *RPCEncode) execute(conn *Conn) []interface{} {
	r.buffer.Write([]byte("\r\n"))
	line := r.buffer.Bytes()
	size := r.getLength(uint32(len(line)))
	line = append(size, line...)
	conn.WriteOnce(line)
	body := conn.ReadOnce()
	d := &RPCDecode{data: body}
	return d.extract()
}

func main() {
	times, _ := strconv.Atoi(os.Args[1])
	p := NewPool(10, "127.0.0.1:4004", DialTCP)
	wg.Add(times)
	start := time.Now()
	for i := 0; i < times; i++ {
		k := strconv.Itoa(i % 3000)
		go Choke(k, 8, 0.1, p)
	}
	wg.Wait()
	log.Println(time.Since(start))
}
