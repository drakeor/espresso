package template

import "io/ioutil"
import "regexp"
import "strings"

func Using(path string) string {
	p := WebRoot + "/" + strings.Replace(path, ".", "/", -1) + ".lua"

	return ParseTemplate(p)
}

/**
 * Parse a HTML/Lua Template.
 * file: file name to parse
 * return: lua code, parsed
 */
func ParseTemplate(file string) string {
	text, _ := ioutil.ReadFile(file)

	// Match Lua using macro.
	rmac := regexp.MustCompile(`<\?lua([^\0]*?)\?>`)
	using := regexp.MustCompile(`using\(["']([^\)]*)["']\);`)
	codeBlocks := rmac.FindAllString(string(text), -1)

	using_fixed := string(text)
	for _, block := range codeBlocks {
		usingUnits := using.FindAllString(block, -1)
		usingUnitsIndex := using.FindAllStringIndex(block, -1)
		numDiff := 0
		for j, unit := range usingUnits {
			textBefore := using_fixed[0 : usingUnitsIndex[j][0]+numDiff]
			textAfter := using_fixed[usingUnitsIndex[j][1]+numDiff : len(using_fixed)]
			//Unit = using("...");
			//len(using(") = 7. len(");) = 3.
			importedUnit := Using(unit[7 : len(unit)-3])

			numDiff = numDiff + (len(importedUnit) - len(unit))

			using_fixed = textBefore + importedUnit + textAfter
		}
	}

	text = []byte(using_fixed)

	// Match blocks of html buffer code
	reg := regexp.MustCompile(`\?>([^\0]*?)<\?lua`)
	r2 := regexp.MustCompile(`\[\[`)
	r2b := regexp.MustCompile(`\]\]`)
	bufferBlocks := reg.FindAllString(string(text), -1)
	bufferBlocksIndex := reg.FindAllIndex(text, -1)

	// If there aren't any lua tags, just print the entire thing.
	dCheck := regexp.MustCompile(`<\?lua`).FindAllIndex(text, -1)
	if len(dCheck) == 0 {
		fixed_block := r2.ReplaceAllString(string(text), `[ [`)
		fixed_block = r2b.ReplaceAllString(fixed_block, `] ]`)
		fixed_final := "print [[" + fixed_block + "]];"
		return fixed_final
	}

	fixed := string(text)
	numDiff := 0
	// Loop through the blocks and parse out the "'s
	for i, block := range bufferBlocks {
		// Split the whole file into three segments, the characters before
		// the characters after
		// and the block itself.
		textBefore := fixed[0 : bufferBlocksIndex[i][0]+numDiff]
		textAfter := fixed[bufferBlocksIndex[i][1]+numDiff : len(fixed)]
		fixedBlock := r2.ReplaceAllString(block, `[ [`)
		fixedBlock = r2b.ReplaceAllString(fixedBlock, `] ]`)

		// We're adding characters, this invalidates the length of the index.
		// Adding this fixes that.
		numDiff = numDiff + (len(fixedBlock) - len(block))

		// Combine and move on.
		fixed = string(textBefore) + fixedBlock + string(textAfter)
	}

	// Convert the buffered blocks into Lua code.
	// TODO: Implement buffer data structure.
	fixed_temp1 := reg.ReplaceAllString(fixed, `print [[$1]];`)

	// Start parsing the beginning and end of the file.
	var fixed_final string = fixed_temp1

	// If the file starts with a <?lua, there is no buffered code before the
	// open lua brace.
	if fixed_temp1[0:5] == "<?lua" {
		// Erase the <?lua brace.
		fixed_final = regexp.MustCompile(`<\?lua`).ReplaceAllString(fixed_temp1, ``)
	} else {
		// There is buffered code before the first brace.
		// Match the buffered code.
		r3 := regexp.MustCompile(`([^\0]*?)<\?lua`)
		fixedIndexes := r3.FindAllStringIndex(fixed_temp1, -1)

		// String of the buffered code
		block := fixed_final[0 : fixedIndexes[0][1]-5]

		// Parse out the quotes and convert to Lua code.
		fixed_block := r2.ReplaceAllString(block, `[ [`)
		fixed_block = r2b.ReplaceAllString(fixed_block, `] ]`)
		fixed_block = "print [[" + fixed_block + "]];"

		// Combine.
		fixed_final = fixed_block + fixed_final[fixedIndexes[0][1]:len(fixed_final)]
	}

	// If the file ends with a ?>, there is no buffered code at the end.
	if fixed_final[len(fixed_final)-2:len(fixed_final)] == "?>" {
		// Erase the ?> brace.
		fixed_final = regexp.MustCompile(`\?>`).ReplaceAllString(fixed_final, ``)
	} else {
		// There is buffered code at the end of the file.
		// Match that buffered code.
		r3 := regexp.MustCompile(`\?>([^\0]*?)`)
		fixedIndexes := r3.FindAllStringIndex(fixed_final, -1)

		// String of the buffered code
		block := fixed_final[fixedIndexes[0][1]:len(fixed_final)]

		// Parse out the quote and convert to Lua code.
		fixed_block := r2.ReplaceAllString(block, `[ [`)
		fixed_block = r2b.ReplaceAllString(fixed_block, `] ]`)
		fixed_block = "print [[" + fixed_block + "]];"

		// Combine.
		fixed_final = fixed_final[0:fixedIndexes[0][1]-2] + fixed_block
	}

	// Return the Lua code.
	return fixed_final
}
