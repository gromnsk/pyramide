// To Do
// -------
// Make delete a little cleaner
// Add traversal

package main

import (
	"github.com/gromnsk/tree"
	"github.com/takama/router"
	"strconv"
	// "net/http"
	"fmt"
	// "net/url"
	// "encoding/json"
)

type Result struct {
	Success bool

}

func main() {
	root := new(tree.Tree)
	data := tree.Data{Id: 1}
	node := root.Insert(data)
	
	fmt.Println("Started!")
	r := router.New()
	
	//put node somewhere from root
	r.PUT("/node/:id", func(c *router.Control) {
		c.UseTimer()
		Id, err := strconv.Atoi(c.Get(":id"))
		if err != nil {
			// TODO catch error
		}
		data := tree.Data{Id: Id}
		node.Insert(data)
        c.Code(200).Body(data)
    })

	// PUT node to ReferrerId
    r.PUT("/node/:referrerId/:id", func(c *router.Control) {
		c.UseTimer()
		Id, err := strconv.Atoi(c.Get(":id"))
		referrerId, refErr := strconv.Atoi(c.Get(":referrerId"))
		if err != nil || refErr != nil {
			// TODO catch error
		}

		referrerNode := node.Search(referrerId)

		data := tree.Data{Id: Id}
		referrerNode.Insert(data)
        c.Code(200).Body(data)
    })

    //Get node by ID 
    r.GET("/node/:id", func(c *router.Control) {
		c.UseTimer()
    	Id, err := strconv.Atoi(c.Get(":id"))
    	if err != nil {
			// TODO catch error
		}
    	findedNode := node.Search(Id)
        c.Code(200).Body(findedNode.GetData())
    })

    // GET all nodes from some node
    r.GET("/nodes/:id", func(c *router.Control) {
		c.UseTimer()
    	Id, err := strconv.Atoi(c.Get(":id"))
    	if err != nil {
			// TODO catch error
		}

		referrerNode := node.Search(Id)

    	referrerNode.Print()
        c.Code(200).Body(referrerNode.GetData())
    })

    // GET balance for current node
    r.GET("/balance/:id", func(c *router.Control) {
		c.UseTimer()
    	Id, err := strconv.Atoi(c.Get(":id"))
    	if err != nil {
			// TODO catch error
		}

    	findedNode := node.Search(Id)

		counter := 0

		ch := make(chan int)
		go func() {
			findedNode.Count(ch)
			close(ch)
		}()

		for value := range ch {
			counter += value
		}

        c.Code(200).Body(counter)
    })

    r.Listen(":3333")
}