# プラグインAPI仕様書
<a name="top"></a>

## インデックス
- [API仕様](#API仕様)
{{range .Files}}
{{$file_name := .Name}}  - [{{.Name}}](#{{.Name}})
  {{range .Messages}}    - [{{.LongName}}](#{{.FullName}})
  {{end}}
  {{- if .Enums }}
  {{range .Enums}}    - [{{.LongName}}](#{{.FullName}})
  {{end}}
  {{ end }}
  {{- if .Extensions }}
  {{range .Extensions}}    - [File-level Extensions](#{{$file_name}}-extensions)
  {{end}}
  {{ end -}}
  {{range .Services}}    - [{{.Name}}](#{{.FullName}})
  {{end}}
{{end}}
- [スカラー値型](#スカラー値型)

## API仕様
{{range .Files}}
{{$file_name := .Name}}
<a name="{{.Name}}"></a>
<p align="right"><a href="#top">Top</a></p>

### {{.Name}}
{{.Description}}

{{range .Messages}}
<a name="{{.FullName}}"></a>

#### {{.LongName}}
{{.Description}}

{{if .HasFields}}
| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
{{range .Fields -}}
  | {{.Name}} | [{{.LongType}}](#{{.FullType}}) | {{.Label}} | {{nobr .Description}}{{if .DefaultValue}} Default: {{.DefaultValue}}{{end}} |
{{end}}
{{end}}

{{if .HasExtensions}}
| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
{{range .Extensions -}}
  | {{.Name}} | {{.LongType}} | {{.ContainingLongType}} | {{.Number}} | {{nobr .Description}}{{if .DefaultValue}} Default: {{.DefaultValue}}{{end}} |
{{end}}
{{end}}

{{end}} <!-- end messages -->

{{range .Enums}}
<a name="{{.FullName}}"></a>

#### {{.LongName}}
{{.Description}}

| Name | Number | Description |
| ---- | ------ | ----------- |
{{range .Values -}}
  | {{.Name}} | {{.Number}} | {{nobr .Description}} |
{{end}}

{{end}} <!-- end enums -->

{{if .HasExtensions}}
<a name="{{$file_name}}-extensions"></a>

#### File-level Extensions
| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
{{range .Extensions -}}
  | {{.Name}} | {{.LongType}} | {{.ContainingLongType}} | {{.Number}} | {{nobr .Description}}{{if .DefaultValue}} Default: `{{.DefaultValue}}`{{end}} |
{{end}}
{{end}} <!-- end HasExtensions -->

{{range .Services}}
<a name="{{.FullName}}"></a>

#### {{.Name}}
{{.Description}}

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
{{range .Methods -}}
  | {{.Name}} | [{{.RequestLongType}}](#{{.RequestFullType}}){{if .RequestStreaming}} stream{{end}} | [{{.ResponseLongType}}](#{{.ResponseFullType}}){{if .ResponseStreaming}} stream{{end}} | {{nobr .Description}} |
{{end}}
{{end}} <!-- end services -->

{{end}}

## スカラー値型

| .proto Type | Notes | Go Type | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | -------- | --------- | ----------- |
{{range .Scalars -}}
  | <a name="{{.ProtoType}}" /> {{.ProtoType}} | {{.Notes}} | {{.GoType}} | {{.CppType}} | {{.JavaType}} | {{.PythonType}} |
{{end}}