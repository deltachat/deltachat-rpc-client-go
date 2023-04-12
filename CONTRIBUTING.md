## Contributing

After doing your modifications make sure all tests pass, and add tests to ensure code coverage of your
new modifications.

### Running the test suite

To run the integration tests, you need to have a local fake email server running.
The easiest way to do that is with Docker:

```
$ docker pull ghcr.io/deltachat/mail-server-tester:release
$ docker run -it --rm -p 3025:25 -p 3110:110 -p 3143:143 -p 3465:465 -p 3993:993 ghcr.io/deltachat/mail-server-tester
```

Leave the fake email server running, open a new shell and run:

```
$ ./scripts/run_tests.sh
```
The `run_tests.sh` script will install `deltachat-rpc-server` (if needed) and run all tests.

To run all `Account` tests:

```
go test -v ./... -run TestAccount
```

To run a single test, for example `TestChat_SetName`:

```
go test -v ./... -run TestChat_SetName
```
