package helpers

func Bind(app, instance string) string {
	bindingName := RandomName("binding-%s")
	CF("bind-service", app, instance, "--binding-name", bindingName)
	return bindingName
}
