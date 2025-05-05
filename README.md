# tuf

`tuf` is the Terraform Un-F*ck tool (or, the Terraform Unweaver Framework)

`tuf` makes the act of splitting traditional terraform workspaces simple. 

`tuf` will:

- move HCL blocks (resources, modules, etc) from a workspace to existing or new workspaces
- remediate the state between the two workspaces
- hard-code references to module or resource outputs that existed in the source workspace, but not in the new workspace

In general, `tuf` is a framework to split up Terraform workspaces *as fast as possible* while needing near-zero intervention.

# Opinions

to make life as easy as possible, `tuf` exists with a pretty thorough list of opinions. Some of these opinions are listed below.

**Opinion: Comment relevancy + the "prettiness" of removed HCL blocks**

The default (and only) behavior is to bring along comments that are immediately connected to module/resource calls when moving them.

**Opinion: Data Block Portability**

Data blocks are automatically brought-along (in copy mode) when moving blocks between workspaces when they are deemed necessary.

**Opinion: Provider Same-ness**
It's (temporarily) implicitly assumed that provider blocks with *the same* aliases are completely interoperable.

# Disclaimer

`tuf` is in pre-release; use with caution. Opinions and functionality may change at any time.

# How it Works
`tuf` acts in three stages: `init`, `move`, `finalize`.

During the initialization phase, `tuf` will verify all connected workspaces and store potentially relevant values and metadata in a local database. 

This is done to ease (1) the size of the state files that may need to be loaded into memory during the finalize step, and (2) to keep track of multiple workspaces and dedupe common module and resource names. Ie- what is the behavior of moving a resource that uses the output `module.foo.bar`, if `module.foo.bar` is also a valid value in the destination workspace?

During the `move` phase, `tuf` works solely with the local workspaces and tracks all operations within the local migration database.

During the `finalize` phase, `tuf` will remediate all move operations according to the information stored in the db. This includes updating all related state files, replacing references to no-longer-existing outputs with their hardcoded values, fixing resource addresses that are duplicated as a result of the move (first, by attempting to append `_tuf`, then by appending a series of random characters) and doing one final verification that all initial resources and modules are accounted for in the final state.

`tuf` will *not* use a terraform/open tofu binary to run state commands by default. It will attempt to determine the version of the state file and use one of many implementations for that particular state file. The trade-off here is the speed in which you can move resources and modules vs. maintainability. Since there are not many APIs that hook into state, and the local migration database provides an abstract layer, for now this is an acceptable trade-off.

# Pre-release Checklist:
- [ ] `tuf init`
- [ ] `tuf move`
  - [x] file manipulation
  - [ ] move tracking
- [ ] `tuf finalize`

# License

MIT