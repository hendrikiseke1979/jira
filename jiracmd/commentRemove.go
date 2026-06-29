package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	jira "github.com/hendrikiseke1979/jira"
	"github.com/hendrikiseke1979/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type CommentRemoveOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Issue                 string `yaml:"issue,omitempty" json:"issue,omitempty"`
	CommentID             string `yaml:"comment-id,omitempty" json:"comment-id,omitempty"`
}

func CmdCommentRemoveRegistry() *jiracli.CommandRegistryEntry {
	opts := CommentRemoveOptions{}

	return &jiracli.CommandRegistryEntry{
		"Delete comment from issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdCommentRemoveUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdCommentRemove(o, globals, &opts)
		},
	}
}

func CmdCommentRemoveUsage(cmd *kingpin.CmdClause, opts *CommentRemoveOptions) error {
	cmd.Arg("ISSUE", "Issue to delete comment from").Required().StringVar(&opts.Issue)
	cmd.Arg("COMMENT-ID", "Comment id to delete").Required().StringVar(&opts.CommentID)
	return nil
}

func CmdCommentRemove(o *oreo.Client, globals *jiracli.GlobalOptions, opts *CommentRemoveOptions) error {
	if err := jira.IssueRemoveComment(o, globals.Endpoint.Value, opts.Issue, opts.CommentID); err != nil {
		return err
	}

	if !globals.Quiet.Value {
		fmt.Printf("OK Deleted comment %s from %s\n", opts.CommentID, opts.Issue)
	}
	return nil
}
