package integration

var bashSource = `
if [[ ! "${PROMPT_COMMAND:-}" == *__dad_prompt_command* ]]; then
  PROMPT_COMMAND="__dad_prompt_command; ${PROMPT_COMMAND:-}"
fi
`
