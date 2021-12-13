coverage:
	go test -v ./pkg/... -coverprofile=c.out && go tool cover -html=c.out && rm c.out