# glint

A Go linter for the [G framework](https://github.com/enetx/g).

## Installation

```bash
go install github.com/enetx/glint/cmd/glint@latest
```

## Vim ALE Integration

```lua
return {
    "dense-analysis/ale",

    config = function()
        vim.g.ale_disable_lsp = 1
        vim.g.ale_echo_msg_format = "[%linter%] %s [%severity%]"
        vim.g.ale_set_highlights = 0
        vim.g.ale_set_signs = 0
        vim.g.ale_virtualtext_prefix = "█ "

        vim.g.ale_linters = {
            go = { "glint" },
        }

        vim.cmd([[
            call ale#linter#Define('go', {
            \   'name': 'glint',
            \   'executable': 'glint',
            \   'command': 'glint %t',
            \   'output_stream': 'stdout',
            \   'callback': 'ale#handlers#unix#HandleAsWarning',
            \   'lint_file': 1,
            \})
        ]])
    end,
}
```
