package main

import (
    "encoding/json"
    "log"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)


func main() {
    messages := []int{} 
    n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

        messages = append(messages, int(body["message"].(float64)))

        return n.Reply(msg, map[string]any{
		    "type": "broadcast_ok",
	    })
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

        return n.Reply(msg, map[string]any{
            "messages": messages,
		    "type": "read_ok",
	    })
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

        return n.Reply(msg, map[string]any{
		    "type": "topology_ok",
	    })
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
