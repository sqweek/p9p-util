/* Plumbsh is an experiment that accepts a snippet of shell
 * script to run for each plumb message received on stdin.
 * Plumbsh sets the variables src, dst, wdir, type, data, and
 * any attr variables to the values specified in the message,
 * for the snippet to reference.
 *
 * By default Plumbsh won't execute anything, it just outputs
 * shell code. The expected invocation is something like:
 *  $ 9p read plumb/event | plumbsh 'echo $src $data' | sh
 */
package main

import (
	"strings"
	"bufio"
	"flag"
	"fmt"
	"os"
)

func shquote(s string) string {
	return "'" + strings.Replace(s, "'", "'\\''", -1) + "'"
}

func getln(in *bufio.Reader) string {
	s, err := in.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimRight(s, "\n")
}

func main() {
	flag.Parse()
	args := flag.Args()
	in := bufio.NewReader(os.Stdin)
	for {
		src := getln(in)
		dst := getln(in) //dst
		wdir := getln(in) //wdir
		msgtype := getln(in)
		attrs := getln(in) //attrs
		lenstr := getln(in)
		var n int
		_, err := fmt.Sscan(lenstr, &n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing data length: %s: %v\n", lenstr, err)
			os.Exit(1)
		}
		buf := make([]byte, n)
		n, err = in.Read(buf)
		if err != nil {
			panic(err)
		}
		if msgtype != "text" {
			fmt.Fprintf(os.Stderr, "unknown message type %s\n", msgtype)
			continue
		}
		data := strings.Replace(string(buf), "\n", "↩", -1)
		for _, attr := range(strings.Fields(attrs)) {
			fmt.Println(shquote(attr))
		}
		fmt.Printf("src=%s\n", shquote(src))
		fmt.Printf("dst=%s\n", shquote(dst))
		fmt.Printf("wdir=%s\n", shquote(wdir))
		fmt.Printf("type=%s\n", shquote(msgtype))
		fmt.Printf("data=%s\n", shquote(data))
		for _, arg := range(args) {
			fmt.Print(arg + " ")
		}
		fmt.Println()
	}
}