snake: 
	go build -ldflags "-s -w" -o $@
clean:
	rm snake 2>/dev/null &>&2
