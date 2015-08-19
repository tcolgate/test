package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/kr/text"
)

type node struct {
	label    string
	refine   string
	children []*node
	depth    int
	colour   string
}

type lineData struct {
	indent int
	label  string
	meta   map[string]string
}

func parseLine(l string) lineData {
	var ld lineData

	ls := strings.Split(l, "|")
	lt := strings.Trim(ls[0], " ")
	ld.meta = map[string]string{
		"refine": "OR",
	}

	for _, kv := range ls[1:] {
		ds := strings.SplitN(kv, "=", 2)
		k := ds[0]
		switch {
		case k == "AND", k == "OR": // Refinement
			ld.meta["refine"] = k
		case k == "colour": // Refinement
			ld.meta["colour"] = ds[1]
		case k == "OOAS": // Out of accreditation scope
		case k == "OTHER": // Covered by another attack tree
		case k == "ANALYZE": // Covered by
		}
	}

	ld.indent = len(l) - len(lt)
	ld.label = lt

	return ld
}

func parseAttackTree(r *bufio.Reader, cind int, p *node) (string, error) {
	var err error
	rem := "" // remainder

	for {
		var l string

		if rem == "" {
			b, _, err := r.ReadLine()
			if err != nil {
				break
			}
			l = string(b)
		} else {
			l = rem
			rem = ""
		}

		ld := parseLine(l)
		lt := ld.label
		ind := ld.indent

		switch {
		case lt == "", lt[0] == '#':
			continue
		case ind == cind:
			var peer node
			peer.label = lt
			peer.refine = ld.meta["refine"]
			if col, ok := ld.meta["colour"]; ok {
				peer.colour = col
			}
			p.children = append(p.children, &peer)
		case ind < cind:
			return l, nil
		case ind > cind:
			// We'll add a child to the last node we added
			var child node
			child.label = lt
			child.refine = ld.meta["refine"]
			if col, ok := ld.meta["colour"]; ok {
				child.colour = col
			}
			pnode := p.children[len(p.children)-1]
			pnode.children = append(pnode.children, &child)
			rem, err = parseAttackTree(r, ind, pnode)
			if err != nil {
				return rem, err
			}
		}
	}

	if err == io.EOF {
		err = nil
	}

	return "", nil
}

func prefToString(p []int) string {
	ps := []string{}
	for _, v := range p {
		ps = append(ps, strconv.Itoa(v))
	}
	return "\"" + strings.Join(ps, "_") + "\""
}

func drawnode(f *os.File, r *node, prefix []int) {
	r.depth = len(prefix)
	s := prefToString(prefix)

	kvs := map[string]string{}

	kvs["label"] = fmt.Sprintf("\"%s\"", text.Wrap(r.label, 20))
	if r.colour != "" {
		kvs["style"] = "filled"
		kvs["color"] = r.colour
	}
	var strs []string
	for k, v := range kvs {
		strs = append(strs, k+"="+v)
	}
	fmt.Fprintf(f, "    %s [%s]\n", s, strings.Join(strs, ","))
	if len(r.children) > 0 {
		lastAnd := ""
		for i, c := range r.children {
			newp := append(prefix, i+1)
			if lastAnd != "" {
				fmt.Fprintf(f, "    %s -> %s [label=\"and\",constraint=false]\n", lastAnd, prefToString(newp))
				lastAnd = ""
			}
			if c.refine == "AND" {
				lastAnd = prefToString(newp)
			}
			drawnode(f, c, newp)
			fmt.Fprintf(f, "    %s -> %s\n", s, prefToString(newp))
		}
	} else {
	}
}

func tree2dot(f *os.File, r *node) {
	fmt.Fprintf(f, "digraph G {\n")
	//	fmt.Fprintf(f, "  splines=ortho\n")
	fmt.Fprintf(f, "  edge [dir=none]\n")
	fmt.Fprintf(f, "  node [shape=box]\n")
	for i, t := range r.children {
		fmt.Fprintf(f, "  subgraph  cluster_%d {\n", i)
		drawnode(f, t, []int{0})
		fmt.Fprintf(f, "  }\n")
	}
	fmt.Fprintf(f, "}\n")
}

func main() {
	b := bufio.NewReader(os.Stdin)
	var root node
	_, err := parseAttackTree(b, 0, &root)
	if err != nil {
		log.Fatal(err.Error())
	}

	tree2dot(os.Stdout, &root)
}
