package integration

var zshSource = `
promptcmd() { $("__dad_prompt_command") }
precmd_functions=(promptcmd)
`
