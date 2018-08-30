package go_util

import (
	"html/template"
	"io"
	"net/http"
)

func WriteTemplate(w  io.Writer,templ string){
	t, err := template.New("webpage").Parse(templ)
	CheckErr(err)
	err = t.Execute(w, nil)
	CheckErr(err)
}

func RedirectUrl(writer http.ResponseWriter,url string){
	s:= `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
</body>
<script>
	location.href="`+url+`"
</script>
</html>
	`
	WriteTemplate(writer,s)
}