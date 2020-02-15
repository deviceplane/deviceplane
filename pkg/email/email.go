package email

import (
	"bytes"
	"text/template"
)

type Interface interface {
	Send(request Request) error
}

type Request struct {
	FromName    string
	FromAddress string
	ToName      string
	ToAddress   string
	Subject     string
	Content     Content
}

type Content struct {
	Title       string
	Body        string
	ActionTitle string
	ActionLink  string
}

const Template = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <link 
		href="https://fonts.googleapis.com/css?family=Rubik:400,500&display=swap"
	  	rel="stylesheet"
		type="text/css"
		media="all"
  />
  <style type="text/css" rel="stylesheet" media="all">
  	html {
		font-family: Rubik,Roboto,sans-serif;
		font-size: 18px;
  	}
    body {
		width: 100% !important;
		height: 100% !important;
		margin: 0;
		line-height: 1.4;
		font-weight: 400;
		-webkit-text-size-adjust: none;
		font-family: Rubik,Roboto,sans-serif;
	}
	table {
		background-color: #000 !important;
		color: #fff !important;
		width: 100%;
		border: none;
		border-spacing: 0;
		border-collapse: separate;
	}
	tr {
		width: 100%;
	}
	td {
		text-align: center;
		padding: 16px 0;
	}
	img {
		display: inline-block;
	}
	p {
		font-size: 16px;
		text-align: center;
		color: #C6C6C6;
	}
	.button {
		display: inline-block;
		color: #000000 !important;
		background-color: #6fccff;
		border-radius: 2px;
		padding: 12px 24px;
		text-transform: uppercase;
		font-weight: 500;
		text-decoration: none;
		font-family: Rubik,Roboto,sans-serif !important;
		font-size: 14px;
		text-align: center;
	}
	.title-row td {
		padding: 48px 0 0 0;
	}
	.title-row h1 {
		text-align: center;
		font-size: 24px;
		font-weight: 500;
		color: #ffffff;
		margin: 0;
	}
	.greeting {
		font-size: 24px;
		font-weight: 500;
	}
	.action-container {
		display: flex;
		flex: 1;
		justify-content: center;
	}
	.logo {
		margin-bottom: 24px;
	}
	.banner {
		padding-top: 24px;
	}
	.copyright {
		display: inline-block;
		font-size: 12px;
		color: #fff;
		opacity: .5;
		font-weight: 400;
		text-align: center;
		padding-top: 64px;
	}
  </style>
</head>
<body>
	<table>
		<tbody>
			<tr>
				<td>
				<a href=”” style=”cursor:default;”>
					<img height="42px" src="https://github.com/deviceplane/deviceplane/raw/master/logo/name-white.png" class="banner" />
				</a>
				</td>
			</tr>
			<tr class="title-row">
				<td>
					<h1>{{.Title}}</h1>
				</td>
			</tr>
			<tr>
				<td>
					<p>{{.Body}}</p>
				</td>
			</tr>
			<tr>
				<td>
					<a class="button" href="{{.ActionLink}}">{{.ActionTitle}}</a>
				</td>
			</tr>
			<tr>
				<td>
					<span class="copyright">
						© Deviceplane 2020
					</span>
				</td>
			</tr>
		</tbody>
	</table>
</body>
</html>
`

func GenerateHTML(request Request) (string, error) {
	templateContent := request.Content

	t, err := template.New("email").Parse(Template)

	if err != nil {
		return "", err
	}

	var templateBuffer bytes.Buffer

	if err := t.Execute(&templateBuffer, templateContent); err != nil {
		return "", err
	}

	return templateBuffer.String(), nil
}
