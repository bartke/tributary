local tb = require("tributary")

tb.create_forwarder("fwd1")
tb.create_forwarder("fwd2")

tb.link("ticker_1s", "fwd1")
tb.fanout("fwd1", "filter_even", "fwd2")
tb.fanin("printer", "filter_even", "fwd2")

print("fwd1 exists", tb.node_exists("fwd1"))
print("fwd2 exists", tb.node_exists("fwd2"))
print("fwd3 exists", tb.node_exists("fwd3"))

tb.run("ticker_1s")
tb.run("filter_even")
tb.run("fwd1")
tb.run("fwd2")
tb.run("printer")
