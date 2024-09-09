package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/harshalranjhani/genie/internal/middleware"
	"github.com/harshalranjhani/genie/internal/structs"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify [email]",
	Short: "Verify your support status and get access to extra features.",
	Long:  `If you have donated to the project, you can verify your email to get access to extra features. This command will send a verification email to the provided email address.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]

		status, err := middleware.LoadStatus()
		if err != nil {
			fmt.Println(color.RedString("Error loading status:"), err)
			return
		}

		if status != nil && status.Email == email {
			fmt.Println(color.GreenString("User is already verified."))
			return
		} else if status != nil && status.Email != email {
			fmt.Println(color.YellowString("A different user is already verified. Removing existing status..."))
			// remove the status file
			statusFile, err := middleware.GetStatusFilePath()
			if err != nil {
				fmt.Println(color.RedString("Error getting status file path:"), err)
				return
			}
			err = os.Remove(statusFile)
			if err != nil {
				fmt.Println(color.RedString("Error removing status file:"), err)
				return
			}
		}

		fmt.Println(color.CyanString("Sending verification email... Please do not close this tab until the process is complete."))
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Prefix = color.CyanString("Sending: ")
		s.Start()

		err = middleware.SendVerificationEmail(email)
		s.Stop()
		if err != nil {
			fmt.Println(color.RedString("Error sending verification email:"), err)
			return
		}

		fmt.Println(color.GreenString("Verification email sent. Please check your inbox."))
		fmt.Println(color.CyanString("Waiting for verification..."))

		s.Prefix = color.CyanString("Verifying: ")
		s.Start()
		token, err := middleware.WaitForVerification(email)
		s.Stop()
		if err != nil {
			fmt.Println(color.RedString("Error during verification:"), err)
			return
		}

		expiry := time.Now().Add(30 * 24 * time.Hour).Unix() // Calculate expiry timestamp (30 days)
		status = &structs.UserStatus{Email: email, Token: token, Expiry: expiry}
		err = middleware.SaveStatus(status)
		if err != nil {
			fmt.Println(color.RedString("Error saving status:"), err)
			return
		}

		fmt.Println(color.GreenString("Email verified successfully. You now have access to extra features!"))
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
