package helpers

import (
	"os"

	. "github.com/onsi/gomega"
)

func PushSpringMusicValidator() string {
	const path = "./spring-music-validator"
	originalWorkingDirectory, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())
	defer func() {
		err := os.Chdir(originalWorkingDirectory)
		Expect(err).NotTo(HaveOccurred())
	}()

	err = os.Chdir(path)
	Expect(err).NotTo(HaveOccurred())

	name := RandomName("spring-music-validator-%s")
	CF("push", name, "--no-start", "--no-route")
	return name
}
