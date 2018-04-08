package main

import (
	"encoding/json"
	"io"
	"os"
)

type Data map[string]string

type Store struct {
	file *os.File
	data Data
}

func OpenStore(name string) (st *Store, err error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	st = &Store{
		file: f,
	}
	err = st.readAll()
	if err != nil {
		st.file.Close()
		return nil, err
	}
	if st.data == nil {
		st.data = make(Data)
	}
	return st, nil
}

func (st *Store) readAll() (err error) {
	_, err = st.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	err = json.NewDecoder(st.file).Decode(&st.data)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func (st *Store) writeAll() (err error) {
	_, err = st.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	err = json.NewEncoder(st.file).Encode(st.data)
	if err != nil {
		return err
	}
	pos, err := st.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	err = st.file.Truncate(pos)
	if err != nil {
		return err
	}
	return nil
}

func (st *Store) Get(key string) (value string) {
	return st.data[key]
}

func (st *Store) Set(key, value string) (err error) {
	if st.data[key] == value {
		return nil
	}
	st.data[key] = value
	return st.writeAll()
}

func (st *Store) Close() (err error) {
	err = st.file.Close()
	st.data = nil
	return err
}
