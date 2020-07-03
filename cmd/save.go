/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save message to Drafts folder",
	Run: func(cmd *cobra.Command, args []string) {
		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...)
		defer cancel()

		ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
		defer cancel()

		subject, err := cmd.Flags().GetString("subject")
		if err != nil {
			log.Fatal(err)
		}

		attachFile, err := cmd.Flags().GetString("attach")
		if err != nil {
			log.Fatal(err)
		}

		if err := chromedp.Run(ctx, vmsSaveDraft(subject, attachFile)); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	draftCmd.AddCommand(saveCmd)

	saveCmd.Flags().StringP("subject", "s", "", "Subject")
	saveCmd.Flags().StringP("attach", "a", "", "Attach file")

	if err := saveCmd.MarkFlagRequired("attach"); err != nil {
		panic(err)
	}
}

func vmsSaveDraft(subject, attachFile string) chromedp.Tasks {
	selName := `//input[@id="username"]`
	selPass := `//input[@id="password"]`

	return chromedp.Tasks{
		chromedp.Navigate(viper.GetString("url")),
		chromedp.WaitVisible(selPass),
		chromedp.SendKeys(selName, viper.GetString("username")),
		chromedp.SendKeys(selPass, viper.GetString("password")),
		chromedp.Submit(selPass),
		chromedp.WaitVisible(`//div[@id="z_userName"]`),
		chromedp.Click(`//div[@title="Compose"]`),
		chromedp.WaitVisible(`//div[@id="zb__App__tab_COMPOSE-1"]`),
		chromedp.SendKeys(`//input[@id="zv__COMPOSE-1_subject_control"]`, subject),
		chromedp.SendKeys(`//input[@type="file"]`, attachFile, chromedp.NodeVisible),
		chromedp.WaitVisible(`//a[@class="AttLink"]`),
		chromedp.Click(`//div[@id="zb__COMPOSE-1__SAVE_DRAFT"]`),
	}
}
