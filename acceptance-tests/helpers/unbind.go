package helpers

func Unbind(app, instance string) {
	CF("unbind-service", app, instance)
}
