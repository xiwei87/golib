package net_server

import (
	"net"
	"testing"
)

func TestTcpConnTableAdd(t *testing.T) {
	var tcpTable TcpConnTable

	/* initialize tcp table */
	tcpTable.Init()

	conn, err := net.Dial("tcp", "svn.baidu.com:http")

	if err != nil {
		t.Error("fail to make connection to svn.baidu.com")
	}

	// add to table
	tcpTable.Add(conn)

	// remove from table
	_, err = tcpTable.Remove(conn)
	if err != nil {
		t.Error("err in tcpTable.Remove()")
	}
}
