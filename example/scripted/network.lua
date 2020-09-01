local tb = require("tributary")

-- tb.link("ticker_1s", "filter_even")
-- tb.link("filter_even", "printer")

tb.create_forwarder("fwd1")
tb.fanout("ticker_1s", "filter_even", "fwd1")
tb.fanin("printer", "filter_even", "fwd1")

print(tb.node_exists("fwd1"))
print(tb.node_exists("fwd2"))

tb.run("ticker_1s")
tb.run("filter_even")
tb.run("fwd1")
tb.run("printer")
