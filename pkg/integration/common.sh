# Be careful! This runs in the user shell.

# Mask the command bud with this shell function
# This let us mutate the current shell
bud() {
    # Prepare a file to pass the finalize actions
    local finalizer_file
    finalizer_file="$(mktemp /tmp/bud-finalize-XXXXXX)"

    # Run the actual command
    env BUD_FINALIZER_FILE=$finalizer_file bud $@
    return_code=$?

    # Perform finalizers
    local fin
    while read -r fin; do
        __bud_log_debug "finalizer: ${fin}"

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
    rm -f "${finalizer_file}"

    return ${return_code}
}

__bud_prompt_command() {
    # In shell hook mode, the command will use stderr to print in the console
    # and stdout to mutate the shell (like activating a Python virtualenv)

    # Fail fast if no bud executable is reachable
    which bud > /dev/null || return

    local hook_eval
    hook_eval="$(command bud --shell-hook)"
    __bud_log_debug "Hook eval:\n--------\n${hook_eval}\n--------"
    eval "${hook_eval}"
}

__bud_log_debug() {
    [ -n "${BUD_DEBUG:-}" ] || return
    echo -e "\033[33mBUD_HOOK_DEBUG\033[0m: $@"
}

bud-enable-debug() {
    export BUD_DEBUG=1
    __bud_log_debug "debug log enabled"
}

bud-disable-debug() {
    __bud_log_debug "debug log disable"
    unset BUD_DEBUG
}

if [[ -n "${BUD_DEBUG:-}" ]]; then
    echo "BUD_DEBUG: DevBuddy is now enabled..."
fi
