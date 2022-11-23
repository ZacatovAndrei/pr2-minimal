build: clean
	-cd producer && go build 
	-cd aggregator && go build
	-cd consumer && go build 
clean:
	-cd producer && go clean 
	-cd aggregator && go clean
	-cd consumer && go clean 
