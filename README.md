# Changelog Guardian

## Getting started

Changelog Guardian is the tool that will help you to keep your Changelog safely automated and up to date.

## Installation

We provide a INSTALL.sh file that identifies your SO and architecture, and directly downloads and installs the specific `changelog-guardian` binary into your /usr/local/bin for the latest version available.

```bash
curl -sf https://gitlab.com/aoterocom/changelog-guardian/-/raw/main/INSTALL.sh | sh
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
  github:
    labels:
      Breaking Changes: breaking
      Added: feature
      Changed: perf
      Refactor: refactor
      Fixed: bugfix
      Dependencies: dependencies
      Deprecated: deprecation
      Removed: removal
      Documentation: docs
      Security: security
releaseProvider: git # Release provider
tasksProvider: git # Tasks provider
template: markdown # Changelog template (theming)
releasePipes: [ 'semver' ] # Release pipes
tasksPipes: [] # Task pipes
tasksPipesCfg:
  conventional_commits:
    categories: # Associates Categories with Conventional Commits types
      Breaking Changes: breaking
      Added: feat
      Changed: perf
      Refactor: refactor
      Fixed: fix
      Removed: revert
      Documentation: docs
  inclusions_exclusions:
    labels:
      excluded: [ 'internal' ]
      included: [ '*all' ]
    paths:
      excluded: [ ]
      included: [ '*all' ]
  jira:
    regex: "\\[(?P<key>[A-Z]{2,}[-]{1,}\\d+)]"
    baseUrl: "https://jira.atlassian.net/"
    labels:
      Breaking Changes: breaking
      Added: feature
      Changed: perf
      Refactor: refactor
      Fixed: bugfix
      Dependencies: dependencies
      Deprecated: deprecation
      Removed: removal
      Documentation: docs
      Security: security
initialVersion: 0.1.0 # Initial version for generating an initial release
templatesCfg:
  stylish_markdown:
    categories: # Selects a color and an emoji for Categories on template generated CHANGELOG's
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

Also, Changelog Guardian provides some methods to revert a task already in the CHANGELOG _(see each task provider description above)_

#### Git release provider

Obtains the releases and the project versioning from the git repository tags.

#### Git tasks provider

Each commit is a task. This seems for ninjas, but can be very optimized with task Pipes

_Revert: Git commit message `Revert "Commit message from commit to revert"`_

#### Gitlab release provider

Releases are obtained from the [Gitlab Releases](https://about.gitlab.com/releases/categories/releases/) from your project.

* NOTE: Gitlab automatically detects your repo from the git remotes. No further configuration needed.
* NOTE: If you need to access to a _private repository_, just set the environment variable `GITLAB_TOKEN` with a Personal Access Token

#### Gitlab tasks provider

Fetches already merged Merge Requests as tasks. Remember that you can customize your configuration to modify the labels you will use.

_Revert: Using the REMOVE label or making a Merge Request with the title `Revert "Title from Merge Request to revert"`_

#### GitHub release provider

Releases are obtained from the [Github Releases](https://docs.github.com/es/repositories/releasing-projects-on-github/about-releases) from your project.

* NOTE: GitHub automatically detects your repo from the git remotes. No further configuration needed.
* NOTE: You will need to set the environment variable `GITHUB_TOKEN` with a Personal Access Token.

#### GitHub tasks provider

As the GitHub tasks provider, it fetches already merged Merge Requests as tasks. Remember that you can customize your configuration to modify the labels you will use.

_Revert: Using the REMOVE label or making a Pull Request with the title `Revert "Title from Pull Request to revert"`_

### Pipes

Changelog Guardian Pipes are little fragments of code that filter the releases and the tasks as you need.

The pipes can be combined and Changelog Guardian will make each release/task go through them in an ordered way.

For example, using the `gitlab_resolver` and the `natural_language` task pipes simultaneously

```yml
tasksPipes: [ 'gitlab_resolver', 'natural_language']
```

will result in something like:

`Resolve "Add new provider"` -> `Add new provider` -> `Added new provider`

#### Semver release Pipe

Filter the releases and allows only those which matches Semantic Versioning nomenclatures.

Pipe code: `semver`

#### Gitlab Resolver tasks Pipe

Automatic Gitlab Merge Request merge with a message of the kind `Resolve "Issue tittle"`. This pipe modifies the task title to match only `Issue title`

For example: `Resolve "Add new provider"` -> `Add new provider`

#### Jira tasks Pipe

Detect Jira tickets in your tasks and grab all the info from a pre-configured Jira endpoint

For example: `Working in jira task [JIRAKEY-01]` -> `Title of the Jira ticket`

**The tasks categories, author, links... will be replaced with the info grabbed direclty from the Jira ticket**

_Revert: when the pipe incoming task is marked as a REMOVAL (revert), the Jira tasks Pipe will keep that category with the parsed Jira task_

#### Natural Language tasks Pipe

Usually we name our tasks on an infinitive voice: `Add new feature`. This pipe will filter the task names and will replace some of the most frequently used verbs with a past voice.

Some examples:

`Add new provider` -> `Added new provider`
`Fix new provider` -> `Fixed new provider`
`Refactor new provider` -> `Refactorized new provider`

#### Inclusions&Exclusions tasks pipe

Need to exclude or include only certain labels or paths that tasks must address? With this pipe you can do it

Take a look at the custom configuration section to see the specific configuration to make it happen.
PS: For inclusions, you can use the `*all` wildcard to allow all files/paths.

#### Conventional Commits task pipe

This pipe filters your tasks (mostly for using with the Git task provider) following [Conventional Commits](https://www.conventionalcommits.org/) specification, and appending them to your changelog sections depending on the commit type.


### Changelog Templates

Maybe you would like to maintain a Changelog fed with some emojis. Or perhaps you want to generate a CHANGELOG.adoc in asciidoc. Changelog Templates are meant for exactly cover this ~~whims~~ needs.

#### Markdown Template

Default Changelog Template. 
Follows the [Keep A Changelog](https://keepachangelog.com/en/1.1.0/#how) specification.

Template code: `markdown`

#### Stylish Markdown Template

Markdown Template fed with some emojis and badges. 
Follows the [Keep A Changelog](https://keepachangelog.com/en/1.1.0/#how) specification.

Template code: `stylish_markdown`

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

  Global Flags:
      --changelog-path string    CHANGELOG path
      --config string            Config file path
      --output-template string   Output CHANGELOG template
      --silent                   Logging level
      --template string          CHANGELOG template
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

### Release Notes

Generates the Release Notes for the last released version (or the one specified by the --version flag).

```bash
$> changelog-guardian release-notes
```

```
  Usage:
  Changelog release-notes [flags]

Flags:
  -e, --echo                 Echo Release Notes on screen
  -h, --help                 Prints help
  -o, --output-file string   Output file
  -v, --version string       Version
```

### Calculate Release

Calculates the next release version without makin further changes to the CHANGELOG nor other files and prints it on the screen

```bash
$> changelog-guardian release-notes
```

```
  Usage:
  Changelog calculate-release [flags]

Flags:
      --build string   Build metadata (semver)
  -h, --help           help for calculate-release
  -M, --major          Major Release
  -m, --minor          Minor Release
  -p, --patch          Patch Release
      --pre string     Pre-release string (semver)
```

