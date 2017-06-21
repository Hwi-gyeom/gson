package main

import "os"
import "fmt"
import "strings"
import "testing"
import "bytes"
import "io/ioutil"
import "os/exec"
import "compress/gzip"

var _ = fmt.Sprintf("dummy")
var CMDEXEC = "gson"

var updateref = false

func TestCmdArgs(t *testing.T) {
	testcases := [][]interface{}{
		// json transformations
		[]interface{}{
			[]string{"-inpfile", "example.json", "-json2value"},
			[]byte("Json: \"hello world\"\nValu: hello world\n"),
		},
		[]interface{}{
			[]string{"-inptxt", `"hello world"`, "-json2value"},
			[]byte("Json: \"hello world\"\nValu: hello world\n"),
		},
		[]interface{}{
			[]string{"-inpfile", "example.json", "-json2cbor"},
			[]byte("Json: \"hello world\"\n" +
				"Cbor: [107 104 101 108 108 111 32 119 111 114 108 100]\n" +
				"Cbor: \"khello world\"\n" +
				"Json: \"hello world\"\n"),
		},
		[]interface{}{
			[]string{"-inptxt", `"hello world"`, "-json2cbor"},
			[]byte("Json: \"hello world\"\n" +
				"Cbor: [107 104 101 108 108 111 32 119 111 114 108 100]\n" +
				"Cbor: \"khello world\"\n" +
				"Json: \"hello world\"\n"),
		},
		[]interface{}{
			[]string{"-inpfile", "example.json", "-json2collate"},
			[]byte("Json: \"hello world\"\n" +
				"Coll: \"\\x06hello world\\x00\\x00\"\n" +
				"Coll: [6 104 101 108 108 111 32 119 111 114 108 100 0 0]\n"),
		},
		[]interface{}{
			[]string{"-inptxt", `"hello world"`, "-json2collate"},
			[]byte("Json: \"hello world\"\n" +
				"Coll: \"\\x06hello world\\x00\\x00\"\n" +
				"Coll: [6 104 101 108 108 111 32 119 111 114 108 100 0 0]\n"),
		},
		// json options
		[]interface{}{
			[]string{"-inpfile", "../testdata/typical.json", "-pointers"},
			testdataFile("../testdata/typical_pointers"),
		},
		// cbor transformations
		[]interface{}{
			[]string{"-inpfile", "example.cbor", "-cbor2value"},
			[]byte("Cbor: \"khello world\"\n" +
				"Cbor: [107 104 101 108 108 111 32 119 111 114 108 100]\n" +
				"Valu: hello world\n"),
		},
		[]interface{}{
			[]string{"-inptxt", "khello world", "-cbor2value"},
			[]byte("Cbor: \"khello world\"\n" +
				"Cbor: [107 104 101 108 108 111 32 119 111 114 108 100]\n" +
				"Valu: hello world\n"),
		},
		[]interface{}{
			[]string{"-inpfile", "example.cbor", "-cbor2json"},
			[]byte("Cbor: \"khello world\"\nJson: \"hello world\"\n"),
		},
		[]interface{}{
			[]string{"-inptxt", "khello world", "-cbor2json"},
			[]byte("Cbor: \"khello world\"\nJson: \"hello world\"\n"),
		},
		[]interface{}{
			[]string{"-inpfile", "example.cbor", "-cbor2collate"},
			[]byte("Cbor: \"khello world\"\n" +
				"Coll: \"\\x06hello world\\x00\\x00\"\n" +
				"Coll: [6 104 101 108 108 111 32 119 111 114 108 100 0 0]\n"),
		},
		[]interface{}{
			[]string{"-inptxt", "khello world", "-cbor2collate"},
			[]byte("Cbor: \"khello world\"\n" +
				"Coll: \"\\x06hello world\\x00\\x00\"\n" +
				"Coll: [6 104 101 108 108 111 32 119 111 114 108 100 0 0]\n"),
		},
		[]interface{}{
			[]string{"-inptxt", "khello world", "-cbor2collate"},
			[]byte("Cbor: \"khello world\"\n" +
				"Coll: \"\\x06hello world\\x00\\x00\"\n" +
				"Coll: [6 104 101 108 108 111 32 119 111 114 108 100 0 0]\n"),
		},
		// cbor options
		[]interface{}{
			[]string{"-inptxt", `[10,20]`, "-ct", "stream", "-json2cbor"},
			[]byte(
				"Json: [10,20]\n" +
					"Cbor: [159 251 64 36 0 0 0 0 0 0 251 64 52 0 0 0 0 0 0 255]\n" +
					"Cbor: \"\\x9f\\xfb@$\\x00\\x00\\x00\\x00\\x00\\x00\\xfb@4\\x00\\x00\\x00\\x00\\x00\\x00\\xff\"\n" +
					"Json: [10,20]\n"),
		},
		[]interface{}{
			[]string{"-inptxt", `[10,20]`, "-ct", "lenprefix", "-json2cbor"},
			[]byte(
				"Json: [10,20]\n" +
					"Cbor: [130 251 64 36 0 0 0 0 0 0 251 64 52 0 0 0 0 0 0]\n" +
					"Cbor: \"\\x82\\xfb@$\\x00\\x00\\x00\\x00\\x00\\x00\\xfb@4\\x00\\x00\\x00\\x00\\x00\\x00\"\n" +
					"Json: [10,20]\n"),
		},
		// collate transformations
		[]interface{}{
			[]string{"-inpfile", "example.coll", "-collate2value"},
			[]byte("Coll: \"\\x06hello world\\x00\\x00\"\nValu: hello world\n"),
		},
		[]interface{}{
			[]string{"-inpfile", "example.coll", "-collate2json"},
			[]byte("Coll: \"\\x06hello world\\x00\\x00\"\nJson: \"hello world\"\n"),
		},
		[]interface{}{
			[]string{"-inpfile", "example.coll", "-collate2cbor"},
			[]byte("Coll: \"\\x06hello world\\x00\\x00\"\n" +
				"Cbor: \"khello world\"\n" +
				"Cbor: [107 104 101 108 108 111 32 119 111 114 108 100]\n"),
		},
		// value transformations
		[]interface{}{
			[]string{"-inptxt", `"hello world"`, "-value2json"},
			[]byte("Valu: hello world\nJson: \"hello world\"\n"),
		},
		[]interface{}{
			[]string{"-inptxt", `"hello world"`, "-value2cbor"},
			[]byte("Valu: hello world\nCbor: \"khello world\"\n" +
				"Cbor: [107 104 101 108 108 111 32 119 111 114 108 100]\n"),
		},
		[]interface{}{
			[]string{"-inptxt", `"hello world"`, "-value2collate"},
			[]byte("Valu: hello world\n" +
				"Coll: \"\\x06hello world\\x00\\x00\"\n" +
				"Coll: [6 104 101 108 108 111 32 119 111 114 108 100 0 0]\n"),
		},
	}
	for _, testcase := range testcases {
		args := testcase[0].([]string)
		cmd := exec.Command(CMDEXEC, args...)
		out, _ := cmd.CombinedOutput()
		ref := testcase[1].([]byte)
		if bytes.Compare(out, ref) != 0 {
			t.Logf(strings.Join(args, " "))
			t.Logf("expected %q", ref)
			t.Errorf("got %q", out)
		}
	}
}

func testdataFile(filename string) []byte {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var data []byte
	if strings.HasSuffix(filename, ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(gz)
		if err != nil {
			panic(err)
		}
	} else {
		data, err = ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
	}
	return data
}