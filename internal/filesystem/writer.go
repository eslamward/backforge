package filesystem

import "os"

func writeFile(path, contet string) error {

	return os.WriteFile(path, []byte(contet), 0644)
}
