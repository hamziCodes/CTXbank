package main

import (
	"fmt"
	"os"
	"strings"
)

// ctxCommands is the canonical command list shared by all completion scripts.
var ctxCommands = []string{
	"init", "status", "ui", "dashboard", "list",
	"pause", "checkpoint", "resume", "audit", "ingest",
	"lint-memory", "lint", "prompt-sync", "sync-prompt",
	"serve", "doctor", "completion", "version", "help",
}

var bashCompletion = `# ctx bash completion — install: ctx completion bash > /etc/bash_completion.d/ctx
_ctx_completions() {
    local cur prev cmds
    cmds="` + strings.Join(ctxCommands, " ") + `"
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    case "$prev" in
        ctx) COMPREPLY=($(compgen -W "$cmds" -- "$cur")); return 0 ;;
        pause) COMPREPLY=($(compgen -W "--note --focus --help" -- "$cur")); return 0 ;;
        status|list) COMPREPLY=($(compgen -W "--json --help" -- "$cur")); return 0 ;;
        ui) COMPREPLY=($(compgen -W "--port --no-open --help" -- "$cur")); return 0 ;;
        audit) COMPREPLY=($(compgen -W "--apply --help" -- "$cur")); return 0 ;;
        ingest) COMPREPLY=($(compgen -W "--yes --help" -- "$cur")); return 0 ;;
        lint-memory|lint) COMPREPLY=($(compgen -W "--fix --help" -- "$cur")); return 0 ;;
        prompt-sync|sync-prompt) COMPREPLY=($(compgen -W "--verify --status --help" -- "$cur")); return 0 ;;
        serve) COMPREPLY=($(compgen -W "--mcp --help" -- "$cur")); return 0 ;;
        completion) COMPREPLY=($(compgen -W "bash zsh fish powershell" -- "$cur")); return 0 ;;
    esac
    COMPREPLY=($(compgen -W "--help" -- "$cur"))
}
complete -F _ctx_completions ctx
`

var zshCompletion = `# ctx zsh completion — install: ctx completion zsh > ~/.zsh/completions/_ctx
#compdef ctx
_ctx() {
    local -a cmds
    cmds=(` + zshCmdList() + `)
    _arguments -C \
        '1:command:->cmd' \
        '*:: :->args'
    case $state in
        cmd) _describe 'ctx commands' cmds ;;
        args)
            case $words[2] in
                pause) _arguments '--note[one-line note]:note:' '--focus[focus label]:focus:' ;;
                status|list) _arguments '--json[JSON output]' ;;
                ui) _arguments '--port[port]:port:' '--no-open[do not open browser]' ;;
                audit) _arguments '--apply[apply to memory bank]' ;;
                ingest) _arguments '--yes[skip confirmation]' ;;
                lint-memory|lint) _arguments '--fix[auto-prune]' ;;
                prompt-sync|sync-prompt) _arguments '--verify[verify]' '--status[ledger status]' ;;
                serve) _arguments '--mcp[start MCP server]' ;;
                completion) _arguments '1:shell:(bash zsh fish powershell)' ;;
            esac ;;
    esac
}
_ctx
`

func zshCmdList() string {
	quoted := make([]string, len(ctxCommands))
	for i, c := range ctxCommands {
		quoted[i] = "'" + c + "'"
	}
	return strings.Join(quoted, " ")
}

var fishCompletion = `# ctx fish completion — install: ctx completion fish > ~/.config/fish/completions/ctx.fish
set -l cmds ` + strings.Join(ctxCommands, " ") + `
complete -c ctx -f -n '__fish_use_subcommand' -a "$cmds"
complete -c ctx -f -n '__fish_seen_subcommand_from pause' -l note -d 'One-line note'
complete -c ctx -f -n '__fish_seen_subcommand_from pause' -l focus -d 'Focus label'
complete -c ctx -f -n '__fish_seen_subcommand_from status; and __fish_seen_subcommand_from list' -l json -d 'JSON output'
complete -c ctx -f -n '__fish_seen_subcommand_from ui' -l port -d 'Port'
complete -c ctx -f -n '__fish_seen_subcommand_from ui' -l no-open -d 'Do not open browser'
complete -c ctx -f -n '__fish_seen_subcommand_from audit' -l apply -d 'Apply to memory bank'
complete -c ctx -f -n '__fish_seen_subcommand_from ingest' -l yes -d 'Skip confirmation'
complete -c ctx -f -n '__fish_seen_subcommand_from lint-memory' -l fix -d 'Auto-prune'
complete -c ctx -f -n '__fish_seen_subcommand_from prompt-sync' -l verify -d 'Verify'
complete -c ctx -f -n '__fish_seen_subcommand_from prompt-sync' -l status -d 'Ledger status'
complete -c ctx -f -n '__fish_seen_subcommand_from serve' -l mcp -d 'Start MCP server'
complete -c ctx -f -n '__fish_seen_subcommand_from completion' -a 'bash zsh fish powershell'
`

var powershellCompletion = `# ctx PowerShell completion — install: ctx completion powershell >> $PROFILE
Register-ArgumentCompleter -Native -CommandName ctx -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    $cmds = @(` + psCmdList() + `)
    $tokens = $commandAst.ToString() -split '\s+'
    if ($tokens.Count -le 2) {
        $cmds | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
        }
    } elseif ($tokens[1] -eq 'completion') {
        @('bash','zsh','fish','powershell') | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
        }
    }
}
`

func psCmdList() string {
	quoted := make([]string, len(ctxCommands))
	for i, c := range ctxCommands {
		quoted[i] = "'" + c + "'"
	}
	return strings.Join(quoted, ",")
}

func runCompletion(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: ctx completion <bash|zsh|fish|powershell>\n")
		os.Exit(1)
	}
	switch args[0] {
	case "bash":
		fmt.Print(bashCompletion)
	case "zsh":
		fmt.Print(zshCompletion)
	case "fish":
		fmt.Print(fishCompletion)
	case "powershell":
		fmt.Print(powershellCompletion)
	default:
		fmt.Fprintf(os.Stderr, "Unknown shell %q. Choose: bash, zsh, fish, powershell.\n", args[0])
		os.Exit(1)
	}
}
