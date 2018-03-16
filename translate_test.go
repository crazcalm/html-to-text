package htmltotext

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type TestData struct {
	Name string
}

//InputFile -- return content of input file
func (t *TestData) InputFile() (*os.File, error) {
	fileName := fmt.Sprintf("%s.html", t.Name)
	path := filepath.Join("_test_data", fileName)

	return getFile(path)
}

//OutputFile -- return content of output file
func (t *TestData) OutputFile() ([]byte, error) {
	fileName := fmt.Sprintf("%s.txt", t.Name)
	path := filepath.Join("_test_data", fileName)

	return fileContent(path)
}

func getFile(path string) (file *os.File, err error) {
	//Make sure the file exists
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = fmt.Errorf("file: %s does not exist", path)
		return
	}

	//Open the file
	file, err = os.Open(path)
	if err != nil {
		return
	}
	return
}

func fileContent(path string) ([]byte, error) {
	file, err := getFile(path)
	if err != nil {
		return []byte{}, err
	}

	return ioutil.ReadAll(file)
}

func TestTrandlateError(t *testing.T) {
	tests := []struct {
		FileName string
		Error    string
	}{
		{filepath.Join("_test_data", "pError.html"), "ProcessToken had an error: Stack is Empty"},
	}

	for _, test := range tests {
		file, err := getFile(test.FileName)
		if err != nil {
			t.Errorf("%s", err.Error())
		}

		_, _, err = Translate(file)
		if err == nil {
			t.Errorf("Expected error -- %s --, but did not receive it", test.Error)
		}

		if !strings.EqualFold(err.Error(), test.Error) {
			t.Errorf("Expected error: %s, but got %s", test.Error, err.Error())
		}
	}
}

func TestTranslate(t *testing.T) {
	tests := []struct {
		TestData   TestData
		NumOfLinks int
	}{
		{TestData{"3p"}, 0},
		{TestData{"hs"}, 0},
		{TestData{"ol"}, 0},
		{TestData{"ul"}, 0},
		{TestData{"a"}, 1},
		{TestData{"as"}, 4},
		{TestData{"pAnda"}, 1},
		{TestData{"table"}, 0},
		{TestData{"tableWithLinks"}, 2},
		{TestData{"pAndulWitha"}, 3},
		{TestData{"pAndolWitha"}, 3},
		{TestData{"br"}, 0},
		{TestData{"thead"}, 13},
		{TestData{"notRealTables"}, 4},
		{TestData{"div"}, 0},
		{TestData{"divAnda"}, 1},
		{TestData{"olWithInvalidStart"}, 0},
		{TestData{"tableMissingValues"}, 0},
		{TestData{"pAndbr"}, 0},
		{TestData{"ulWithas"}, 9},
		{TestData{"olWithas"}, 9},
	}

	for _, test := range tests {
		inputFile, err := test.TestData.InputFile()
		if err != nil {
			t.Errorf("Error when getting %s.html file: %s", test.TestData.Name, err.Error())
		}

		outputData, err := test.TestData.OutputFile()
		if err != nil {
			t.Errorf("Error when opening %s.txt: %s", test.TestData.Name, err.Error())
		}

		result, links, err := Translate(inputFile)

		//Debug Statments
		//log.Printf("Testing %s.html", test.TestData.Name)
		//log.Printf("result: %s\n\n", result)
		//log.Printf("OutputData: %s", string(outputData))

		if err != nil {
			t.Errorf("Translate('%s.html') threw an error: %s", test.TestData.Name, err.Error())
		}

		if !strings.EqualFold(string(outputData), result) {
			t.Errorf("Expected:\n%s\n\nBut received:\n%s[end]", string(outputData), result)
		}

		if len(links) != test.NumOfLinks {
			t.Errorf("Expeced %d number of links, but got %d", test.NumOfLinks, len(links))
		}
	}
}
