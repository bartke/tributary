## Tributary

Simple Go event stream processor. Tributary allows to create networks and define events that
propagate through connected nodes. Network nodes process events and can engage external
resources, can decrease or multiply inputs from their in to their out ports. Tributary networks
are created and managed by a lua based runtime, exposing networking primitives and possible
custom node manipulators though the `tributary` lua module.

- flow based, isolates concurrently running network nodes
- network nodes send events through Go channels through their in and out ports
- networks can be created dynamically with a lua based runtime

Simple Sliding Windows using SQL
- includes examples for windowing and query pipelines using Gorm. Working with sqlite/mysql
  - implement sliding windows with time or limit based query conditions
- clear window after alert successfully sent from next network node
- filter outputs to avoid duplicates

Example use cases
- pipe aggregate customer actions back onto a message bus based on query criteria
- alert to slack or telegram

TODOs
- `[]byte` json/octet-stream payloads
- filter reported or clear, uniqueness
- converge window create, query, filter cleanup, add window cleanup
- direct matcher, filter on attribute list, flat queries on one messages, `map[string]interface{}`. attribute ><>
- arbitrary event types?
- network could be created with NATS/Rabbitmq etc
- network stats
- telegram/slack/callback sink
- pubsub/rabbitmq source

Client Distributary:
- manage lua scripts in db
- build graph, multi instance, pick up each script once, select dependencies on graph

### Example

![network](./example/scripted/network.svg)

As per [example/scripted](example/scripted/network.lua), we can set up such network at runtime
with lua scripts, e.g. here

```lua
local tb = require('tributary')

tb.create_tester("debug_print", ".")
tb.create_ticker("ticker_500ms", "500ms")
tb.create_ratelimit("filter_2s", "2s")
tb.create_forwarder("forwarder1")
tb.create_forwarder("forwarder2")
tb.link("ticker_500ms", "forwarder1")
tb.fanout("forwarder1", "filter_2s", "forwarder2")
tb.fanin("debug_print", "filter_2s", "forwarder2")
```

Create a network that module will operate on, load the module into the runtime, then execute the
above script. The script will link up and control network nodes. When executed without error, run
the network nodes.

```go
	// create network and register nodes
	n := network.New()
	// ... add nodes for sources

	// create the tributary module that operates on the network and register the tributary module
	// exports with the runtim
	m := module.New(n)
	// ... add exports to be available in the runtime

	r := runtime.New()
	r.LoadModule(m.Loader)
	// Run will preload the tributary module and execute a script on the VM. We can close it
	// after we called Run() to stop the execution.
	if err := r.Run("./network.lua"); err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// after nodes are linked up and the script is loaded without errors, run the network
	n.Run()
```

We can print the network to a Graphviz output shown above with

```go
fmt.Println(tributary.Graphviz(n)
```
