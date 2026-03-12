package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	"github.com/go-jira/jira/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type SprintAddOptions struct {
	jiradata.EpicIssues `yaml:",inline" json:",inline" figtree:",inline"`
	Project             string `yaml:"project,omitempty" json:"project,omitempty"`
}

func CmdSprintAddRegistry() *jiracli.CommandRegistryEntry {
	opts := SprintAddOptions{}

	return &jiracli.CommandRegistryEntry{
		"Add issues to the current active sprint",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdSprintAddUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			for i := range opts.Issues {
				opts.Issues[i] = jiracli.FormatIssue(opts.Issues[i], opts.Project)
			}
			return CmdSprintAdd(o, globals, &opts)
		},
	}
}

func CmdSprintAddUsage(cmd *kingpin.CmdClause, opts *SprintAddOptions) error {
	cmd.Arg("ISSUE", "Issue(s) to add to the current sprint").Required().StringsVar(&opts.Issues)
	return nil
}

func CmdSprintAdd(o *oreo.Client, globals *jiracli.GlobalOptions, opts *SprintAddOptions) error {
	sprint, err := jira.GetActiveSprint(o, globals.Endpoint.Value, opts.Project)
	if err != nil {
		return jiracli.CliError(err)
	}

	if err := jira.SprintAddIssues(o, globals.Endpoint.Value, sprint.ID, &opts.EpicIssues); err != nil {
		return jiracli.CliError(err)
	}

	if !globals.Quiet.Value {
		for _, issue := range opts.Issues {
			fmt.Printf("OK %s %s\n", issue, jira.URLJoin(globals.Endpoint.Value, "browse", issue))
		}
	}

	return nil
}
