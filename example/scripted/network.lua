local tb = require("tributary")

tb.link("ticker_1s", "filter_even")
tb.link("filter_even", "printer")
