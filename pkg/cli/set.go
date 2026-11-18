package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/coordination"
	yc "github.com/ydb-platform/ydb-go-yc"
)

var setCmd = &cobra.Command{
	Use:   "set [flags]",
	Short: "Set the value from file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv("IAM_TOKEN")
		opts := []ydb.Option{
			yc.WithInternalCA(),
		}
		if token != "" {
			opts = append(opts, ydb.WithAccessTokenCredentials(token))
		} else if saFile != "" {
			opts = append(opts, yc.WithServiceAccountKeyFileCredentials(saFile))
		} else {
			panic("either IAM_TOKEN or sa-file is required")
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		db, err := ydb.Open(ctx, dsn,
			opts...,
		)
		if err != nil {
			panic(err)
		}
		defer func() { _ = db.Close(ctx) }()

		scheme, nodeCfg, err := db.Coordination().DescribeNode(ctx, "/traefik")
		if err != nil {
			fmt.Printf("failed to describe node: %v", err)
			//return
		}
		fmt.Printf("Node scheme: %s\n", scheme)
		fmt.Printf("Node config: %v\n", nodeCfg)

		if scheme == nil {
			err := db.Coordination().CreateNode(ctx, "/traefik", coordination.NodeConfig{
				Path:                     "/traefik",
				SelfCheckPeriodMillis:    1000,
				SessionGracePeriodMillis: 1000,
				ReadConsistencyMode:      0,
				AttachConsistencyMode:    0,
				RatelimiterCountersMode:  0,
			})
			if err != nil {
				fmt.Printf("failed to create node: %v", err)
				return
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
