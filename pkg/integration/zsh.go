package integration

var zshSource = `
promptcmd() { $("__dad_prompt_command") }
if [[ ! "${precmd_functions[@]}" == *__dad_prompt_command* ]]; then
  precmd_functions+=(promptcmd)
fi
`
