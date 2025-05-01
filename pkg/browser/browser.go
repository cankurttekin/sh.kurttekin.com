package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenURL opens the specified URL in the system's default browser
func OpenURL(url string) error {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	
	return cmd.Start()
} 