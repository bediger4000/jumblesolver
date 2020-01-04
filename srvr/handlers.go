package srvr

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"jumble/dictionary"
	"jumble/solver"
)

var indexHTML string = `
<html>
<head>
</head>
<body>
<form name="f" method="post" action="/form">
<input name="word" />
<input type="submit" />
</form>
</body>
</html>
`

var formHTML string = `
<html>
<head>
</head>
<body>
<form name="f" method="post" action="/form">
<p>%s</p>
<input name="word" />
<input type="submit" />
</form>
</body>
</html>
`

var explainHTML string = `<!DOCTYPE html>
<html>
    <head>
    <meta charset="UTF-8">
	</head>
	<body>
`

var errorHTML string = `<!DOCTYPE html>
<html>
    <head>
    <meta charset="UTF-8">
	</head>
	<body>
		<p>I had a problem.</p>
		<p>%s</p>
	</body>
</html>
`

func (s *Srvr) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.Debug {
			fmt.Printf("Enter handleIndex closure\n")
			defer fmt.Printf("Exit handleIndex closure\n")
		}
		w.Write([]byte(indexHTML))
	}
}

func (s *Srvr) handleJumble() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.Debug {
			fmt.Printf("Enter handleJumble closure\n")
			defer fmt.Printf("Exit handleJumble closure\n")
		}

		w.Header().Set("Content-Type", "text/html")

		words, err := readRequestData(r)

		if err != nil {
			w.Write([]byte(fmt.Sprintf(errorHTML, err)))
			return
		}

		rewriteHTML(words, w)
	}
}

func (s *Srvr) handleForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Enter handleForm closure\n")
		defer fmt.Printf("Exit handleForm closure\n")

		w.Header().Set("Content-Type", "text/html")
		x := strings.TrimSpace(r.FormValue("word"))
		fmt.Printf("Form word value:\n%s\n", x)
		text := ""

		if _, alphabetized, valid := dictionary.Alphabetizer([]byte(x)); valid {
			if matches, ok := s.FindWords[alphabetized]; ok {
				spacer := ""
				for _, word := range matches {
					text += fmt.Sprintf("%s%s", spacer, word)
					spacer = ", "
				}
			}
		}
		w.Write([]byte(fmt.Sprintf(formHTML, text)))

	}
}

func readRequestData(r *http.Request) ([]solver.Word, error) {
	wordCount, err := strconv.Atoi(strings.TrimSpace(r.FormValue("wordcount")))
	if err != nil {
		return nil, fmt.Errorf("finding value of wordcount: %v\n", err)
	}

	var words []solver.Word

	for wordNumber := 0; wordNumber < wordCount; wordNumber++ {
		var marks []int
		var word []rune
		for charNumber := 0; charNumber < 10; charNumber++ {
			wordCode := fmt.Sprintf("w%dc%d", wordNumber, charNumber)
			wordChar := strings.TrimSpace(r.FormValue(wordCode))
			if wordChar != "" {
				word = append(word, []rune(wordChar)[0])
			}

			markCode := wordCode + "forward"
			m := strings.TrimSpace(r.FormValue(markCode))
			if m == "on" {
				marks = append(marks, charNumber)
			}

		}

		asIsCode := fmt.Sprintf("w%dasis", wordNumber)
		aic := strings.TrimSpace(r.FormValue(asIsCode))
		asIs := false
		if aic == "on" {
			asIs = true
		}

		if len(word) > 0 {
			wrd := solver.Word{
				Word:        word,
				MarkedChars: marks,
				AsIs:        asIs,
			}
			words = append(words, wrd)
		}
	}

	return words, nil
}

func rewriteHTML(words []solver.Word, w http.ResponseWriter) {
	w.Write([]byte(fmt.Sprintf(headerHTML, len(words))))

	for wordNumber, word := range words {
		w.Write([]byte(`	<table border="1">`))

		// Characters in word
		w.Write([]byte("		<tr>\n"))
		for charNumber, char := range word.Word {
			w.Write([]byte(fmt.Sprintf(characterHTML, wordNumber, charNumber, wordNumber, charNumber, char)))
		}
		w.Write([]byte("		</tr>\n"))

		// Marks for characters to carry forward
		w.Write([]byte("		<tr>\n"))
		for charNumber, _ := range word.Word {
			checked := ""

			for _, markedChar := range word.MarkedChars {
				if markedChar == charNumber {
					checked = "checked"
					break
				}
			}

			w.Write([]byte(fmt.Sprintf(markHTML, wordNumber, charNumber, wordNumber, charNumber, checked)))
		}
		w.Write([]byte("		</tr>\n"))

		// Use-as-is flag
		checked := ""
		if word.AsIs {
			checked = "checked"
		}
		w.Write([]byte(fmt.Sprintf(asIsHTML, len(word.Word), wordNumber, checked)))

		w.Write([]byte(`	</table>`))
	}

	w.Write([]byte(footerHTML))
}

var headerHTML string = `<!DOCTYPE html>
<html>
	<head>
	<meta charset="UTF-8">
	<script language="javascript" type="text/javascript">
		function wordcountchange() {
			var Nwords = document.getElementById("wordcount").value;
			var worddiv = document.getElementById("worddiv");
			var generatedHTML = "";
			for (var i = 0; i < Nwords; ++i) {
				generatedHTML += generateWordHTML(i);
			}
			worddiv.innerHTML = wordsDropDown(Nwords);
			worddiv.innerHTML += generatedHTML;
		}
		function generateWordHTML(wordNumber) {
			var generated = '<table border="1">';
			// The characters of the word
			// <td><input type="text" id="w0c0" name="w0c0" size="1" /></td>
			var wordID = 'w'+wordNumber+'c';
			generated += '<tr>';
			for (var i = 0; i < 5; ++i) {
				var charID = wordID + i;
				var prevChar = document.getElementById(charID);
				var prevCharValue = '';
				if (prevChar != null) {
					prevCharValue = prevChar.value;
				}
				generated += '<td><input type="text" size="1" name="'+charID+'" id="'+charID+'" value="' +prevCharValue + '" /></td>';
			}

			// The marked characters of the word, to use to solve the jumble
			generated += '</tr><tr>';
			for (var i = 0; i < 5; ++i) {
				var markID = wordID + i + 'forward';
				generated += '<td><input type="checkbox" name="'
					+markID+'" id="'+markID+'" '
				var prevMark = document.getElementById(markID);
				if (prevMark != null && prevMark.checked) {
						generated += 'checked ';
				}
				generated += '></td>';
			}
			generated += '</tr>';


			// The use-as-is flag, for words the human knows for a fact are part
			// of the solution.
			var asIsID = 'w' + wordNumber + 'asis';
			var prevAsIs = document.getElementById(asIsID);
			generated += '<tr><td colspan="5">Use as-is: <input type="checkbox" name="'+asIsID+'" id="' +asIsID+'" ';
			if (prevAsIs != null && prevAsIs.checked) {
				generated += 'checked ';
			}
			generated += '></td></tr></table>';
			return generated;
		}
		function wordsDropDown(selected) {
				var html = '<table border="0">';
				html += '<tr><td>Number of words:';
				html += '	<select name="wordcount" id="wordcount" onchange="wordcountchange()">';
				for (var i = 1; i <= 5; ++i) {
					html += '<option value="'+i+'"'
					if (i == selected) {
						html += ' selected="true" ';
					}
					html += '>'+i+'</option>';
				}
				html += '</select>';
				html += '</td></tr>';
				return html;
		}
		function setwordcount() {
			document.getElementById("wordcount").value = document.f.words.value;
		}
	</script>
	</head>
	<body onload="setwordcount()" >
	<form name="f" method="post" action="http://localhost:8012/jumble">

	<input type="hidden" name="words" id="words" value="%d" />

	<div id="worddiv" >
		<table border="0">
		<tr><td>Number of words:
			<select name="wordcount" id="wordcount" onchange="wordcountchange()">
				<option value="1">1</option>
				<option value="2">2</option>
				<option value="3">3</option>
				<option value="4">4</option>
				<option value="5">5</option>
			</select>
		</td></tr>

	<table border="1">
`

var footerHTML string = `</table>
	</div>
	<input type="submit" />
	</form>
	</body>
</html>
`
var asIsHTML string = `
		<tr>
				<td colspan="%d">Use as-is: <input type="checkbox" name="w%dasis" %s></td>
		</tr>
`

var characterHTML string = `			<td><input type="text" id="w%dc%d" name="w%dc%d" size="1" value="%c" /></td>
`
var markHTML string = `			<td><input type="checkbox" id="w%dc%dforward" name="w%dc%dforward" %s /></td>
`
