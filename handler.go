package main

import "net/http"

func IndexHandler(c *Context) {
	c.String(http.StatusOK, "Hello, World!")
}

func HelloHandler(c *Context) {
	name := c.Param("name")
	c.String(http.StatusOK, "Hello, %s!", name)
}
