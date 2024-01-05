if [[ ! "${precmd_functions[@]}" == *__bud_prompt_command* ]]; then
  precmd_functions+=(__bud_prompt_command)
fi
