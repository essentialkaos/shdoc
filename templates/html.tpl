<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset='utf-8'>

    <title>{{ .Title }}</title>

    <link href='https://fonts.googleapis.com/css?family=Roboto:400,300,700|Roboto+Mono' rel='stylesheet' type='text/css'>

    <style type="text/css">
      html,body { color:#222; font-family:Roboto, Verdana, sans-serif; height:100%; margin:0; padding:0 }
      h1,h2,h3 { color:#666; font-weight:100; margin:0; padding:32px 0 8px }
      h1 { border-bottom:1px #DDD solid; font-size:2.2em; padding-bottom:8px }
      h2 { font-size:1.6em }
      h3 { font-size:1.4em }
      code,.mono { font-family:'Roboto Mono', monospace }
      a { color:#222; text-decoration:none }
      p,div { position:relative }
      div.doc { display:block; font-size:.9em; margin-left:auto; margin-right:auto; padding-top:32px; width:800px }
      div.toc { margin:0; padding-top:8px }
      div.entity { margin:0; padding-top:16px }
      div.method { margin:0; padding-top:48px }
      div.entity::before,div.toc::before,div.method::before { color:#AAA; content:attr(data-loc); font-size:.9em; margin-right:12px; margin-top:4px; position:absolute; right:100% }
      div.method-data { margin-left:24px }
      div.argument { padding-top:2px }
      div.arguments,div.result,div.example { padding-top:16px }
      div.example-code { background-color:#f5f5f5; border:1px solid #CCC; border-radius:4px; color:#444; font-size:.9em; margin-top:8px; padding:16px; white-space:pre-wrap }
      span.badge { border-radius:4px; color:#FFF; cursor:default; font-size:.6em; font-weight:700; padding:2px 4px; vertical-align:middle }
      span.number { background-color:#DEAF57 }
      span.string { background-color:#5598E2 }
      span.boolean { background-color:#50C449 }
      span.dot { font-weight:700; vertical-align:middle }
      span.dot-number { color:#DEAF57 }
      span.dot-string { color:#5598E2 }
      span.dot-boolean { color:#50C449 }
      span.desc { color:#444 }
      span.variable { font-size:.9em }
      span.optional { background-color:#BBB }
      div.footer { color:#999; font-size:.9em; padding:64px 0 40px; text-align:center }
      div.footer a { border-bottom:1px solid #666; color:#666 }
      span.equals,span.title { color:#888 }
    </style>
  </head>
  <body>
    <div class="doc">
      <h1>{{ .Title }}</h1>
      {{ if .HasAbout }}
      <h2>About</h2>
      <div>
        {{ range .About }}<p>{{ . }}</p>{{ end }}
      </div>
      {{ end }}

      <!-- TOC -->

      {{ if .HasConstants }}
      <h3>Constants</h3>
      {{ range .Constants }}
      <div data-loc="{{ .Line }}" class="toc"><a class="mono" href="#{{ .Line }}">{{ .Name }}</a> <span class="dot dot-{{ .TypeName 1 }}">•</span></div>
      {{ end }}
      {{ end }}

      {{ if .HasVariables }}
      <h3>Global Variables</h3>
      {{ range .Variables }}
      <div data-loc="{{ .Line }}" class="toc"><a class="mono" href="#{{ .Line }}">{{ .Name }}</a> <span class="dot dot-{{ .TypeName 1 }}">•</span></div>
      {{ end }}
      {{ end }}

      {{ if .HasMethods }}
      <h3>Methods</h3>
      {{ range .Methods }}
      <div data-loc="{{ .Line }}" class="toc"><a class="mono" href="#{{ .Line }}">{{ .Name }}</a></div>
      {{ end }}
      {{ end }}

      <!-- CONSTANTS -->

      {{ if .HasConstants }}
      <h2>Constants</h2>
      {{ range .Constants }}
      <div data-loc="{{ .Line }}" id="{{ .Line }}" class="entity">
        <div>
          <a class="mono" href="#{{ .Line }}">{{ .Name }}</a> <span class="equals">=</span> <span class="code">{{ .Value }}</span> <span class="badge {{ .TypeName 1 }}">{{ .TypeName 2 }}</span>
        </div>
        <div>
          <span class="variable desc">{{ .UnitedDesc }}</span>
        </div>
      </div>
      {{ end }}
      {{ end }}

      <!-- VARIABLES -->

      {{ if .HasVariables }}
      <h2>Global Variables</h2>
      {{ range .Variables }}
      <div data-loc="{{ .Line }}" id="{{ .Line }}" class="entity">
        <div>
          <a class="mono" href="#{{ .Line }}">{{ .Name }}</a> <span class="equals">=</span> <span class="mono">{{ .Value }}</span> <span class="badge {{ .TypeName 1 }}">{{ .TypeName 2 }}</span>
        </div>
        <div>
          <span class="variable desc">{{ .UnitedDesc }}</span>
        </div>
      </div>
      {{ end }}
      {{ end }}

      <!-- METHODS -->

      {{ if .HasMethods }}
      <h2>Methods</h2>
      {{ range .Methods }}
      <div data-loc="{{ .Line }}" id="{{ .Line }}" class="method">
        <div>
          <a class="mono" href="#{{ .Line }}">{{ .Name }}</a><span class="desc"> — {{ .UnitedDesc }}</span>
        </div>
        <div class="method-data">
          {{ if .HasArguments }}
          <div class="arguments">
            {{ range .Arguments }}
            <div class="argument">
              <span class="variable title">{{ .Index }}.</span> <span class="variable desc">{{ .Desc }}</span> <span class="badge {{ .TypeName 1 }}">{{ .TypeName 2 }}</span> {{ if .IsOptional }}<span class="badge optional">OPTIONAL</span>{{ end }}
            </div>
            {{ end }}
          </div>
          {{ end }}
          {{ if .ResultCode }}
          <div class="result">
            <span class="variable title">Code:</span> <span class="variable desc">0 - ok, 1 - not ok</span>
          </div>
          {{ end }}
          {{ if .HasEcho }}
          <div class="result">
            <span class="variable title">Echo:</span> <span class="variable desc">{{ .ResultEcho.UnitedDesc }}</span> <span class="badge {{ .ResultEcho.TypeName 1 }}">{{ .ResultEcho.TypeName 2 }}</span>
          </div>
          {{ end }}
          {{ if .HasExample }}
          <div class="example">
            <span class="variable title">Example:</span>
            <div class="example-code">{{ range .Example }}<code>{{ . }}<br/></code>{{ end }}</div>
          </div>
          {{ end }}
        </div>
      </div>
      {{ end }}
      {{ end }}
    </div>

    <!-- FOOTER -->

    <div class="footer">Generated with ❤ by <a href="https://kaos.sh/shdoc">SHDoc</a></div>
  </body>
</html>
