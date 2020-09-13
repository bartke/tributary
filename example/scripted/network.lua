local tb = require('tributary')

tb.create_tester("tick_printer", ".")
tb.create_ticker("ticker_500ms", "500ms")
tb.create_ratelimit("filter_2s", "2s")
tb.create_forwarder("forwarder1")
tb.create_forwarder("forwarder2")

tb.link("ticker_500ms", "forwarder1")
tb.fanout("forwarder1", "filter_2s", "forwarder2")
tb.fanin("tick_printer", "filter_2s", "forwarder2")
