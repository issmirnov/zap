package main

import (
	"bytes"
	"container/list"
	"fmt"
	"strings"

	"github.com/Jeffail/gabs"
)

const (
	expand = iota
	query
	port
)

// tokenize takes a string delimited by slashes and splits it up into tokens
// returns a linked list.
func tokenize(path string) *list.List {

	// Creates the list.
	l := list.New()
	for _, tok := range strings.Split(path, "/") {
		l.PushBack(tok)
	}
	return l
}

// Returns string to write to result, boolean flag indicating whether to advance
// token, and error if needed.
// The option to advance the token is needed when we want to suppress the slash separator.
func getPrefix(c *gabs.Container) (string, int, error) {
	d := c.Path(expandKey).Data()
	if d != nil {
		s, oks := d.(string)
		i, oki := d.(float64)
		if oks {
			return s, expand, nil
		} else if oki {
			return fmt.Sprintf("%.f", i), expand, nil
		}
		return "", 0, fmt.Errorf("unexpected type of expansion value, got %T instead of int or string", d)
	}
	q := c.Path(queryKey).Data()
	if q != nil {
		if s, ok := q.(string); ok {
			return s, query, nil
		}
		return "", 0, fmt.Errorf("casting query key to string failed for %T:%v", q, q)
	}

	p := c.Path(portKey).Data()
	if p != nil {
		if s, ok := p.(float64); ok {
			return fmt.Sprintf(":%.f", s), port, nil
		}
		return "", 0, fmt.Errorf("casting port key to float64 failed for %T:%v", p, p)
	}

	return "", 0, fmt.Errorf("error in config, no key matching 'expand', 'query' or 'port' in %s", c.String())
}

// expandPath takes a config, list of tokens (parsed from request) and the results buffer
// At each level of recursion, it matches the token to the action described in the config, and writes it
// to the result buffer. There is special care needed to handle slashes correctly, which makes this function
// quite nontrivial. Tests are crucial to ensure correctness.
func expandPath(c *gabs.Container, token *list.Element, res *bytes.Buffer) {
	if token == nil {
		return
	}
	children, _ := c.ChildrenMap()
	child, ok := children[token.Value.(string)]
	if ok {
		p, action, err := getPrefix(child)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		switch action {
		case expand: // Generic case, write slash followed by expanded token.
			res.WriteString("/")
			res.WriteString(p)

		case query: // Write a slash + query string expansion, then perform token skipahead in order to have correct slashes.
			res.WriteString("/")
			res.WriteString(p)
			if token.Next() != nil {
				res.WriteString(token.Next().Value.(string))
				token = token.Next()
			}

		case port: // A little bit of a special case - unlike "expand", we don't want a leading slash.
			res.WriteString(p)

		default:
			panic("Programmer error, this should never happen.")
		}
		expandPath(child, token.Next(), res)
		return
	}

	// if tokens left over, append the rest
	for e := token; e != nil; e = e.Next() {
		res.WriteString("/")
		res.WriteString(e.Value.(string))
	}
}
