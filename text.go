package main

import (
	"bytes"
	"container/list"
	"fmt"
	"strings"

	"github.com/Jeffail/gabs"
)

// tokenize takes a string delimited by slahes and splits it up into tokens
// returns a linkedlist.
func tokenize(path string) *list.List {

	// Creates the list.
	l := list.New()
	for _, tok := range strings.Split(path, "/") {
		l.PushBack(tok)
	}
	return l
}

// Returns string to write to result, boolean flag indicating whther to advance
// token, and error if needed.
func getPrefix(c *gabs.Container) (string, bool, error) {
	d := c.Path(expandKey).Data()
	if d != nil {
		s, oks := d.(string)
		i, oki := d.(float64)
		if oks {
			return s, false, nil
		} else if oki {
			return fmt.Sprintf("%.f", i), false, nil
		}
		return "", false, fmt.Errorf("unexpected type of expansion value, got %T instead of int or string", d)
	}
	q := c.Path(queryKey).Data()
	if q != nil {
		if s, ok := q.(string); ok {
			return s, true, nil
		}
		return "", false, fmt.Errorf("Casting query key to string failed for %T:%v", q, q)
	}
	return "", false, fmt.Errorf("error in config, no expand or query key in %s", c.String())
}
func expand(c *gabs.Container, token *list.Element, res *bytes.Buffer) {
	// base case
	if token == nil {
		return
	}
	res.WriteString("/")
	children, _ := c.ChildrenMap()
	for key, child := range children {
		if key == token.Value {
			p, skip, err := getPrefix(child)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			res.WriteString(p)
			if skip {
				token = token.Next()
				if token == nil {
					return
				}
				res.WriteString(token.Value.(string))
			}
			expand(child, token.Next(), res)
			return
		}
	}
	// handle base case if no keys matched
	res.WriteString(token.Value.(string))

	// if tokens left over, append the rest
	for e := token.Next(); e != nil; e = e.Next() {
		res.WriteString("/")
		res.WriteString(e.Value.(string))
	}
}
