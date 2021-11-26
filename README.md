# Changelog Guardian

## Getting started

Changelog Guardian is the tool that will help you to keep your Changelog safely automated and up to date.

## Installation

Use [Go Binaries](https://gobinaries.com/) to automatically compile, retrieve and install.

```bash
curl -sf https://gobinaries.com/aoterolorenzo/changelog-guardian | sh
```

## Configuration

A series of configuration parameters are internally provided by default, so no initial or custom configuration is needed to run Changelog Guardian

### Custom configurations: .cg-config.yml

The default configuration above-mentioned can be customised and overwritten on a repository scope `.cg-config.yml` file. 

This file could contain only one or various configuration parameters. Here is an example of all the parameters that forms the Changelog Guardian configuration:

```yml
changelogPath: ./CHANGELOG.md # Path to the project's Changelog File
releaseNotesPath: ./RELEASE-NOTES.md # Path to the file where to save generated Release Notes
mainBranch: main # Main branch of the repository
defaultBranch: develop # Default/develop branch of the repository
providers: # Internal configuration of the Changelog Guardian providers
  gitlab:
    labels: # Relation between Gitlab labels and tasks sections. 
      Breaking Changes: kind::breaking
      Added: kind::feature
      Changed: kind::perf
      Refactor: kind::refactor
      Fixed: kind::bugfix
      Dependencies: kind::dependencies
      Deprecated: kind::deprecation
      Removed: kind::removal
      Documentation: kind::docs
      Security: kind::security
releaseProvider: git # Release provider
tasksProvider: git # Tasks provider
style: markdown # Changelog style (theming)
releasePipes: [ 'semver' ] # Release pipes
taskPipes: [] # Task pipes
initialVersion: 0.1.0 # Initial version for generating an initial release
stylesCfg:
  stylish_markdown:
    categories: // Selects a color and an emoji for Categories on template generated CHANGELOG's
      Breaking Changes: ['f70000','🚨']
      Added: ['5ccb31','✨']
      Changed: ['31cb7d','✒️']
      Refactor: ['cba531','🏗']
      Fixed: ['cb3131','🐛']
      Dependencies: ['cb6b31','📦']
      Deprecated: ['4e31cb','✖️']
      Removed: ['7631cb','❌']
      Documentation: ['3188cb','📖']
      Security: ['b841a0','🔒']
```

### Providers

Changelog Guardian implements different providers for both releases and tasks. This means that you can choose from where to get the versioning of your project.

¿Do you want to retrieve the releases from the repository git tags, the Gitlab releases, or maybe from an Artifactory or Registry (pending implementation)?  ¿Do you want to retrieve your tasks from git commits, from Gitlab's Merge Requests, from Github Pull request (pending)...? It's up to you.

#### Git release provider

Obtains the releases and the project versioning from the git repository tags.

#### Git task provider

Each commit is a task. This seems for ninjas, but can be very optimized with task Pipes

#### Gitlab release provider

Releases are obtained from the [Gitlab Releases](https://about.gitlab.com/releases/categories/releases/) from your project.

* NOTE: Gitlab automatically detects your repo from the git remotes. No further configuration needed. 
* NOTE: If you need to access to a _private repository_, just set the environment variable `GITLAB_TOKEN` with a Personal Access Token

#### Gitlab task provider

Fetch already merged Merge Requests as tasks. Remember that you can customize your configuration to modify the labels you will use.

### Pipes

Changelog Guardian Pipes are little fragments of code that filter the releases and the tasks as you need.

The pipes can be combined and Changelog Guardian will make each release/task go through them in an ordered way.

For example, using the `gitlab_resolver` and the `natural_language` task pipes simultaneously

```yml
taskPipes: [ 'gitlab_resolver', 'natural_language']
```

will result in something like:

`Resolve "Add new provider"` -> `Add new provider` -> `Added new provider`

#### Semver release Pipe

Filter the releases and allows only those which matches Semantic Versioning nomenclatures.

Pipe code: `semver`

#### Gitlab Resolver task Pipe

Automatic Gitlab Merge Request merge with a message of the kind `Resolve "Issue tittle"`. This pipe modifies the task title to match only `Issue title`

For example: `Resolve "Add new provider"` -> `Add new provider`

#### Natural Language task Pipe

Usually we name our tasks on an infinitive voice: `Add new feature`. This pipe will filter the task names and will replace some of the most frequently used verbs with a past voice.

Some examples: 

`Add new provider` -> `Added new provider`
`Fix new provider` -> `Fixed new provider`
`Refactor new provider` -> `Refactorized new provider`

#### Conventional Commits task pipe (WIP)

This pipe will filter your tasks (mostly for using with the Git task provider) following [Conventional Commits](https://www.conventionalcommits.org/) specification, and appending them to your changelog sections depending on the commit type and scope.


### Changelog Styles

Maybe you would like to maintain a Changelog fed with some emojis. Or perhaps you want to generate a CHANGELOG.adoc in asciidoc. Changelog Styles are meant for exactly cover this ~~whims~~ needs.

#### Markdown Style

Default Changelog Style. 
Follows the [Keep A Changelog](https://keepachangelog.com/en/1.1.0/#how) specification.

Style code: `markdown`

#### Stylish Markdown Style

Markdown Style fed with some emojis and badges. 
Follows the [Keep A Changelog](https://keepachangelog.com/en/1.1.0/#how) specification.

Style code: `stylish_markdown`

## Usage

The usage is pretty straight forward, just use `changelog-guardian` from your terminal.


```bash
$> changelog-guardian [command]
```

```
Usage:
  changelog-guardian [flags]
  changelog-guardian [command]

Available Commands:
  help        Help about any command
  insert      Inserts a task in CHANGELOG
  release     Generates a new Release
  yank        Yank release

Flags:
      --config string    CLI config file
      -h, --help        help for changelog-guardian
```

### Generating a changelog

You can generate a changelog with the base command from a folder containing a git repository. This will generate a `CHANGELOG.md` file automatically with the default settings.

```bash
$> changelog-guardian
```

```
Usage:
  changelog-guardian [flags]

Flags:
 
  --template             CHANGELOG template
  
  -h, --help             Prints help
```

### Release

After ensure and check all the tasks to release, automatically calculates (following semver) and releases a new version 

```bash
$> changelog-guardian release
```

```
Usage:
  changelog-guardian release [flags]

Flags:
  
  -M, --major            Major Release
  -m, --minor            Minor Release
  -p, --patch            Patch Release
  -v, --version string   Specific version to release
  -f, --force            Forces the versioning altough differs from the calculated one
      --pre string       Pre-release string (semver)
      --build string     Build metadata (semver)
  --template             CHANGELOG template
  
  -h, --help             Prints help
```


### Insert

Inserts a task in the `Unreleased` section, and checks all the data against the selected **tasks provider**

```bash
$> changelog-guardian insert
```

```
Usage:
  changelog-guardian insert [flags]

Flags:
  -i, --id string             Task ID
  -t, --title string          Task title
    -l, --link string           Task link
  -f, --author string         Task author
  -v, --authorLink string     Task author link
  -c, --category string       Task category (default "Added")
  -s, --skip-autocompletion   Skip autocompletion from providerUsed to check the task data from it through the provided --id
  --template             CHANGELOG template
  
  -h, --help                  Prints help
```

### Yank

Yanks a version and move all the tasks it contained to the immediately superior (or to the `Unreleased` section if the yanked version is the last). Changelog Guardian will yank the last release if no flags are provided

```bash
$> changelog-guardian yank
```

```
  changelog-guardian yank [flags]

Flags:
  -v, --version string   Version to yank
  --template             CHANGELOG template
  
  -h, --help             Prints help
```

### Yank

Generates the Release Notes for the last released version (or the one specified by the --version flag).

```bash
$> changelog-guardian yank
```

```
  changelog-guardian yank [flags]

Flags:
  -v, --version string      Version to yank
  -o, --output-file string   Output file
  -e, --echo                Echo Release Notes on screen
  --template                CHANGELOG template
  
  -h, --help             Prints help
```