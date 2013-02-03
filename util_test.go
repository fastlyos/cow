package main

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func TestReadLine(t *testing.T) {
	testData := []struct {
		text  string
		lines []string
	}{
		{"nihao\nwelcome\r", []string{"nihao", "welcome"}},
		{"hello\r\nworld\n", []string{"hello", "world"}},
		{"hello\r\r\nworld\r\r", []string{"hello", "world"}},
		{"hello", []string{"hello"}},
		{"hello\r\r", []string{"hello"}},
		{"\r\n", []string{""}},
		{"\r", []string{""}},
		{"\n", []string{""}},
		{"\r\r\n", []string{""}},
	}
	for _, td := range testData {
		raw := bytes.NewBufferString(td.text)
		rd := bufio.NewReader(raw)
		for i, line := range td.lines {
			l, err := ReadLine(rd)
			if err != nil {
				t.Fatalf("%d read error %v got: %s\ntext: %s\n", i+1, err, l, td.text)
			}
			if line != l {
				t.Fatalf("%d read got: %s should be: %s\ntext: %s\n", i+1, l, line, td.text)
			}
		}
		_, err := ReadLine(rd)
		if err != io.EOF {
			t.Error("ReadLine past end should return EOF")
		}
	}
}

func TestIsDigit(t *testing.T) {
	for i := 0; i < 10; i++ {
		digit := '0' + byte(i)
		letter := 'a' + byte(i)

		if IsDigit(digit) != true {
			t.Errorf("%c should return true", digit)
		}

		if IsDigit(letter) == true {
			t.Errorf("%c should return false", letter)
		}
	}
}

func TestCopyN(t *testing.T) {
	testStr := "hello world"
	src := bytes.NewBufferString(testStr)
	dst := new(bytes.Buffer)
	buf := make([]byte, 5)

	copyN(src, dst, nil, len(testStr), buf, nil, nil)
	if dst.String() != "hello world" {
		t.Error("copy without pre and end failed, got:", dst.String())
	}

	src.Reset()
	dst.Reset()
	src.WriteString(testStr)
	copyN(src, dst, nil, len(testStr), buf, []byte("by cyf "), nil)
	if dst.String() != "by cyf hello world" {
		t.Error("copy with pre no end failed, got:", dst.String())
	}

	src.Reset()
	dst.Reset()
	src.WriteString(testStr)
	copyN(src, dst, nil, len(testStr), buf, []byte("by cyf "), []byte(" welcome"))
	if dst.String() != "by cyf hello world welcome" {
		t.Error("copy with both pre and end failed, got:", dst.String())
	}

	src.Reset()
	dst.Reset()
	src.WriteString(testStr)
	copyN(src, dst, nil, len(testStr), buf, []byte("pre longer then buffer "), []byte(" welcome"))
	if dst.String() != "pre longer then buffer hello world welcome" {
		t.Error("copy with long pre failed, got:", dst.String())
	}

	src.Reset()
	dst.Reset()
	testStr = "34"
	src.WriteString(testStr)
	copyN(src, dst, nil, len(testStr), buf, []byte("12"), []byte(" welcome"))
	if dst.String() != "1234 welcome" {
		t.Error("copy len(pre)+size<bufLen failed, got:", dst.String())
	}

	src.Reset()
	dst.Reset()
	testStr = "2"
	src.WriteString(testStr)
	copyN(src, dst, nil, len(testStr), buf, []byte("1"), []byte("34"))
	if dst.String() != "1234" {
		t.Error("copy len(pre)+size+len(end)<bufLen failed, got:", dst.String())
	}
}

func TestIsFileExists(t *testing.T) {
	exists, err := isFileExists("testdata")
	if err == nil {
		t.Error("should return error is path is directory")
	}
	if exists {
		t.Error("directory should return false")
	}

	exists, err = isFileExists("testdata/none")
	if exists {
		t.Error("BOOM! You've found a non-existing file!")
	}
	if err != nil {
		t.Error("Not existing file should just return false, on error")
	}

	exists, err = isFileExists("testdata/file")
	if !exists {
		t.Error("testdata/file exists, but returns false")
	}
	if err != nil {
		t.Error("Why error for existing file?")
	}
}

func TestNewNbitIPv4Mask(t *testing.T) {
	mask := []byte(NewNbitIPv4Mask(32))
	for i := 0; i < 4; i++ {
		if mask[i] != 0xff {
			t.Error("NewNbitIPv4Mask with 32 error")
		}
	}
	mask = []byte(NewNbitIPv4Mask(5))
	if mask[0] != 0xf8 || mask[1] != 0 || mask[2] != 0 {
		t.Error("NewNbitIPv4Mask with 5 error:", mask)
	}
	mask = []byte(NewNbitIPv4Mask(9))
	if mask[0] != 0xff || mask[1] != 0x80 || mask[2] != 0 {
		t.Error("NewNbitIPv4Mask with 9 error:", mask)
	}
	mask = []byte(NewNbitIPv4Mask(23))
	if mask[0] != 0xff || mask[1] != 0xff || mask[2] != 0xfe || mask[3] != 0 {
		t.Error("NewNbitIPv4Mask with 23 error:", mask)
	}
	mask = []byte(NewNbitIPv4Mask(28))
	if mask[0] != 0xff || mask[1] != 0xff || mask[2] != 0xff || mask[3] != 0xf0 {
		t.Error("NewNbitIPv4Mask with 28 error:", mask)
	}
}

func TestHost2Domain(t *testing.T) {
	var testData = []struct {
		host   string
		domain string
	}{
		{"www.google.com", "google.com"},
		{"google.com", "google.com"},
		{"com.cn", "com.cn"},
		{"sina.com.cn", "sina.com.cn"},
		{"www.bbc.co.uk", "bbc.co.uk"},
		{"apple.com.cn", "apple.com.cn"},
		{"simplehost", ""},
		{"192.168.1.1", ""},
	}

	for _, td := range testData {
		dm := host2Domain(td.host)
		if dm != td.domain {
			t.Errorf("%s got domain %v should be %v", td.host, dm, td.domain)
		}
	}
}
