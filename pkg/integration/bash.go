package integration

var bashSource = `
# Be careful! This runs in the user shell.

# Mask the command dad with this shell function
# This let us mutate the current shell
dad() {
    # Prepare a file to pass the finalize actions
    local finalizer_file
    finalizer_file="$(mktemp /tmp/dad-finalize-XXXXXX)"
    trap "rm -f '${finalizer_file}'" EXIT

    # Run the actual command
    env DAD_FINALIZER_FILE=$finalizer_file dad $@
    return_code=$?

    # Perform finalizers
    local fin
    while read -r fin; do
        [ -n "${DAD_DEBUG:-}" ] && echo "DAD_DEBUG: finalizer: ${fin}"

        case "${fin}" in
            cd:*)
                cd "${fin//cd:/}"
                ;;
            setenv:*)
                export "${fin//setenv:/}"
                ;;
            *)
                ;;
        esac
    done < "${finalizer_file}"

    return ${return_code}
}

__dad_prompt_command() {
    # In shell hook mode, the command will use stderr to print in the console
    # and stdout to mutate the shell (like activating a Python virtualenv)

    # Fail fast if no dad executable is reachable
    which dad > /dev/null || return

    local hook_eval
    hook_eval="$(dad --shell-hook)"
    [ -n "${DAD_DEBUG:-}" ] && echo -e "DAD_DEBUG: Hook eval:\n${hook_eval}\n---"
    eval "${hook_eval}"
}

if [[ ! "${PROMPT_COMMAND:-}" == *__dad_prompt_command* ]]; then
  PROMPT_COMMAND="__dad_prompt_command; ${PROMPT_COMMAND:-}"
fi

dad-enable-debug() {
    export DAD_DEBUG=1
    echo "DAD_DEBUG: enabled"
}

dad-disable-debug() {
    unset DAD_DEBUG
    echo "DAD_DEBUG: disable"
}

if [[ -n "${DAD_DEBUG:-}" ]]; then
    echo "DAD_DEBUG: Dad is now enabled..."
fi
`
