package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync/atomic"
	"time"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type SnowflakeGenerator struct {
	epoch      int64
	machineID  int64
	sequence   uint64
	lastTime   int64
}

func getSnowflakeGenerator() *SnowflakeGenerator {
	rand.Seed(time.Now().UnixNano())
	machineID := rand.Int63n(1024)

	return &SnowflakeGenerator{
		epoch:     1420070400000,
		machineID: machineID,
	}
}

func (g *SnowflakeGenerator) getUniqueId() int64 {
	now := time.Now().UnixNano()
	seq := atomic.AddUint64(&g.sequence, 1)
	seq &= 0xFFF

	return ((now - g.epoch) << 22) |
		   (g.machineID << 12) |
		   int64(seq)
}

func main() {

    n := maelstrom.NewNode()
    gen := getSnowflakeGenerator()
	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "generate_ok"
		body["id"] = gen.getUniqueId()
 		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
