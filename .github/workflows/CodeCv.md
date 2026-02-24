            - name: Codecov
  # You may pin to the exact commit or the version.
  # uses: codecov/codecov-action@671740ac38dd9b0130fbe1cec585b89eea48d3de
  uses: codecov/codecov-action@v5.5.2
  with:
    # The base SHA to select. This is only used in the "pr-base-picking" run command
    base_sha: # optional
    # The file location of a pre-downloaded version of the CLI. If specified, integrity checking will be bypassed.
    binary: # optional
    # The location of the codecov.yml file. This is crrently ONLY used for automated test selection (https://docs.codecov.com/docs/getting-started-with-ats). Note that for all other cases, the Codecov yaml will need to be located as described here: https://docs.codecov.com/docs/codecov-yaml#can-i-name-the-file-codecovyml
    codecov_yml_path: # optional
    # SHA (with 40 chars) of what should be the parent of this commit.
    commit_parent: # optional
    # Folder to search for coverage files. Default to the current working directory
    directory: # optional
    # Disable file fixes to ignore common lines from coverage (e.g. blank lines or empty brackets). Read more here https://docs.codecov.com/docs/fixing-reports
    disable_file_fixes: # optional, default is false
    # Disable search for coverage files. This is helpful when specifying what files you want to upload with the files option.
    disable_search: # optional, default is false
    # Disable setting safe directory. Set to true to disable.
    disable_safe_directory: # optional, default is false
    # Disable sending telemetry data to Codecov. Set to true to disable.
    disable_telem: # optional, default is false
    # Don't upload files to Codecov
    dry_run: # optional, default is false
    # Environment variables to tag the upload with (e.g. PYTHON | OS,PYTHON)
    env_vars: # optional
    # Comma-separated list of folders to exclude from search.
    exclude: # optional
    # On error, exit with non-zero code
    fail_ci_if_error: # optional, default is false
    # Comma-separated list of explicit files to upload. These will be added to the coverage files found for upload. If you wish to only upload the specified files, please consider using disable_search to disable uploading other files.
    files: # optional
    # Comma-separated list of flags to upload to group coverage metrics.
    flags: # optional
    # Only used for empty-upload run command
    force: # optional
    # Override the git_service (e.g. github_enterprise)
    git_service: # optional, default is github
    # Extra arguments to pass to gcov
    gcov_args: # optional
    # gcov executable to run. Defaults to 'gcov'
    gcov_executable: # optional, default is gcov
    # Paths to ignore during gcov gathering
    gcov_ignore: # optional
    # Paths to include during gcov gathering
    gcov_include: # optional
    # If no coverage reports are found, do not raise an exception.
    handle_no_reports_found: # optional, default is false
    #
    job_code: # optional
    # Custom defined name of the upload. Visible in the Codecov UI
    name: # optional
    # Specify a filter on the files listed in the network section of the Codecov report. This will only add files whose path begin with the specified filter. Useful for upload-specific path fixing.
    network_filter: # optional
    # Specify a prefix on files listed in the network section of the Codecov report. Useful to help resolve path fixing.
    network_prefix: # optional
    # Override the assumed OS. Options available at cli.codecov.io
    os: # optional
    # Specify the branch to be displayed with this commit on Codecov
    override_branch: # optional
    # Specify the build number manually
    override_build: # optional
    # The URL of the build where this is running
    override_build_url: # optional
    # Commit SHA (with 40 chars)
    override_commit: # optional
    # Specify the pull request number manually. Used to override pre-existing CI environment variables.
    override_pr: # optional
    # Comma-separated list of plugins to run. Specify `noop` to turn off all plugins
    plugins: # optional
    # Whether to enumerate files inside of submodules for path-fixing purposes. Off by default.
    recurse_submodules: # optional, default is false
    # The code of the report if using local upload. If unsure, leave default. Read more here https://docs.codecov.com/docs/the-codecov-cli#how-to-use-local-upload
    report_code: # optional
    # The type of file to upload, coverage by default. Possible values are "test_results", "coverage".
    report_type: # optional
    # Root folder from which to consider paths on the network section. Defaults to current working directory.
    root_dir: # optional
    # Choose which CLI command to run. Options are "upload-coverage", "empty-upload", "pr-base-picking", "send-notifications". "upload-coverage" is run by default.
    run_command: # optional, default is upload-coverage
    # Skip integrity checking of the CLI. This is NOT recommended.
    skip_validation: # optional, default is false
    # [Required when using the org token] Set to the owner/repo slug used instead of the private repo token. Only applicable to some Enterprise users.
    slug: # optional
    # Specify the swift project name. Useful for optimization.
    swift_project: # optional
    # Repository Codecov token. Used to authorize report uploads
    token: # optional
    # Set to the Codecov instance URl. Used by Dedicated Enterprise Cloud customers.
    url: # optional
    # Use the legacy upload endpoint.
    use_legacy_upload_endpoint: # optional, default is false
    # Use OIDC instead of token. This will ignore any token supplied
    use_oidc: # optional, default is false
    # Use the pypi version of the CLI instead of from cli.codecov.io
    use_pypi: # optional, default is false
    # Enable verbose logging
    verbose: # optional, default is false
    # Which version of the Codecov CLI to use (defaults to 'latest')
    version: # optional, default is latest
    # Directory in which to execute codecov.sh
    working-directory: # optional
