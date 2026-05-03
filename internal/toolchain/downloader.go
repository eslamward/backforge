package toolchain

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadGo(url, outpath string) error {
	fmt.Println("Downloading go tool chain....")

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}
