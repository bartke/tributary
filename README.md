# Tributary

Simple **event stream processor** written in Go and Lua. Tributary allows to create networks and
define events that propagate through connected nodes. Network nodes process events and can engage
external resources, can decrease or multiply inputs from their in to their out ports. Tributary
networks are created and managed by a lua runtime that exposes networking primitives and allows
for custom node manipulators that can be registered through the `tributary` module.

- flow based, isolates concurrently running network nodes
- network nodes send events through Go channels through their in and out ports
- networks can be created dynamically with a lua based runtime

Possibility of simple **sliding windows** and event aggregate and join queries using SQL Databases
- includes examples for windowing and query pipelines using Gorm. Working with sqlite/mysql
  - implement sliding windows with time or limit based query conditions
- clear window after alert successfully sent from next network node
- filter outputs to avoid duplicates

Example **use cases**
- ingest customer actions from distributed systems that communicate event driven over a message bus
- pipe aggregate customer actions back onto a message bus based on query criteria
- send alerts to slack/telegram, http callbacks
- start workflows based on aggregate customer actions

## Examples

### Concepts

![network](./example/scripted/network.svg)

As per [example/scripted](example/scripted/network.lua), we can set up such network at runtime
with lua scripts, e.g. here, showing a a ticker source, a debug printer, and node connections via
direct **link**, via **fan-out** and **fan-in**.

```lua
local tb = require('tributary')

tb.create_tester("tick_printer", ".")
tb.create_ticker("ticker_500ms", "500ms")
tb.create_ratelimit("filter_2s", "2s")
tb.create_forwarder("forwarder1")
tb.create_forwarder("forwarder2")

tb.link("ticker_500ms", "forwarder1")
tb.fanout("forwarder1", "filter_2s", "forwarder2")
tb.fanin("tick_printer", "filter_2s", "forwarder2")
```

Create a network that module will operate on, load the module into the runtime, then execute the
above script. The script will link up and control network nodes. When executed without error, run
the network nodes.

```go
n := network.New() // ... add custom nodes
m := module.New(n) // ... add custom exports to be available in the runtime
r := runtime.New()
r.LoadModule(m.Loader) // register the tributary module with the runtime
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

### Custom Nodes and Module Exports

Custom network nodes can be added by either implementing the **Source**, **Pipeline** or **Sink**
interfaces, or by utilizing provided abstractions for **Injectors**, **Interceptors** and
**Handlers** function types. A handler example would be a simple print function, e.g. we can
route messages to a `printer` node that is added as follows

```go
out := func(e tributary.Event) {
	fmt.Println(string(e.Payload()))
}
n := network.New()
n.AddNode("printer", handler.New(out))
```

Equally, we can add additional exports on the tributary lua module:

```go
myCustomFn := func(l *lua.LState) int {
	firstStringArg := l.CheckString(1)
	secondIntegerArg := l.CheckInt(2)
	// ... do sth with string arg 1 and integer arg 2
	// add boolean return value
	l.Push(module.LuaConvertValue(l, true))
	return 1
}
m := module.New(n)
m.Export("my_custom_fn", myCustomFn)
```

In the networking script we can now call `my_custom_fn` on the required module:

```lua
local tb = require('tributary')
tb.my_custom_fn("arg1", 2)
-- ...
```

### Available module functions

The following lua script prints all available functions, exported from the `tributary` module.

```lua
local tb = require("tributary")

print("functions:")
for i,v in pairs(tb) do
    if type(v) == "function" then
        print(i)
    end
end
```

The output is

```
functions:
sliding_window_time
create_filter
node_exists
fanin
create_ratelimit
create_tester
create_window
query_window
link
fanout
create_ticker
create_forwarder
create_discarder
```

### Sliding Window CEP

Create a windower on top of a mysql database and Gorm. With mysql we can span sliding windows
with time ranges and limits if desired on an incrementally built dataset.

```go
db := mysql.Open("root:root@tcp(localhost:3306)/tb")
gormCfg := &gorm.Config{
	SkipDefaultTransaction: true,
	PrepareStmt:            true,
	NamingStrategy: schema.NamingStrategy{
		TablePrefix: "example_",
	},
}
window, err := gormwindow.New(db, gormCfg, standardevent.New,
	&event.Bet{},
	&event.Selection{})
if err != nil {
	log.Fatal(err)
}
```

Add lua module exports for the windower with a struct, annotated for gorm tables

```go
m := module.New(n)
m.AddWindowExports(window, &event.Bet{})
```

We can then create and query a window on the runtime and link it up. The following creates a
window `bets_window` from a streaming ingest of type `event.Bet`, fowards the data stream to
`window_query` and the output ultimately to a printer.

```lua
local tb = require("tributary")

tb.create_window("bets_window")
tb.link("streaming_ingest", "bets_window")

-- (mysql) select aggregate customer liability if > 130
local query = [[
select
	round(sum(stake*exchange_rate*(b.odds -1)), 2) as liability,
	customer_uuid,
	game_id
from
	bets b
join selections s on
	b.uuid = s.bet_uuid
where
	FROM_UNIXTIME(b.create_time / 1000000000) >= now() - interval 10 second
group by
	customer_uuid,
	game_id
having
	liability > 130
]]
tb.query_window("window_query", query)
tb.link("bets_window", "window_query")
tb.link("window_query", "printer")
```

This example can be found under [example/advanced](example/advanced).

### Todos & Notes

- error if _ node is added
- converge window create, query, filter cleanup, add window cleanup
- direct matcher, filter on attribute list, flat queries on one messages, `map[string]interface{}`. attribute ><>
- graceful network node shutdown
- network with NATS/Rabbitmq etc
- network stats?
- telegram/slack/callback sink
- pubsub/rabbitmq source
- client
	- manage runtime networks in db
	- build graph, multi instance, pick up each script once, select dependencies
