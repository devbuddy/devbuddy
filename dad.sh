function _dad_help {
    echo "Usage: dad [command] ..."
    echo ""
    echo "Commands:"
    echo ""
    echo "  new NAME          Create a new project"
    echo ""
    echo "  cd NAME           Go to a project and activate environment"
    echo ""
    echo "  up                Prepare your development environment"
    echo ""
    echo "  test              Run the test suite (alias: t)"
    echo ""
    echo "  server            Run the server"
    echo ""
}

function _dad_up {
    local reset=$(tput sgr0);
    local white=$(tput setaf 15);
    local yellow=$(tput setaf 136);
    local cyan=$(tput setaf 37);

    if [[ -e 'Pipfile' ]]; then
        cmd="pipenv install --dev"
        echo "${yellow}★ ${cyan}Running ${white}${cmd}${reset}"
        $cmd
    fi

    if [[ -e 'requirements.txt' ]]; then
        cmd="pip install -r requirements.txt"
        echo "${yellow}★ ${cyan}Running ${white}${cmd}${reset}"
        $cmd
    fi
}

function _dad_test {
    py.test $@
}

function _dad_server {
    pserve development.ini $@
}

function _dad_new {
    local reset=$(tput sgr0);
    local white=$(tput setaf 15);
    local yellow=$(tput setaf 136);
    local cyan=$(tput setaf 37);

    local path="$HOME/src/github.com/pior/$1"
    if [[ ! -e "$path" ]]; then
        echo "${yellow}★ ${cyan}Creating project in ${white}${path}${reset}"
        mkdir $path
    fi
    _dad_cd $1
}

function _dad_cd_cleanup {
    local reset=$(tput sgr0);
    local yellow=$(tput setaf 136);
    local cyan=$(tput setaf 37);

    if [[ -n "${PIPENV_ACTIVE}" ]]; then
        echo "${yellow}★ ${cyan}De-activating PipEnv ${reset}${VIRTUAL_ENV}"
        unset PIPENV_ACTIVE
        deactivate
    fi

    if [[ -n "${VIRTUAL_ENV}" ]]; then
        echo "${yellow}★ ${cyan}De-activating venv ${reset}${VIRTUAL_ENV}"
        deactivate
    fi

}

function _dad_cd {
    local reset=$(tput sgr0);
    local white=$(tput setaf 15);
    local yellow=$(tput setaf 136);
    local cyan=$(tput setaf 37);

    if [[ -z "$1" ]]; then
        if [[ -z "$DAD_ACTIVE_PROJECT_PATH" ]]; then
            echo "Usage: dad cd [NAME]"
            return
        fi

        echo "${yellow}★ ${cyan}Jumping to root of active project${reset}"
        cd ${DAD_ACTIVE_PROJECT_PATH}
        return
    fi

    if [[ -e "$HOME/src/github.com/pior/$1" ]]; then
        _dad_cd_cleanup

        echo "${yellow}★ ${cyan}Jumping to ${white}pior/$1${reset}"
        local path=~/src/github.com/pior/$1
        cd ${path}
        DAD_ACTIVE_PROJECT_PATH=${path}

        if [[ -e 'Pipfile' ]]; then
            local venv=$(pipenv --venv)
            if [[ -n "${venv}" ]]; then
                echo "${yellow}★ ${cyan}Activating PipEnv ${reset}${venv}"
                export PIPENV_ACTIVE=1
                source "${venv}/bin/activate"
            fi
        fi

        return
    fi

    echo "Unknown project $1"
}

function _dad_completion {
    local cur=${COMP_WORDS[COMP_CWORD]}
    local prev=${COMP_WORDS[COMP_CWORD-1]}

    if [[ ${COMP_CWORD} -eq 1 ]]; then
        COMPREPLY=($(compgen -W "help cd" ${cur}))
        return
    fi

    if [[ ${COMP_CWORD} -ne 2 ]]; then
        COMPREPLY=()
        return
    fi

    case ${prev} in
        cd)
            cd $HOME/src/github.com/pior
            local targets=$(ls -d */ | cut -f1 -d'/' | tr '\n' ' ')
            cd $OLDPWD

            COMPREPLY=($(compgen -W "${targets}" ${cur}))
            ;;

        *)
            COMPREPLY=()
            ;;
    esac
}


function dad {
    local cmd=$1
    shift 1
    local args=$*

    case $cmd in
        cd)
            _dad_cd $args
            ;;
        up)
            _dad_up $args
            ;;
        test|t)
            _dad_test $args
            ;;
        server|s)
            _dad_server $args
            ;;
        new)
            _dad_new $args
            ;;
        *)
            _dad_help
            ;;
    esac
}

complete -F _dad_completion dad
