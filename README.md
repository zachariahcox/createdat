### summary
Use the GitHub CLI to add a "Created Date" column to a GitHub project.

### developer loop

1. Install the gh cli
1. Login `gh auth refresh` and follow the steps
1. Ensure you add the project scope `gh auth -s project`

```bash
make build
make test
make clean
```

### usage

provide the fully-qualified url to your project.

```bash
createdat --debug=true --project=https://github.com/orgs/<org>/project/<number>
```