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

// Returns error if neither value present.
func getPrefix(c *gabs.Container) (string, error) {
	d := c.Path(expandKey).Data()
	if d != nil {
		s, oks := d.(string)
		i, oki := d.(float64)
		if oks {
			return fmt.Sprintf("%s", s), nil
		} else if oki {
			return fmt.Sprintf("%.f", i), nil
		}
		return "", fmt.Errorf("unexpected type of expansion value, got %T instead of int or string", d)
	}
	q := c.Path(queryKey).Data()
	if q != nil {
		if s, ok := q.(string); ok {
			return s, nil
		}
		return "", fmt.Errorf("Casting query key to string failed for %T:%v", q, q)
	}
	return "", fmt.Errorf("error in config, no expand or query key in %s", c.String())
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
			p, err := getPrefix(child)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			res.WriteString(p)
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
