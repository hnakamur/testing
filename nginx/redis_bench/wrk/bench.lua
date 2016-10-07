

local requests = {
	"http://127.0.0.1/get?key=a&value=hello",
	"http://127.0.0.1/get?key=b&value=hola",
	"http://127.0.0.1/get?key=c&value=salut",
	"http://127.0.0.1/get?key=a",
	"http://127.0.0.1/get?key=b",
	"http://127.0.0.1/get?key=c",
}

request = function()
	local n = math.random(#requests)
	return wrk.format(nil, requests[n])
end
