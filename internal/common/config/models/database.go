package models

import "time"

type DataBasePool struct {
	MaxOpenConn     int
	MaxIdleConn     int
	ConnMaxLifeTime time.Duration
	ConnMaxIdleTime time.Duration
}

type DataBasePoolBuilder struct {
	*DataBasePool
	err error
}

func Builder() *DataBasePoolBuilder {
	b := new(DataBasePoolBuilder)
	/*
		some default setting
	*/
	return b
}

func (b *DataBasePoolBuilder) MaxOpenConn(conn int) *DataBasePoolBuilder {
	if b.err != nil {
		return b
	}
	// verify conn
	b.DataBasePool.MaxOpenConn = conn
	return b
}
func (b *DataBasePoolBuilder) MaxIdleConn(conn int) *DataBasePoolBuilder {
	if b.err != nil {
		return b
	}
	b.DataBasePool.MaxIdleConn = conn
	return b
}
func (b *DataBasePoolBuilder) ConnMaxLifeTime(d time.Duration) *DataBasePoolBuilder {
	if b.err != nil {
		return b
	}
	b.DataBasePool.ConnMaxLifeTime = d
	return b
}
func (b *DataBasePoolBuilder) ConnMaxIdleTime(d time.Duration) *DataBasePoolBuilder {
	if b.err != nil {
		return b
	}
	b.DataBasePool.ConnMaxIdleTime = d
	return b
}

func (b *DataBasePoolBuilder) Build() *DataBasePool {
	return b.DataBasePool
}
