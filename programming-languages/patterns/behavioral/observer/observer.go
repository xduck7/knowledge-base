package observer

import "fmt"

// Observer

type Subscriber interface {
	Notify(videoTitle string)
	ID() string
}

type User struct {
	name string
}

func (u *User) Notify(videoTitle string) {
	fmt.Printf("[%s] получил уведомление: новое видео \"%s\"\n", u.name, videoTitle)
}

func (u *User) ID() string { return u.name }

// Subject (канал)

type Channel struct {
	name        string
	subscribers map[string]Subscriber
}

func NewChannel(name string) *Channel {
	return &Channel{
		name:        name,
		subscribers: make(map[string]Subscriber),
	}
}

func (c *Channel) Subscribe(s Subscriber) {
	c.subscribers[s.ID()] = s
}

func (c *Channel) Unsubscribe(s Subscriber) {
	delete(c.subscribers, s.ID())
}

func (c *Channel) Publish(videoTitle string) {
	fmt.Printf("Канал %q опубликовал видео: %s\n", c.name, videoTitle)

	for _, sub := range c.subscribers {
		sub.Notify(videoTitle)
	}
}

func Example() {
	channel := NewChannel("Golang Channel")

	alex := &User{name: "Alex"}
	ira := &User{name: "Ira"}

	channel.Subscribe(alex)
	channel.Subscribe(ira)

	channel.Publish("Паттерн Observer в Go") // уведомятся оба

	channel.Unsubscribe(ira)

	channel.Publish("Паттерн Strategy в Go") // уведомится только Alex
}
