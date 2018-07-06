DEPS= \
	src/game/game.go \
	src/logger/logger.go \
	src/character/character.go \
	src/spells/spells.go \
	src/spells/spell_provider.go \
	src/utils/utils.go

wizard: src/wizard.go $(DEPS)
	GOPATH=$$(PWD):$$GOPATH go build src/wizard.go
