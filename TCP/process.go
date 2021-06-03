package TCP

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)


type result struct {
	v []byte
	e error
}
func (s *Server)readKey(r *bufio.Reader) (string, error) {
	klen, err := readLen(r)
	if err != nil {
		return "", err
	}
	//fmt.Println("key len is", klen)
	k := make([]byte, klen)
	_, err  = io.ReadFull(r, k)
	if err != nil {
		return "", err
	}
	key := string(k)
	addr, ok := s.ShouldProcess(key)
	if !ok {
		return "", errors.New("redirct "+addr)
	}
	return key, nil
}

func (s *Server)readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	//fmt.Println("读取数据流")
	klen, err := readLen(r)
	if err != nil {
		return "", nil, err
	}

	vlen, err := readLen(r)
	if err != nil {
		return "", nil, err
	}
	//fmt.Println("len of key is", klen, "len of value is", vlen)
	k := make([]byte, klen)
	// readfull表示根据所设置的buf的长度进行读取数据
	_, err = io.ReadFull(r, k)
	if err != nil {
		return "", nil, err
	}
	key := string(k)
	addr, ok := s.ShouldProcess(key)
	if !ok {
		return "", nil, errors.New("redirct "+addr)
	}

	v := make([]byte, vlen)
	_, err = io.ReadFull(r, v)
	if err != nil {
		return "", nil, err
	}
	//fmt.Println("len of key is", klen, "len of value is", vlen, "key is", string(k), "values is", string(v))
	return key, v, nil
}

func readLen(r *bufio.Reader) (int, error) {
	// readstring读取数据直到指定结束位置
	tmp, err := r.ReadString(' ')
	if err != nil {
		return 0, err
	}
	// strings.TrimSpace(s string)会返回一个string类型的slice，并将最前面和最后面的ASCII定义的空格去掉，中间的空格不会去掉，如果遇到了\0等其他字符会认为是非空格。
	//fmt.Println(">>>>>>>>>>>>>>", strings.TrimSpace(tmp))
	l, err := strconv.Atoi(strings.TrimSpace(tmp))
	if err != nil {
		return 0, err
	}
	//log.Println("tmp is", tmp, "l is", l)
	return l, nil
}

/*
	sendResponse的作用是返回操作是否成功的信息，格式如下：
	response = error|bytes-array
	成功：返回bytes-array
	失败：-bytes-array
 */
func sendResponse(value []byte, err error, conn net.Conn) error {
	//fmt.Println("send response>>>>>>>>>>>>>>>>")
	if err != nil {
		// 说明set的时候有报错
		errString := err.Error()
		tmp := fmt.Sprintf("-%d", len(errString)) + errString
		fmt.Println("set 操作报错了，错误信息为:", tmp)
		_, e := conn.Write([]byte(tmp))
		return e
	}
// 说明set的时候没有报错
// 不对啊，这个地方value的值为nil，手动写死的，长度自然为0，有何意义
	vlen := fmt.Sprintf("%d", len(value))
	//fmt.Println("set操作成正常，vlen is", vlen)
	//fmt.Println([]byte(vlen), value, string([]byte(vlen)), string(value))
	//data := []byte(vlen)
	//data = append(data, value...)
	//fmt.Println("返回信息为", data, "string is", string(data))

	_, e := conn.Write([]byte(vlen))
	return e
}

func (s *Server) get(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, err := s.readKey(r)
	if err != nil {
		c <- &result{nil, err}
		return
	}
	go func() {
		v, e := s.Get(k)
		c <- &result{v: v, e: e}
	}()

	//return sendResponse(v, e, conn)
}

func (s *Server)set(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	// 解析数据流，这个很关键，
	k, v, e := s.readKeyAndValue(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	//fmt.Println("进行set操作,key is", k, "value is", v)
	// 存放数据
	go func() {
		c <- &result{nil, s.Set(k, v)}
	}()
	// 返回响应
	//return sendResponse(nil, err, conn)
}

func (s *Server) del(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, e := s.readKey(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		c <- &result{nil, s.Del(k)}
	}()
	//return sendResponse(nil, s.Del(k), conn)
}

func (s *Server) process(conn net.Conn)  {
	//fmt.Println("process.................")
	r := bufio.NewReader(conn)
	resultCh := make(chan chan *result, 5000)
	defer close(resultCh)
	go reply(conn, resultCh)
	num := 1
	for {
		fmt.Println("开始处理第", num, "个请求")
		// readbyte好像是获取net.conn数据流的第一个字符
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("close connection due to error:", e)
			}
			return
		}
		//fmt.Println("op is", string(op))
		if op == 'S' {
			s.set(resultCh, r)
		}else if op == 'G' {
			s.get(resultCh, r)
		} else if op == 'D' {
			s.del(resultCh, r)
		} else {
			log.Println("close connection due to invalid operation:", op)
			return
		}
		if e != nil {
			log.Println("close connection due to error:", e)
			return
		}
		num += 1
 	}
}

func reply(conn net.Conn, ch chan chan *result) {
	defer conn.Close()
	for {
		c, open := <- ch
		if !open {
			// 说明resultCh关闭了
			return
		}
		r := <- c
		e := sendResponse(r.v, r.e, conn)
		if e != nil {
			log.Println("reply function close connection due to err:", e)
			return
		}
	}
}