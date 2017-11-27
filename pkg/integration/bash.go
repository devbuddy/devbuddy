package integration

var bash_source = `
# Mask the command dad with this shell function
# This let us mutate the current shell
dad() {
    [ -n "${DAD_DEBUG}" ] && echo "DAD_DEBUG: called with integration enabled"

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
        [ -n "${DAD_DEBUG}" ] && echo "DAD_DEBUG: finalizer: ${fin}"

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

complete -C 'dad --completion' dad

[ -n "${DAD_DEBUG}" ] && echo "DAD_DEBUG: Dad is now enabled..."
`
