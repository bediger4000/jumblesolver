package srvr

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"jumblesolver/dictionary"
	"jumblesolver/solver"
)

var indexHTML = `
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

var solveHTML = `
<html>
<head>
    <meta charset="UTF-8">
</head>
<body>
<h1>Solve it</h1>
<form name="f" method="post" >
`
var solveHTML2 = `</form>
</body>
</html>
`

var formHTML = `
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

var errorHTML = `<!DOCTYPE html>
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
		fmt.Fprintf(w, indexHTML)
	}
}

func (s *Srvr) handleSolve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.Debug {
			fmt.Printf("Enter handleSolve closure\n")
			defer fmt.Printf("Exit handleSolve closure\n")
		}

		w.Header().Set("Content-Type", "text/html")

		alternates, finalwords, finalwordsizes, err := readSolveData(s.FindWords, r, s.Debug)

		if err != nil {
			fmt.Fprintf(w, errorHTML, err)
			return
		}
		for i, alts := range alternates {
			if len(alts) == 0 {
				fmt.Fprintf(w, errorHTML, fmt.Errorf("no letters at position %d", i))
				return
			}
		}

		if s.Debug {
			fmt.Printf("%d letters in solution\n", len(alternates))
			for _, alts := range alternates {
				fmt.Printf("%v\n", alts)
			}
			fmt.Printf("%d final words\n", finalwords)
		}

		if finalwords == 1 {
			s.singleWordSolution(w, alternates)
			return
		}

		// Multi-word solution
		s.multiWordSolution(w, alternates, finalwords, finalwordsizes)
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
			fmt.Fprintf(w, errorHTML, err)
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
		fmt.Fprintf(w, formHTML, text)

	}
}

// readSolveData returns an array of []rune. The array is len the number
// of "marked" characters in the jumbled words, which is the number of
// characters in the solution. Since a jumbled word can have more than
// a single matching real word, each position in the array is a []rune,
// so that each position can hold alternate characters.
func readSolveData(dict dictionary.Dictionary, r *http.Request, debug bool) ([][]rune, int, []int, error) {
	words, _, err := readRequestData(r, debug)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("reading unjumbled words: %v", err)
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
	if debug {
		fmt.Printf("%d elements in matches array\n", len(matches))
		for i, match := range matches {
			fmt.Printf("\tjumbled word %d has %d matches\n", i, len(match))
		}
	}
	for i, match := range matches {
		if len(match) == 0 {
			return nil, 0, nil, fmt.Errorf("no unjumbled word at position %d", i)
		}
	}

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

	// Find out how many final words there are
	finalWords, finalWordSizes := determineFinalWords(r)
	if finalWords == 1 {
		finalWordSizes = []int{len(jumbledChars)}
	}
	if debug {
		fmt.Printf("%d words in solution\n", finalWords)
		fmt.Printf("solution word sizes: %v\n", finalWordSizes)
	}

	return jumbledChars, finalWords, finalWordSizes, nil
}

func determineFinalWords(r *http.Request) (int, []int) {
	wrdCntStr := r.FormValue("soluWrdCnt")
	solutionCount := 1
	if wrdCntStr != "" {
		n, err := strconv.Atoi(wrdCntStr)
		if err == nil {
			solutionCount = n
		}
	}
	var solutionSizes []int
	if solutionCount > 1 {
		soluWrdSz := r.FormValue("soluWrdSz")
		if soluWrdSz != "" {
			fields := strings.Fields(soluWrdSz)
			for _, szStr := range fields {
				n, err := strconv.Atoi(szStr)
				if err == nil {
					solutionSizes = append(solutionSizes, n)
				}
			}
		}
	}
	return solutionCount, solutionSizes
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
		return nil, false, fmt.Errorf("finding value of wordcount: %v", err)
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
				if debug {
					fmt.Printf("\tletter %d, mark name %q: %q\n", charNumber, markCode, r.FormValue(markCode))
				}
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

	fmt.Fprintf(w, headerHTML, len(words))

	for wordNumber, word := range words {
		fmt.Fprintf(w, "	<table border='1'>\n")

		// Characters in word
		fmt.Fprintf(w, "		<tr id='w%drow'>\n", wordNumber)
		for charNumber, char := range word.Word {
			fmt.Fprintf(w, characterHTML, wordNumber, charNumber, wordNumber, charNumber, char, wordNumber, charNumber)
		}
		fmt.Fprintf(w, "		</tr>\n")

		// Marks for characters to carry forward
		fmt.Fprintf(w, "		<tr id='w%dmarks'>\n", wordNumber)
		for charNumber := range word.Word {
			checked := ""

			for _, markedChar := range word.MarkedChars {
				if markedChar == charNumber {
					checked = "checked"
					break
				}
			}

			fmt.Fprintf(w, markHTML, wordNumber, charNumber, wordNumber, charNumber, checked)
		}
		fmt.Fprintf(w, "		</tr>\n")

		// Use-as-is flag
		checked := ""
		if word.AsIs {
			checked = "checked"
		}
		fmt.Fprintf(w, asIsHTML, len(word.Word)-1, wordNumber, wordNumber, checked)
		fmt.Fprintf(w, addLetterHTML, wordNumber, wordNumber)

		// words that might match
		fmt.Fprintf(w, "\t\t<tr>\n\t\t\t<td colspan=%d>%s</td>\n\t\t</tr>\n", len(word.Word), strings.Join(matches[wordNumber], ", "))

		fmt.Fprintf(w, "	</table>\n")
	}

	fmt.Fprintf(w, footerHTML)
}

func (s *Srvr) multiWordSolution(w http.ResponseWriter, alternates [][]rune, finalcount int, finalsizes []int) {
	if s.Debug {
		fmt.Printf("Enter multiWordSolution, %d solution words, %v\nAlternates:\n", finalcount, finalsizes)
		defer fmt.Printf("Exit multiWordSolution\n")
	}

	fmt.Fprintf(w, "<h4>Possible Alternates</h4>\n<table border='1'>\n")
	for _, alternate := range alternates {
		fmt.Fprintf(w, "\t<tr><td>%s</td></tr>\n", string(alternate))
	}
	fmt.Fprintf(w, "</table>\n")

	fmt.Fprintf(w, solveHTML)
	fmt.Fprintf(w, "<h3>Possible multi-word solutions</h3>\n")
	keyCombos := solver.GenerateKeyCombos(s.Debug, alternates, finalcount, finalsizes)
	if s.Debug {
		fmt.Printf("multi-word solutions %d keyCombos\n", len(keyCombos))
	}
	multiWordSolutions := solver.SolutionsFromKeyCombos(s.Debug, keyCombos, s.FindWords)
	fmt.Fprintf(w, "<h4>%d solutions</h4>\n", len(multiWordSolutions))
	sort.Sort(StringSliceSlice(multiWordSolutions))
	for _, solution := range multiWordSolutions {
		fmt.Fprintf(w, "<p>%s</p>\n", strings.Join(solution, " "))
	}
	fmt.Fprintf(w, solveHTML2)
}

func (s *Srvr) singleWordSolution(w http.ResponseWriter, alternates [][]rune) {

	fmt.Fprintf(w, solveHTML)
	fmt.Fprintf(w, "<h4>Possible Alternates</h4>\n<table border='1'>\n")
	for _, alternate := range alternates {
		fmt.Fprintf(w, "\t<tr><td>%s</td></tr>\n", string(alternate))
	}
	fmt.Fprintf(w, "</table>\n")

	uniquematches := solver.FindUniqueMatches(alternates, s.FindWords)

	fmt.Fprintf(w, "<h2>Found %d unique keys</h2>\n", len(uniquematches))

	fmt.Fprintf(w, "<h2>Possible Unique Single Word Solutions</h2>\n")
	for key, matches := range uniquematches {
		if len(matches) > 0 {
			fmt.Fprintf(w, "<h4>Key %q</h4>\n<pre>", key)
			for _, match := range matches {
				fmt.Fprintf(w, "%s\n", match)
			}
			fmt.Fprintf(w, "</pre>\n")
		}
	}

	fmt.Fprintf(w, solveHTML2)
}

func noWordsHTML(w http.ResponseWriter) {
	fmt.Fprintf(w, headerHTML, 1)
	fmt.Fprintf(w, "	<table border='1'>\n")
	fmt.Fprintf(w, "		<tr id='w0row'>\n")
	for charNumber := 0; charNumber < 6; charNumber++ {
		fmt.Fprintf(w, emptyCharacterHTML, charNumber, charNumber, 0, charNumber)
	}
	fmt.Fprintf(w, "		</tr>\n")
	fmt.Fprintf(w, "		<tr id='w0marks'>\n")
	for charNumber := 0; charNumber < 6; charNumber++ {
		fmt.Fprintf(w, emptyMarkHTML, charNumber, charNumber)
	}
	fmt.Fprintf(w, asIsHTML, 4, 0, 0, "")
	fmt.Fprintf(w, addLetterHTML, 0, 0)
	fmt.Fprintf(w, "		</tr>\n	</table>\n")
	fmt.Fprintf(w, footerHTML)
}

var headerHTML = `<!DOCTYPE html>
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
				charCount = 6;
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
				generated += '<td><input type="text" size="1" name="'
					+charID+'" id="'+charID+'" value="' +prevCharValue
					+ '" oninput="letterchanged('+wordNumber+','+i+')" /></td>';
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
				for (var i = 1; i <= 7; ++i) {
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

		function letterchanged(wordno, charno) {
			var wordcode = 'w'+wordno+'c'+charno;
			var cell = document.getElementById(wordcode);
			var val = cell.value;
			// if val is one character, move to the next cell
			if (val.length == 1) {
				var nextwordcode = 'w'+wordno+'c'+(charno+1);
				var nextcell = document.getElementById(nextwordcode);
				if (nextcell != null) {
					nextcell.focus();
				}
			}
			// if val is more than one character, trim it,
			// move to next cell, put next character there
			if (val.length > 1) {
				cell.value = val[0];
				var nextwordcode = 'w'+wordno+'c'+(charno+1);
				var nextcell = document.getElementById(nextwordcode);
				if (nextcell != null) {
					nextcell.value = val[1];
					nextcell.focus();
				}
			}
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
				<option value="6">6</option>
				<option value="7">7</option>
			</select>
		</td></tr>
		</table>
`

var footerHTML = `
	</div>
	<input type="button" value="Unjumble Words" onclick="submitjumble()" />
	<br />
	<br />
	<p>Solution words: <input type="text" id="soluWrdCnt" name="soluWrdCnt" size="1" /></p>
	<p>Solution sizes: <input type="text" id="soluWrdSz"  name="soluWrdSz" size="8" /></p>
	<br />
	<input type="submit" value="Solve" onclick="submitsolve()" />
	<br />
	<br />
	<br />
	<input type="submit" value="Reset" onclick="submitreset()" />
	</form>
	</body>
</html>
`
var asIsHTML = `
		<tr>
			<td colspan="%d">Use as-is: <input type="checkbox" name="w%dasis" id="w%dasis" %s></td>
`
var addLetterHTML = `
			<td></td><td><input type="button" name="w%db" value="Add letter" onclick="addletter(%d)" /></td>
		</tr>
`

var characterHTML = `			<td><input type="text" id="w%dc%d" name="w%dc%d" size="1" value="%c" oninput="letterchanged(%d, %d)" /></td>
`
var emptyCharacterHTML = `			<td><input type="text" id="w0c%d" name="w0c%d" size="1"  oninput="letterchanged(%d, %d)" /></td>
`
var markHTML = `			<td><input type="checkbox" id="w%dc%dforward" name="w%dc%dforward" %s /></td>
`
var emptyMarkHTML = `			<td><input type="checkbox" id="w0c%dforward" name="w0c%dforward" /></td>
`

type StringSliceSlice [][]string

func (sss StringSliceSlice) Len() int { return len(sss) }
func (sss StringSliceSlice) Less(i, j int) bool {
	for idx := 0; idx < len(sss[i]); idx++ {
		if sss[i][idx] != sss[j][idx] {
			return sss[i][idx] < sss[j][idx]
		}
	}
	return false
}
func (sss StringSliceSlice) Swap(i, j int) {
	sss[i], sss[j] = sss[j], sss[i]
}
