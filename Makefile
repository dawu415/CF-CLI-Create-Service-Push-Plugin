build:
	go build .

cf:
	cf uninstall-plugin Create-Service-Push || true
	cf install-plugin CF-CLI-Create-Service-Push-Plugin

it: build cf
