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

var solveHTML string = `
<html>
<head>
    <meta charset="UTF-8">
</head>
<body>
<h1>Solve it</h1>
<form name="f" method="post" >
`
var solveHTML2 string = `</form>
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

func (s *Srvr) handleSolve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.Debug {
			fmt.Printf("Enter handleSolve closure\n")
			defer fmt.Printf("Exit handleSolve closure\n")
		}

		w.Header().Set("Content-Type", "text/html")

		alternates, err := readSolveData(s.FindWords, r, s.Debug)

		if err != nil {
			w.Write([]byte(fmt.Sprintf(errorHTML, err)))
			return
		}

		keys := solver.CreateKeys(alternates, len(alternates))
		fmt.Printf("Found %d unique keys\n", len(keys))
		if len(keys) < 30 {
			for _, key := range keys {
				fmt.Printf("%q\n", key)
				if matches, ok := s.FindWords[key]; ok {
					w.Write([]byte("<h2>Possible Single Word Solutions</h2>\n"))
					for _, match := range matches {
						w.Write([]byte(fmt.Sprintf("<p><keyboard>%s</keyboard></p>\n", match)))
					}
				}
			}
		}

		if s.Debug {
			for i, alternate := range alternates {
				fmt.Printf("Solution letter %d: %s\n", i, string(alternate))
			}
		}

		w.Write([]byte(solveHTML))
		for _, alternate := range alternates {
			w.Write([]byte(fmt.Sprintf("<p>%s</p>\n", string(alternate))))
		}
		w.Write([]byte(solveHTML2))

	}
}

func (s *Srvr) handleJumble() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.Debug {
			fmt.Printf("Enter handleJumble closure\n")
			defer fmt.Printf("Exit handleJumble closure\n")
		}

		w.Header().Set("Content-Type", "text/html")

		words, reset, err := readRequestData(r, s.Debug)

		if err != nil {
			w.Write([]byte(fmt.Sprintf(errorHTML, err)))
			return
		}

		if reset {
			rewriteHTML(nil, nil, w)
			return
		}

		guesses := solver.LookupWords(s.FindWords, words)

		rewriteHTML(words, guesses, w)
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

		if _, alphabetized, valid := dictionary.Alphabetizer([]rune(x)); valid {
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

// readSolveData does some stuff
func readSolveData(dict dictionary.Dictionary, r *http.Request, debug bool) ([][]rune, error) {
	words, _, err := readRequestData(r, debug)
	if err != nil {
		return nil, fmt.Errorf("reading unjumbled words: %v\n", err)
	}
	if debug {
		fmt.Printf("%d jumbled words\n", len(words))
		for i, w := range words {
			fmt.Printf("\tword %d: %q, use as-is: %v\n", i, string(w.Word), w.AsIs)
		}
	}

	// Find out how many marked characters exist in the jumbled words
	markedCount := 0
	for _, word := range words {
		markedCount += len(word.MarkedChars)
	}
	if debug {
		fmt.Printf("%d marked solution letters total\n", markedCount)
	}

	// The solution has markedCount number of letters in it.
	// Each of those markedCount may have more than 1 alternative.

	matches := solver.LookupWords(dict, words)

	jumbledChars := make([][]rune, markedCount)
	jumbledCharNum := 0

	for wordNum, word := range words {
		for _, markNum := range word.MarkedChars {
			if word.AsIs {
				jumbledChars[jumbledCharNum] = []rune{word.Word[markNum]}
				jumbledCharNum++
				continue
			}
			matching := matches[wordNum]
			var matchChars []rune
			for _, match := range matching {
				rm := []rune(match)
				foundit := false
				for _, alreadyPresent := range matchChars {
					if alreadyPresent == rm[markNum] {
						foundit = true
						break
					}
				}
				if !foundit {
					matchChars = append(matchChars, rm[markNum])
				}
			}
			jumbledChars[jumbledCharNum] = matchChars
			jumbledCharNum++
		}
	}

	return jumbledChars, nil
}

func readRequestData(r *http.Request, debug bool) ([]solver.Word, bool, error) {
	if r.FormValue("wordcount") == "" {
		if debug {
			fmt.Printf("form input wordcount has zero-length string value\n")
		}
		return []solver.Word{}, true, nil
	}

	wordCount, err := strconv.Atoi(strings.TrimSpace(r.FormValue("wordcount")))
	if err != nil {
		return nil, false, fmt.Errorf("finding value of wordcount: %v\n", err)
	}

	if debug {
		fmt.Printf("wordcount %d\n", wordCount)
	}

	var words []solver.Word

	for wordNumber := 0; wordNumber < wordCount; wordNumber++ {
		if debug {
			fmt.Printf("Jumbled word %d:\n", wordNumber)
		}
		var marks []int
		var word []rune
		for charNumber := 0; true; charNumber++ {
			wordCode := fmt.Sprintf("w%dc%d", wordNumber, charNumber)
			wordChar := strings.TrimSpace(r.FormValue(wordCode))
			if wordChar != "" {
				if debug {
					fmt.Printf("\tletter %d, field name %q: '%c'\n", charNumber, wordCode, []rune(r.FormValue(wordCode))[0])
				}
				word = append(word, []rune(wordChar)[0])

				markCode := wordCode + "forward"
				m := strings.TrimSpace(r.FormValue(markCode))
				fmt.Printf("\tletter %d, mark name %q: %q\n", charNumber, markCode, r.FormValue(markCode))
				if m == "on" {
					marks = append(marks, charNumber)
				}
				continue
			}
			break
		}

		asIsCode := fmt.Sprintf("w%dasis", wordNumber)
		aic := strings.TrimSpace(r.FormValue(asIsCode))
		asIs := false
		if aic == "on" {
			asIs = true
		}
		if debug {
			fmt.Printf("Use-as-is code %q, value %v\n", asIsCode, asIs)
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

	return words, false, nil
}

func rewriteHTML(words []solver.Word, matches [][]string, w http.ResponseWriter) {

	// Called without any POST data
	if len(words) == 0 {
		noWordsHTML(w)
		return
	}

	w.Write([]byte(fmt.Sprintf(headerHTML, len(words))))

	for wordNumber, word := range words {
		w.Write([]byte("	<table border='1'>\n"))

		// Characters in word
		w.Write([]byte(fmt.Sprintf("		<tr id='w%drow'>\n", wordNumber)))
		for charNumber, char := range word.Word {
			w.Write([]byte(fmt.Sprintf(characterHTML, wordNumber, charNumber, wordNumber, charNumber, char)))
		}
		w.Write([]byte("		</tr>\n"))

		// Marks for characters to carry forward
		w.Write([]byte(fmt.Sprintf("		<tr id='w%dmarks'>\n", wordNumber)))
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
		w.Write([]byte(fmt.Sprintf(asIsHTML, len(word.Word)-1, wordNumber, wordNumber, checked)))
		w.Write([]byte(fmt.Sprintf(addLetterHTML, wordNumber, wordNumber)))

		// words that might match
		w.Write([]byte(fmt.Sprintf("\t\t<tr>\n\t\t\t<td colspan=%d>%s</td>\n\t\t</tr>\n", len(word.Word), strings.Join(matches[wordNumber], ", "))))

		w.Write([]byte("	</table>\n"))
	}

	w.Write([]byte(footerHTML))
}

func noWordsHTML(w http.ResponseWriter) {
	w.Write([]byte(fmt.Sprintf(headerHTML, 1)))
	w.Write([]byte("	<table border='1'>\n"))
	w.Write([]byte("		<tr id='w0row'>\n"))
	for charNumber := 0; charNumber < 5; charNumber++ {
		w.Write([]byte(fmt.Sprintf(emptyCharacterHTML, charNumber, charNumber)))
	}
	w.Write([]byte("		</tr>\n"))
	w.Write([]byte("		<tr id='w0marks'>\n"))
	for charNumber := 0; charNumber < 5; charNumber++ {
		w.Write([]byte(fmt.Sprintf(emptyMarkHTML, charNumber, charNumber)))
	}
	w.Write([]byte(fmt.Sprintf(asIsHTML, 4, 0, 0, "")))
	w.Write([]byte(fmt.Sprintf(addLetterHTML, 0, 0)))
	w.Write([]byte("		</tr>\n	</table>\n"))
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
			generated += '<tr id="w'+wordNumber+'row">';

			// Count characters in this word
			var charCount = 0;
			for (var i = 0; true; ++i) {
				var prevChar = document.getElementById(wordID+i);
				if (prevChar == null) {
					break;
				}	
				++charCount;
			}
			if (charCount == 0) {
				charCount = 5;
			}

			// redo the existing characters, or fill in charCount
			// blank spaces for user to input new characters.
			for (var i = 0; i < charCount; ++i) {
				var charID = wordID + i;
				var prevChar = document.getElementById(charID);
				var prevCharValue = '';
				if (prevChar != null) {
					prevCharValue = prevChar.value;
				}
				generated += '<td><input type="text" size="1" name="'+charID+'" id="'+charID+'" value="' +prevCharValue + '" /></td>';
			}

			// The marked characters of the word, to use to solve the jumble
			generated += '</tr><tr id="w'+wordNumber+'marks">';
			for (var i = 0; i < charCount; ++i) {
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
			generated += '<tr><td colspan="'+charCount+'">Use as-is: <input type="checkbox" name="'+asIsID+'" id="' +asIsID+'" ';


			var prevAsIs = document.getElementById(asIsID);
			if (prevAsIs != null && prevAsIs.checked) {
				generated += 'checked ';
			}
			generated += '/></td>';

			// The add-a-letter button
			generated += '<td><input type="button" name="w'+wordNumber+'b" value="Add letter" onclick="addletter('+wordNumber+')" /></td></tr></table>';

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
		function addletter(wordno) {
			var row = document.getElementById("w"+wordno+"row");
			var name = "w"+wordno+"c"+row.childElementCount;
			var newdatum = row.insertCell(-1);
			newdatum.innerHTML = '<input type="text" name="'+name+'" id="'+name+'" size="1" />';

			// The corresponding "carry this character forward" mark.
			var markrow = document.getElementById("w"+wordno+"marks");
			var markdatum = markrow.insertCell(-1);
			markdatum.innerHTML = '<input type="checkbox" id="'+name+'forward" name="'+name+'forward" />';
		}

		function setwordcount() {
			document.getElementById("wordcount").value = document.f.words.value;
		}
		function submitjumble() {
			document.f.action = "/jumble";
			document.f.submit();
		}
		function submitsolve() {
			document.f.action = "/solve";
			document.f.submit();
		}
		function submitreset() {
			document.f.action = "/jumble";
			document.getElementById("wordcount").value = "";
			document.f.submit();
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
		</table>
`

var footerHTML string = `
	</div>
	<input type="button" value="Unjumble Words" onclick="submitjumble()" />
	<br />
	<input type="submit" value="Solve" onclick="submitsolve()" />
	<br />
	<br />
	<input type="submit" value="Reset" onclick="submitreset()" />
	</form>
	</body>
</html>
`
var asIsHTML string = `
		<tr>
			<td colspan="%d">Use as-is: <input type="checkbox" name="w%dasis" id="w%dasis" %s></td>
`
var addLetterHTML string = `
			<td><input type="button" name="w%db" value="Add letter" onclick="addletter(%d)" /></td>
		</tr>
`

var characterHTML string = `			<td><input type="text" id="w%dc%d" name="w%dc%d" size="1" value="%c" /></td>
`
var emptyCharacterHTML string = `			<td><input type="text" id="w0c%d" name="w0c%d" size="1" /></td>
`
var markHTML string = `			<td><input type="checkbox" id="w%dc%dforward" name="w%dc%dforward" %s /></td>
`
var emptyMarkHTML string = `			<td><input type="checkbox" id="w0c%dforward" name="w0c%dforward" /></td>
`
