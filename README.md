# airflow-cli

Cli to interact with you airflow instance using the REST API when you can't ssh into the airflow server.


The code as been greatly inspired by [glab](https://gitlab.com/gitlab-org/cli/-/tree/main) and [jira-cli](https://github.com/ankitpokhrel/jira-cli).

## Quickstart

### Create a new profile 

#### Using the CLI

Create a profile using :

```sh
airflow-cli profile create
```

you can choose between two kind of authentication :
- user/password
- jwt auth

In case of a jwt auth you can choose to use a shell command to retrieve the jwt token instead of directly providing it.

#### Editing the config file

The config file is stored in `~/.config/.airflow/.config`.
You can edit it directly.

```
; classic airflow authentication
[default]
user = airflow
password = airflow
url = http://0.0.0.0:8080/api/v1/

; auth across a proxy via jwt
[token]
token = secrettoken
url = https://some.server:8080/path/to/proxy

; you can retrieve the token via any shell command
[shell]
isShell = true
token = cat file_with_token
url = https://some.server:8080/path/to/proxy


```

## Commands

### Dag

The `dag` command allow you to access dags and dag runs resources. They are by default sorted by `-date_start`.

```sh
# list dags
airflow-cli dag list

# list dag run (prompt will ask for which dag_id)
airflow-cli dag list-run

# list dag run for dag_id 
airflow-cli dag list-run -d dag_id

# trigger a dag run for dag_id 
airflow-cli dag trigger -d dag_id

# get status of dag runs, tasks associated and log of a task instance
airflow-cli dag trigger -d dag_id
```


### Dag graph

To get a [graphviz](https://graphviz.org/) representation of the graph.
```sh 
airflow-cli dag graph
```

You can use `graph-easy` to plot an asci art in your terminal. Installation instruction [here](https://stackoverflow.com/questions/3211801/graphviz-and-ascii-output/3391213).

```sh
airflow-cli dag graph -d tutorial_dag | graph-easy --from=graphviz --as=boxart

              tutorial_dag

╭─────────╮     ╭───────────╮     ╭──────╮
│ extract │ ──▶ │ transform │ ──▶ │ load │
╰─────────╯     ╰───────────╯     ╰──────╯

```

### Dag run grid

Get a grid view similar to airflow UI.

```sh
go run airflow-cli/main.go dag grid -d example_bash_operator
example_bash_operator    v v v v v v
------------------------------------
also_run_this            v v v v v v
runme_0                  v v v v v v
runme_1                  v v v v v v
runme_2                  v v v v v v
run_after_loop           v v v v v v
this_will_skip           - - - - - -
run_this_last            - - - - - -
                                   |
                               2023-06-24
                               00:00:00
```


### Task

```sh
# list tasks of a dag (prompt will ask to choose a dag)
airflow-cli task list

# list tasks instance of a dag run
airflow-cli task list

# list tasks instance of a dag run with no prompt
airflow-cli task list -d dag_id -r dag_run_id

# get logs from a task instance
airflow-cli task logs -d dag_id -r dag_run_id
```


## Contribute

### Tools

To run linting, install [golang-lint](https://golangci-lint.run/usage/install/).

```sh
make lint
```

### Test

Use test makefile to spin up any airflow version.

```sh
 make up VERSION=2.5.3
 ```

## Todo 
- Choose a name
- Reoganise required module
- Wrap error and sys exit [x]
- Correct Select descriptions servey [x]
- Correctly escape with crtl+d [x]
- Write quickstart README [x]
- Output DAG with [diagon](https://github.com/ArthurSonzogni/Diagon) and [natural C biding](https://pkg.go.dev/cmd/cgo)
	- Graph-Easy : https://stackoverflow.com/a/3391213
