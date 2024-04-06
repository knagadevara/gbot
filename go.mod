module github.com/knagadevara/gbot

go 1.22.1

require golang.org/x/crypto v0.22.0

require (
	gopkg.in/yaml.v3 v3.0.1
)

require golang.org/x/sys v0.19.0 // indirect

replace github.com/knagadevara/gbot/utl => ./utl

replace github.com/knagadevara/gbot/commands => ./commands
