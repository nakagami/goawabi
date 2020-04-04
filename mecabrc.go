/*****************************************************************************
MIT License

Copyright (c) 2020 Hajime Nakagami

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*****************************************************************************/

package goawabi

import (
	"bufio"
	"errors"
	"os"
	"regexp"
)

func find_mecabrc() (path string, err error) {
	pathes := []string{"/usr/local/etc/mecabrc", "/etc/mecabrc"}
	for _, s := range pathes {
		_, e := os.Stat(s)
		if !os.IsNotExist(e) {
			path = s
			return path, err
		}
	}

	err = errors.New("Can't find mecabrc")
	return path, err
}

func get_mecabrc_map(path string) (mecabrc_map map[string]string, err error) {
	mecabrc_map = make(map[string]string)

	if path == "" {
		path, err = find_mecabrc()
		if err != nil {
			return mecabrc_map, err
		}
	}

	fp, err := os.Open(path)
	if err != nil {
		return mecabrc_map, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	re := regexp.MustCompile(`^(\S+)\s*=\s*(\S+)`)
	for scanner.Scan() {
		group := re.FindAllStringSubmatch(scanner.Text(), -1)
		if len(group) == 1 { // at most 1
			mecabrc_map[group[0][1]] = group[0][2]
		}
	}
	return mecabrc_map, err
}
