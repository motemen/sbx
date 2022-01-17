sbx
===

A Scrapbox CLI in the wild

Usage
-----

### project

    sbx project show <project>

Prints JSON for _project_.

Requests `https://scrapbox.io/api/projects/<project>`

### page

    sbx page list [-L <limit>] <project>

Prints JSON array of pages in _project_.

Requests `https://scrapbox.io/api/pages/<project>`


Options
-------

* `--session <session>`
  * Specify `connect.sid` cookie value on scrapbox.io, for querying private projects.
* `--jq <query>`
  * Run a jq query on results.

Configuration
-------------

Put a JSON file like below at `~/.config/sbx/config.json` to give a default value for `--session`.

```json
{
  "projects": {
    "<project>": {
      "session": {
        "command": "<command to print session>", // or
        "value": "<constant value>"
      }
    }
  },

  // default value
  "default": {
    "session": ...
  }
}
```
