### summary
Use the GitHub CLI to add a "Created Date" column to a GitHub project.

### developer loop

1. Install the gh cli
1. Login `gh auth refresh` and follow the steps
1. Ensure you add the project scope `gh auth refresh -s project`

```bash
make build
make test
make clean
```

### usage

1. find the fully-qualified url to your project.
2. run the following script with to dry-run the operation:

```bash
createdat --project=https://github.com/orgs/<org>/project/<number>
```

3. When you're happy with the output run with `--debug=false`
