ZMQ
It is also possible to subscribe to a bitcoin node and be notified about new transactions and new blocks via the node's ZMQ interface.

First, create a ZMQ instance:

zmq := bitcoin.NewZMQ("localhost", 28332)
Then create a buffered or unbuffered channel of strings and a goroutine to consume the channel:

	ch := make(chan string)

	go func() {
		for c := range ch {
			log.Println(c)
		}
	}()
Finally, subscribe to "hashblock" or "hashtx" topics passing in your channel:

	err := zmq.Subscribe("hashblock", ch)
	if err != nil {
		log.Fatalln(err)
	}