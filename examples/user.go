package examples

//go:generate goor -type=User -getter -setter
type User struct {
	id   int    `goor:"constructor:-"`
	name string `goor:"getter:-"`
	age  int    `goor:"constructor:-;setter:-"`
}
