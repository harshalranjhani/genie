package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

func startPayment() error {
	paymentURL := "https://rzp.io/rzp/JUFHvoxw"

	var err error
	switch os := runtime.GOOS; os {
	case "linux":
		err = exec.Command("xdg-open", paymentURL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", paymentURL).Start()
	case "darwin":
		err = exec.Command("open", paymentURL).Start()
	default:
		err = fmt.Errorf("Unable to open payment page on your platform. Please open the payment page manually: %s", paymentURL)
	}

	return err
}

var supportCmd = &cobra.Command{
	Use:   "support",
	Short: "Support the tool by donating to the project.",
	Long:  `Support genie and its development by donating to the project. This would allow you lifetime access to the tool and its updates. You can donate any amount you like.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Redirecting to the payment page...")
		err := startPayment()
		if err != nil {
			fmt.Println("Error starting payment:", err)
			return
		}

		fmt.Println("Please complete the payment in your browser.")
	},
}

func init() {
	rootCmd.AddCommand(supportCmd)
}
