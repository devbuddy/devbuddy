package integration

var zshSource = `
prmptcmd() { $("__dad_prompt_command") }
precmd_functions=(prmptcmd)
`
