package errhtml

import "html/template"

var (
	errTemplate       *template.Template
	errTemplateSource = `
<html>
<head>
	<title>{{.Title}}: {{.Message}}</title>

	<style type="text/css">
		html, body {
			margin: 0;
			padding: 0;
			font-family: Helvetica, Arial, Sans;
			background: #EEEEEE;
		}

		.block {
			padding: 20px;
			border-bottom: 1px solid #aaa;
		}

		#header h1 {
			font-weight: normal;
			font-size: 28px;
			margin: 0;
		}

		#more {
			color: #666;
			font-size: 80%;
			border: none;
		}

		#header {
			background: #fcd2da;
		}

		#header p {
			color: #333;
		}

		#header .message {
			font-style: italic;
			padding-left: 10px;
		}

		.details {
			background: #f6f6f6;
		}

		.details h2 {
			font-weight: normal;
			font-size: 18px;
			margin: 0 0 10px 0;
		}

		.details .lineNumber {
			float: left;
			display: block;
			width: 40px;
			text-align: right;
			margin-right: 10px;
			font-size: 14px;
			font-family: monospace;
			background: #333;
			color: #fff;
		}

		.details .line {
			clear: both;
			color: #333;
			margin-bottom: 1px;
		}

		.details .line.trace {
			margin-bottom: 10px;
		}

		.details .line.trace .location {
			font-weight: bold;
		}

		.details .line.trace .location span.separator:after {
			content: "/";
			padding: 0 1px 0 1px;
		}
		.details .line.trace .location span.separator:last-child:after {
			content: "";
			padding: 0;
		}

		.details pre {
			font-size: 14px;
			margin: 0;
			overflow-x: hidden;
		}

		.details .line.trace pre {
			margin: 5px 0 0 25px;
		}

		.details .error {
			color: #c00 !important;
		}

		.details .error .lineNumber {
			background: #c00;
		}

		.details a {
			text-decoration: none;
		}

		.details a:hover * {
			cursor: pointer !important;
		}

		.details a:hover pre {
			background: #FAFFCF !important;
		}

		.details em {
			font-style: normal;
			text-decoration: underline;
			font-weight: bold;
		}

		.details strong {
			font-style: normal;
			font-weight: bold;
		}
	</style>
</head>

<body>
	<div id="header" class="block">
		<h1>
			{{.Title}}
		</h1>
		<p class="message">
			{{.Message}}
		</p>
	</div>

	{{if .SourceContext}}
		<div class="details block">
			{{if and .Source (not .SourceTrace)}}
				<h2>
					{{.Source.AbbreviatedFilePath}}{{if .Source.Line}}, line {{.Source.Line}}{{end}}
				</h2>
				<br/>
			{{end}}

			{{range .SourceContext}}
				<div class="line {{if .Highlight}}error{{end}}">
					<span class="lineNumber">{{.Line}}</span>
					<pre>{{.Text}}</pre>
				</div>
			{{end}}
		</div>
	{{end}}

	{{if .SourceTrace}}
		<div class="details block">
			{{range .SourceTrace}}
				<div class="line trace">
					<div class="location">
						{{range $part := .AbbreviatedFilePathDirectories}}
							<span class="directory">{{$part}}</span>
							<span class="separator"></span>
						{{end}}
						<span>{{.FileName}}:{{.Line}}</span>
					</div>
					<pre>{{.Text}}</pre>
				</div>
			{{end}}
		</div>
	{{end}}

	{{if .MetaError}}
		<div class="details block">
			<h2>Additionally, an error occurred while handling this error.</h2>
			<div class="line error">
				{{.MetaError}}
			</div>
		</div>
	{{end}}

	<script>

	</script>
</body>
</html>
`
)

func init() {
	errTemplate = template.Must(template.New("errhtml").Parse(errTemplateSource))
}
