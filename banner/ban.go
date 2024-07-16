package banner

import (
	"fmt"
	"os"
)

func ShowBanner() {
	ban, readErr := os.ReadFile("banner.txt")
	if readErr != nil {
		fmt.Fprintf(os.Stderr, "gohttp: Cannot read banner.txt: %v\n", readErr)
		return
	}
	fmt.Println(string(ban) + "\n")
}
