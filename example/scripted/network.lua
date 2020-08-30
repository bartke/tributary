local tb = require("tributary")

tb.connect("ticker_1s", "filter_even")
tb.connect("filter_even", "printer")
