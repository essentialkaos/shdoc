<!DOCTYPE html>
<html>
  <head>
    <meta charset='utf-8'>
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta http-equiv="Content-Language" content="en">

    <title>{{ .Title }}</title>

    <link href='https://fonts.googleapis.com/css?family=Roboto:400,300,700|Roboto+Mono' rel='stylesheet' type='text/css'>

    <style type="text/css">
      html, body {font-family: 'Roboto', Verdana, sans-serif; height:100%; padding:0; margin:0; color: #222 }
      h1, h2, h3, h4, h5, h6 { font-weight: 100; color: #666; padding-top: 16px }
      h1 { border-bottom: 1px #DDD solid; padding-bottom: 8px }
      code, #code { font-family: 'Roboto Mono', monospace }
      a { color: #222; text-decoration: none }
      p { position: relative; }
      p::before { content: attr(data-loc); position: absolute; right: 100%; margin-right: 12px; margin-top: 2px; font-size: 80%; color: #AAA }
      #doc { width:800px; display:block; margin-left:auto; margin-right:auto; padding:16px 56px 96px }
      .badge { font-size: 60%; font-weight: 700; vertical-align: middle; color: #FFF; padding: 2px 4px 2px 4px; border-radius: 4px }
      .number { background-color: #DEAF57 }
      .string { background-color: #5598E2 }
      .boolean { background-color: #50C449 }
      .equals { color: #888 }
      .desc { color: #444 }
      .variable { font-size: 80%; }
      .title { margin-left: 24px; color: #888 }
      .optional { background-color: #BBB }
      .example { font-size: 80%; color: #444; margin-left: 24px; padding: 16px; border: 1px solid #CCC; border-radius: 4px; background-color: #f5f5f5; white-space:pre-wrap }
      #footer { font-size: 80%; text-align: center; color: #999 }
      #footer a { color: #666 }
    </style>
  </head>
  <body>
    <div id="doc">
      <h1>{{ .Title }}</h1>
      {{ if .HasAbout }}
      <h2>About</h3>
      {{ range .About }}<p>{{ . }}</p>{{ end }}
      {{ end }}
      
      {{ if .HasConstants }}
      <h2>Constants</h2>
      {{ range .Constants }}
      <p data-loc="{{ .Line }}">
        <span class="code"><a name="{{ .Line }}" href="#{{ .Line }}">{{ .Name }}</a></span> <span class="code equals">=</span> <span class="code">{{ .Value }}</span> <span class="badge {{ .TypeName 1 }}">{{ .TypeName 2 }}</span><br/>
        <span class="variable desc">{{ .UnitedDesc }}</span>
      </p>
      {{ end }}
      {{ end }}

      {{ if .HasVariables }}
      <h2>Global Variables</h2>
      {{ range .Variables }}
      <p data-loc="{{ .Line }}">
        <span class="code"><a name="{{ .Line }}" href="#{{ .Line }}">{{ .Name }}</a></span> <span class="code equals">=</span> <span class="code">{{ .Value }}</span> <span class="badge {{ .TypeName 1 }}">{{ .TypeName 2 }}</span><br/>
        <span class="variable desc">{{ .UnitedDesc }}</span>
      </p>
      {{ end }}
      {{ end }}

      {{ if .HasMethods }}
      <h2>Methods</h2>
      {{ range .Methods }}
      <p data-loc="{{ .Line }}">
        <span class="code"><a name="{{ .Line }}" href="#{{ .Line }}">{{ .Name }}</a></span><span class="desc"> - {{ .UnitedDesc }}</span><br/>
  
        {{ if .HasArguments }}
        <br/>
        {{ range .Arguments }}<span class="variable title">{{ .Index }}.</span> <span class="variable desc">{{ .Desc }}</span> <span class="badge {{ .TypeName 1 }}">{{ .TypeName 2 }}</span> {{ if .IsOptional }}<span class="badge optional">OPTIONAL</span>{{ end }} <br/>{{ end }}
        {{ end }}

        {{ if .ResultCode }}<br/><span class="variable title">Code:</span> <span class="variable desc">0 - ok, 1 - not ok</span><br/>{{ end }}
        {{ if .HasEcho }}<br/><span class="variable title">Echo:</span> <span class="variable desc">{{ .ResultEcho.UnitedDesc }}</span> <span class="badge {{ .ResultEcho.TypeName 1 }}">{{ .ResultEcho.TypeName 2 }}</span><br/>{{ end }}
        {{ if .HasExample }}
        <br/><span class="variable title">Example:</span>
        <div class="example">{{ range .Example }}<code>{{ . }}<br/></code>{{ end }}</div>
        {{ end }}
        <br/>
      </p>
      {{ end }}
      {{ end }}
    </div>
    <p id="footer">Genereated by <a href="https://github.com/essentialkaos/shdoc">SHDoc</a><br/><br/><p>
  </body>
</html>
