build:
	go build .

cf:
	cf uninstall-plugin CreateServicePush || true
	cf install-plugin CF-CLI-Create-Service-Push-Plugin

it: build cf
