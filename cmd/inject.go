/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	protoFiles "github.com/stevezaluk/protoc-go-inject-tag/file"
	"github.com/stevezaluk/protoc-go-inject-tag/inject"
	"log/slog"
	"os"
)

var injectCmd = &cobra.Command{
	Use:   "inject",
	Short: "Inject tags into your protobuf models",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("verbose") {
			slog.SetLogLoggerLevel(slog.LevelDebug)
		}

		if viper.GetString("input") == "" {
			slog.Error("input path is empty. A valid UNIX path must be passed")
			os.Exit(1)
		}

		inject.InitRegex()
	},
	Run: func(cmd *cobra.Command, args []string) {
		protoFiles.WalkDir(viper.GetString("input"))
	},
}

func init() {
	rootCmd.AddCommand(injectCmd)

	injectCmd.Flags().StringP("input", "i", "", "The input path you want to search for protobufs with")
	viper.BindPFlag("input", injectCmd.Flags().Lookup("input"))

	injectCmd.Flags().BoolP("verbose", "v", false, "Enable extended verbosity in logging")
	viper.BindPFlag("verbose", injectCmd.Flags().Lookup("verbose"))

	injectCmd.Flags().StringP("file-ext", "f", ".pb.go", "The file extensions that should be considered for injection")
	viper.BindPFlag("tag.file-ext", injectCmd.Flags().Lookup("file-ext"))

	injectCmd.Flags().Bool("remove-comments", false, "Remove comments from generated protobufs")
	viper.BindPFlag("tag.remove-comments", injectCmd.Flags().Lookup("remove-comments"))

	injectCmd.Flags().String("comment-prefix", "gotags", "The prefix of the comment that protoc-go-inject-tag should search for when looking for tags to inject. A @ will be prefixed to this value")
	viper.BindPFlag("tag.comment-prefix", injectCmd.Flags().Lookup("comment-prefix"))

	injectCmd.Flags().StringSlice("tags", nil, "Additional tags that should be applied to all fields")
	viper.BindPFlag("tag.additional-tags", injectCmd.Flags().Lookup("tags"))
}
