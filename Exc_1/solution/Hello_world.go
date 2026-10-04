package main

import "fmt"

type Message struct {
	Greeting, Target string
}

func GreetUser(h Message) {
	fmt.Println(h.Greeting, h.Target)
}

func main() {
	greeting := Message{"Hello", "World"}
	GreetUser(greeting)

}
