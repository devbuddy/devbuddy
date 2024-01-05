if [[ ! "${PROMPT_COMMAND:-}" == *__bud_prompt_command* ]]; then
  PROMPT_COMMAND="__bud_prompt_command; ${PROMPT_COMMAND:-}"
fi
