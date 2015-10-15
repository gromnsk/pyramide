// To Do
// -------
// Make delete a little cleaner
// Add traversal

package main

import (
	"github.com/gromnsk/tree"
	"github.com/takama/router"
	"strconv"
	"encoding/json"
	"io/ioutil"
	// "net/http"
	"fmt"
	// "net/url"
	// "encoding/json"
)

var path = "/var/lib/pyramide/nodes.json"

type Result struct {
	Success bool

}

func main() {
	root := new(tree.Tree)
	data := tree.Data{Id: 1}
	node := root.Insert(data)

	// for index := 2; index < 50000; index++ {
	// 	data := tree.Data{Id: index}
	// 	node.Insert(data)
	// }
	
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

    // GET all nodes from some node with level limit
    r.GET("/nodes/:id/:level", func(c *router.Control) {
		c.UseTimer()
    	Id, err := strconv.Atoi(c.Get(":id"))
    	if err != nil {
			// TODO catch error
		}

    	Level, err := strconv.Atoi(c.Get(":level"))
    	if err != nil {
			// TODO catch error
		}

		referrerNode := node.Search(Id)

        c.Code(200).Body(referrerNode.GetNodes(Level))
    })

    // GET balance for current node
    r.GET("/balance/:id", func(c *router.Control) {
		c.UseTimer()
    	Id, err := strconv.Atoi(c.Get(":id"))
    	if err != nil {
			// TODO catch error
		}
		
		referrerNode := node.Search(Id)
		data := referrerNode.GetNodes(8)

        c.Code(200).Body(len(data))
    })

    // GET balance for current node
    r.GET("/dump", func(c *router.Control) {
    	c.UseTimer()
		data := node.GetNodes(10000)
		j, jerr := json.MarshalIndent(data, "", "  ")
		if jerr != nil {
			fmt.Println("jerr:", jerr.Error())
		}

		ioutil.WriteFile(path, j, 0644)

        c.Code(200).Body(len(data))
    })

    // GET balance for current node
    r.GET("/restore", func(c *router.Control) {
    	c.UseTimer()
		data, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("error reading json: ", err.Error())
		}
		var restored []*tree.Result
		jerr := json.Unmarshal(data, &restored)
		if jerr != nil {
			fmt.Println("jerr:", jerr.Error())
		}

		node.SetAllNodes(restored)

        c.Code(200).Body(node)
    })

    r.Listen(":3333")
}