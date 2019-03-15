// Copyright © 2019 James Rivett-Carnac
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

const template = `# Background

Enter the background to his ADR here.  Setup the context
so we are on the same page.  Keep it simple and easy to
follow. Don't tell me the problem

# Complication

Now tell me where the problem/complication is that this ADR
is addressing.

# Options Considered

1. This was one
2. This was another

# Decision

What did we decided

# Outcome
`

const baseTemplate = `# {{ .Title }}

*Number:* {{ .Filename }}
*Created:* {{ .Created }}
*Status:* {{ .Status }}
*Authors:*
{{ range $a := .Authors -}}
- {{ $a }}
{{- end -}}
{{ if .Related }}
*Related*:
{{ range $r := .Related -}}
- {{ $r }}
{{- end -}}
{{- end}}

{{ .Text }}
`
