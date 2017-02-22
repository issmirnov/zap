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

func expand(c *gabs.Container, token *list.Element, res *bytes.Buffer) {
    // base case
    if token == nil {
        return
    }
    res.WriteString("/")
    children, _ := c.ChildrenMap()
    for key, child := range children {
        if key == token.Value {
            d := child.Path(expandKey).Data()
            if d == nil {
                fmt.Printf("error in config, no expand key in %s\n", c.String())
            }
            s, oks := d.(string)
            i, oki := d.(float64)
            if oks {
                res.WriteString(s)
            } else if oki {
                res.WriteString(fmt.Sprintf("%.f", i))
            } else {
                fmt.Printf("Unexpected type of expansion value, got %T instead of int or string.", d)
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
