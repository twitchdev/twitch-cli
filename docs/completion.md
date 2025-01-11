# Completion

Generate autocompletions for your desired shell.

## Arguments

| Argument   | Description                                       |
| ---------- | ------------------------------------------------- |
| bash       | Generate the autocompletion script for bash       |
| fish       | Generate the autocompletion script for fish       |
| powershell | Generate the autocompletion script for powershell |
| zsh        | Generate the autocompletion script for zsh        |

## Usage

### Bash

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

```bash
source <(twitch completion bash)
```

To load completions for every new session, execute once:

#### Linux:

```bash
twitch completion bash > /etc/bash_completion.d/twitch
```

#### macOS:

```bash
twitch completion bash > $(brew --prefix)/etc/bash_completion.d/twitch
```

You will need to start a new shell for this setup to take effect.

### Zsh

If shell completion is not already enabled in your environment you will need
to enable it. You can execute the following once:

```zsh
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

To load completions in your current shell session:

```zsh
source <(twitch completion zsh)
```

To load completions for every new session, execute once:

#### Linux:

```zsh
twitch completion zsh > "${fpath[1]}/_twitch"
```

#### macOS:

```zsh
twitch completion zsh > $(brew --prefix)/share/zsh/site-functions/_twitch
```

You will need to start a new shell for this setup to take effect.

### Fish

To load completions in your current shell session:

```fish
twitch completion fish | source
```

To load completions for every new session, execute once:

```fish
twitch completion fish > ~/.config/fish/completions/twitch.fish
```

You will need to start a new shell for this setup to take effect.

### PowerShell

To load completions in your current shell session:

```powershell
twitch completion powershell | Out-String | Invoke-Expression
```

To load completions for every new session, add the output of the above command
to your powershell profile.
