package helpers

func DeleteApp(name string) {
	CF("delete", "-f", name)
}
