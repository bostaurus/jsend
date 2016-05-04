package jsend

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/bostaurus/jsend"
	"net/http"
	"os"
	"testing"
)

var urlno int

////////////////////////////////////////////////////////////////////////////////
// Memory buffer compare using Read/Write/WriteFormatted
////////////////////////////////////////////////////////////////////////////////

func memCompare(t *testing.T, jsw *jsend.JSend, formatted bool) {
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)
	var err error
	if formatted {
		err = jsw.WriteFormatted(bw)
	} else {
		err = jsw.Write(bw)
	}
	if err != nil {
		t.Errorf("Error writing JSend message to memory buffer")
		return
	}
	bw.Flush()

	br := bufio.NewReader(&b)
	jsr, err := jsend.Read(br)
	if err != nil {
		t.Errorf("Error reading JSend message from memory buffer")
		return
	}

	if !jsr.IsValid() {
		t.Errorf("Invalid JSend structure read back")
		return
	}

	if jsr.Status != jsw.Status {
		t.Errorf("Mismatch in JSend.Status, sent %v, got %v", jsw.Status, jsr.Status)
		return
	}

	if jsr.Code != jsw.Code {
		t.Errorf("Mismatch in JSend.Code, sent %v, got %v", jsw.Code, jsr.Code)
		return
	}

	if jsr.Message != jsw.Message {
		t.Errorf("Mismatch in JSend.Message, sent %v, got %v", jsw.Message, jsr.Message)
		return
	}

	// Note this won't work with maps, as they are randomly reordered in go
	wdata := fmt.Sprintf("%+v", jsw.Data)
	rdata := fmt.Sprintf("%+v", jsr.Data)
	if wdata != rdata {
		t.Errorf("Mismatch in JSend.Data, sent %s, got %s", wdata, rdata)
	}
}

////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// HTTP compare using Recieve/Send/SendFormatted
////////////////////////////////////////////////////////////////////////////////

func httpCompare(t *testing.T, jsw *jsend.JSend, formatted bool) {
	// Generate dynamic URLs and handlers
	url := fmt.Sprintf("/jsend%d", urlno)
	urlno++
	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) { jsw.Send(w) })

	request, err := http.NewRequest("POST", "http://localhost:8180"+url, nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
		return
	}
	if response.StatusCode != 200 {
		t.Errorf("HTTP request failed: %d", response.StatusCode)
		return
	}

	jsr, err := jsend.Receive(response)
	if err != nil {
		t.Errorf("Error reading JSend message from http Response (%v)", err)
		return
	}

	if !jsr.IsValid() {
		t.Errorf("Invalid JSend structure read back")
		return
	}

	if jsr.Status != jsw.Status {
		t.Errorf("Mismatch in JSend.Status, sent %v, got %v", jsw.Status, jsr.Status)
		return
	}

	if jsr.Code != jsw.Code {
		t.Errorf("Mismatch in JSend.Code, sent %v, got %v", jsw.Code, jsr.Code)
		return
	}

	if jsr.Message != jsw.Message {
		t.Errorf("Mismatch in JSend.Message, sent %v, got %v", jsw.Message, jsr.Message)
		return
	}

	// Note this won't work with maps, as they are randomly reordered in go
	wdata := fmt.Sprintf("%+v", jsw.Data)
	rdata := fmt.Sprintf("%+v", jsr.Data)
	if wdata != rdata {
		t.Errorf("Mismatch in JSend.Data, sent %s, got %s", wdata, rdata)
	}

}

////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Test success message
////////////////////////////////////////////////////////////////////////////////

func TestSuccess(t *testing.T) {
	jss := jsend.Fail("This was a triumph!")
	if !jss.IsValid() {
		t.Errorf("Invalid JSend structure")
	}
	memCompare(t, jss, false)
	memCompare(t, jss, true)
	httpCompare(t, jss, false)
	httpCompare(t, jss, true)
}

////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Test fail message
////////////////////////////////////////////////////////////////////////////////

func TestFail(t *testing.T) {
	jss := jsend.Fail("For the good of all of us, Except the ones who are dead")
	if !jss.IsValid() {
		t.Errorf("Invalid JSend structure")
	}
	memCompare(t, jss, false)
	memCompare(t, jss, true)
	httpCompare(t, jss, false)
	httpCompare(t, jss, true)
}

////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Test error message
////////////////////////////////////////////////////////////////////////////////

func TestError(t *testing.T) {
	jss := jsend.Error("But there's no sense crying over every mistake")
	if !jss.IsValid() {
		t.Errorf("Invalid JSend structure")
	}
	memCompare(t, jss, false)
	memCompare(t, jss, true)
	httpCompare(t, jss, false)
	httpCompare(t, jss, true)
}

////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Test error message with code
////////////////////////////////////////////////////////////////////////////////

func TestErrorCode(t *testing.T) {
	jss := jsend.ErrorCode("You just keep on trying 'til you run out of cake", 1234)
	if !jss.IsValid() {
		t.Errorf("Invalid JSend structure")
	}
	memCompare(t, jss, false)
	memCompare(t, jss, true)
	httpCompare(t, jss, false)
	httpCompare(t, jss, true)
}

////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Test error message with data
////////////////////////////////////////////////////////////////////////////////

func TestErrorData(t *testing.T) {
	jss := jsend.ErrorData("I'm not even angry...", "I'm being so sincere right now")
	if !jss.IsValid() {
		t.Errorf("Invalid JSend structure")
	}
	memCompare(t, jss, false)
	memCompare(t, jss, true)
	httpCompare(t, jss, false)
	httpCompare(t, jss, true)
}

////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Test error message with code and data
////////////////////////////////////////////////////////////////////////////////

func TestErrorCodeWithData(t *testing.T) {
	sarr := [2]string{"And tore me to pieces", "And threw every piece into a fire"}
	jss := jsend.ErrorCodeWithData("Even though you broke my heart, and killed me", 1234, sarr)
	if !jss.IsValid() {
		t.Errorf("Invalid JSend structure")
	}
	memCompare(t, jss, false)
	memCompare(t, jss, true)
	httpCompare(t, jss, false)
	httpCompare(t, jss, true)
}

////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// Startup the HTTP server in the background
////////////////////////////////////////////////////////////////////////////////

func TestMain(m *testing.M) {
	urlno = 1
	go http.ListenAndServe(":8180", nil)
	os.Exit(m.Run())
}

////////////////////////////////////////////////////////////////////////////////
