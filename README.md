# airflow-cli

Cli to interact with you airflow instance using the REST API when you can't ssh into the airflow server.


The code as been greatly inspired by [glab](https://gitlab.com/gitlab-org/cli/-/tree/main) and [jira-cli](https://github.com/ankitpokhrel/jira-cli).

## Todo 
- Choose a name
- Reoganise required module [x]
- testing
- Wrap error and sys exit [x]
- Correct Select descriptions servey [x]
- Correctly escape with crtl+d [x]
- Write quickstart README [x]
- Output DAG with [diagon](https://github.com/ArthurSonzogni/Diagon) and [natural C biding](https://pkg.go.dev/cmd/cgo)
	- Graph-Easy : https://stackoverflow.com/a/3391213

## Quickstart

### Create a new profile 

Create a profile using :

```sh
airflow-cli profile create
```

you can choose between two kind of authentication :
- user/password
- jwt auth

In case of a jwt auth you can choose to use a shell command to retrieve the jwt token instead of directly providing it.

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

#### Airflow server

Use test makefile to spin up any airflow version.

```sh
 make up VERSION=2.5.3
 ```

#### Golden files

This repository test the behavior of the CLI.
It stores on disk (so called golden file) the expected output of the command under test.

When the behaviour of a file change, it's good practice to update the golden file:
```sh
go test integration/cli_test.go -update
```

Largely inspired from [lucapette](https://lucapette.me/writing/writing-integration-tests-for-a-go-cli-application/) and [sobyte](https://www.sobyte.net/post/2022-07/go-setup-and-teardown/).

