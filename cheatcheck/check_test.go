package cheatcheck

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func read(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

func readDataFromFile(path string) []int {
	data := read(path)
	steps := strings.Split(data, "\n")

	stepInts := make([]int, 0, len(steps))
	for i := 0; i < len(steps); i++ {
		if steps[i] != "" {
			j, err := strconv.Atoi(steps[i])
			if err == nil {
				stepInts = append(stepInts, j)
			}
		}
	}

	return stepInts
}

func TestCheatCheck(t *testing.T) {

	err := filepath.Walk("./test_data", func(path string, info os.FileInfo, err error) error {
		fmt.Printf("\n\nfile: %s\n", path)

		data := readDataFromFile(path)

		cheatIndexes, err1 := StepsCheatCheck(data)
		if err1 != nil {
			t.Errorf("error: %s", err1.Error())
		}

		for _, item := range cheatIndexes {
			fmt.Printf("[%d, %d], ", item[0], item[1])
		}
		return nil
	})

	if err != nil {
		t.Error(err.Error())
	}
}
