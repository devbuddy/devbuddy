package integration

var zshSource = `
promptcmd() { $("__bud_prompt_command") }
if [[ ! "${precmd_functions[@]}" == *__bud_prompt_command* ]]; then
  precmd_functions+=(promptcmd)
fi
`
