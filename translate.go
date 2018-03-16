package htmltotext

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func tableBoarder(wordSize, numOfRows, wordSpacing int) string {
	wall := strings.Repeat("-", wordSize+wordSpacing)
	corner := "+"

	result := fmt.Sprintf("%s%s", corner, wall)
	result = strings.Repeat(result, numOfRows)
	result += corner
	return result
}

func tableContent(item string, totalSpaces, leftSpacing int) string {
	rightSpacing := totalSpaces - leftSpacing - len(item)

	result := fmt.Sprintf("%s%s%s", strings.Repeat(" ", leftSpacing), item, strings.Repeat(" ", rightSpacing))
	return result
}

func tagInList(tokenBytes []byte, tags []Tag) bool {
	result := false
	for _, tag := range tags {
		if bytes.HasPrefix(tokenBytes, tag.Byte()) {
			result = true
			break
		}
	}
	return result
}

func processToken(token html.Token, stack Stack, tempt, result string, links []string, listCount int, tableRows int, tableColumns int, tableItems []string, ignoreToken bool) (Stack, string, string, []string, int, int, int, []string, bool, error) {
	var err error
	tagsThatPopTheStack := []Tag{CloseDivTag, CloseH1Tag, CloseH2Tag, CloseH3Tag, CloseH4Tag, CloseH5Tag, CloseH6Tag, ClosePTag, CloseOLTag, CloseULTag, CloseATag, CloseTableTag, CloseLITag}
	tokenString := strings.TrimSpace(token.String())
	tokenBytes := []byte(tokenString)

	//Used to ignore empty strings
	if len(tokenString) > 0 {

		//Check for Tag
		if bytes.HasPrefix(tokenBytes, []byte("<")) && bytes.HasSuffix(tokenBytes, []byte(">")) {
			if tagInList(tokenBytes, tagsThatPopTheStack) {
				stack, _, err = stack.Pop()
				if err != nil {
					return stack, tempt, result, links, listCount, tableRows, tableColumns, tableItems, ignoreToken, err
				}
			}

			switch {
			case bytes.HasPrefix(tokenBytes, OpenDivTag.Byte()):
				stack = stack.Push(OpenDivTag)

			case bytes.HasPrefix(tokenBytes, CloseDivTag.Byte()):
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s\n", result, tempt)
					tempt = ""
				} else {
					tempt = fmt.Sprintf("%s\n", tempt)
				}

			case bytes.HasPrefix(tokenBytes, OpenH1Tag.Byte()):
				stack = stack.Push(OpenH1Tag)

			case bytes.HasPrefix(tokenBytes, OpenH2Tag.Byte()):
				stack = stack.Push(OpenH2Tag)

			case bytes.HasPrefix(tokenBytes, OpenH3Tag.Byte()):
				stack = stack.Push(OpenH3Tag)

			case bytes.HasPrefix(tokenBytes, OpenH4Tag.Byte()):
				stack = stack.Push(OpenH4Tag)

			case bytes.HasPrefix(tokenBytes, OpenH5Tag.Byte()):
				stack = stack.Push(OpenH5Tag)

			case bytes.HasPrefix(tokenBytes, OpenH6Tag.Byte()):
				stack = stack.Push(OpenH6Tag)

			case tagInList(tokenBytes, []Tag{CloseH1Tag, CloseH2Tag, CloseH3Tag, CloseH4Tag, CloseH5Tag, CloseH6Tag}):
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s\n", result, tempt)
					tempt = ""
				} else {
					tempt = fmt.Sprintf("%s\n\n", tempt)
				}

			case bytes.HasPrefix(tokenBytes, OpenPTag.Byte()):
				stack = stack.Push(OpenPTag)

			case bytes.HasPrefix(tokenBytes, ClosePTag.Byte()):
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s\n\n", result, tempt)
					tempt = ""
				} else {
					result = fmt.Sprintf("%s%s\n", result, tempt)
					tempt = ""
				}

			case bytes.HasPrefix(tokenBytes, OpenOLTag.Byte()):
				if stack.Contains(OpenPTag) {
					tempt = fmt.Sprintf("%s\n", tempt)
				}

				for _, attr := range token.Attr {
					if strings.EqualFold(attr.Key, "start") {
						listCount, err = strconv.Atoi(attr.Val)
						if err != nil {
							listCount = 1
						}
					}
				}
				stack = stack.Push(OpenOLTag)

			case bytes.HasPrefix(tokenBytes, CloseOLTag.Byte()):
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s\n", result, tempt)
					tempt = ""
				}

				//Reset variable
				listCount = 1

			case bytes.HasPrefix(tokenBytes, OpenULTag.Byte()):
				if stack.Contains(OpenPTag) {
					tempt = fmt.Sprintf("%s\n", tempt)
				}

				stack = stack.Push(OpenULTag)

			case bytes.HasPrefix(tokenBytes, CloseULTag.Byte()):
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s\n", result, tempt)
					tempt = ""
				}

				//Reset variable
				listCount = 1

			case bytes.HasPrefix(tokenBytes, OpenLITag.Byte()):
				stack = stack.Push(OpenLITag)

				if stack.Contains(OpenOLTag) {
					tempt = fmt.Sprintf("%s  %d. ", tempt, listCount)
					listCount++
				} else if stack.Contains(OpenULTag) {
					tempt = fmt.Sprintf("%s  * ", tempt)
				}

			case bytes.HasPrefix(tokenBytes, CloseLITag.Byte()):
				tempt = fmt.Sprintf("%s\n", tempt)

			case bytes.HasPrefix(tokenBytes, OpenATag.Byte()):
				stack = stack.Push(OpenATag)

				//Capture the href
				for _, attr := range token.Attr {
					if strings.EqualFold(attr.Key, "href") {
						links = append(links, attr.Val)
					}
				}

			case bytes.HasPrefix(tokenBytes, CloseATag.Byte()):
				var addedLink bool // Used to denote it the link was added to the text

				if stack.Contains(OpenTableTag) {
					//Case: a tag nested in a table tag
					tableItems[len(tableItems)-1] += fmt.Sprintf("[%d]", len(links))
					addedLink = true

					//Case: Table with no rows -- A.K.A not a real table...
					tempt = fmt.Sprintf("%s[%d]", tempt, len(links))
				}

				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s[%d]", result, tempt, len(links))
					tempt = ""
					addedLink = true
				}

				//Makes sure that the link was added
				if !addedLink {
					tempt = fmt.Sprintf("%s[%d]", tempt, len(links))
				}

			case bytes.HasPrefix(tokenBytes, OpenTableTag.Byte()):
				stack = stack.Push(OpenTableTag)

			case bytes.HasPrefix(tokenBytes, CloseTableTag.Byte()):
				//Need to create the table
				//Step 1: find the longest item in the table
				var wordLength int
				for _, item := range tableItems {
					if wordLength < len(item) {
						wordLength = len(item)
					}
				}

				//Step 2: Create the Table
				var count int
				var tableString string
				wordSpacing := 2 //  The amount of space padding I want for largest item -- left/right padding
				leftSpacing := wordSpacing / 2
				maxSpaces := wordLength + wordSpacing //  Total number of spaces in each box

				//If this is an actual table with rows then do this
				if tableRows != 0 {

					for y := 0; y < tableColumns+1; y++ {
						tableString = fmt.Sprintf("%s%s\n", tableString, tableBoarder(wordLength, tableRows, wordSpacing))

						for x := 0; x < tableRows; x++ {

							//Need to end the loop
							if y == tableColumns {
								break
							}

							if len(tableItems) > count {
								tableString = fmt.Sprintf("%s|%s", tableString, tableContent(tableItems[count], maxSpaces, leftSpacing))
							} else {
								tableString = fmt.Sprintf("%s|%s", tableString, tableContent("", maxSpaces, leftSpacing))
							}
							count++
						}
						//I Need this to stop one iterval prior to the end
						if y < tableColumns {
							tableString = fmt.Sprintf("%s|\n", tableString)
						}
					}
				}

				//Step 3: add table to result
				if tableRows != 0 {
					result += fmt.Sprintf("%s\n", tableString)
				} else {
					result += fmt.Sprintf("%s\n", tempt)
				}
				//Step 4: reset variables
				tableRows = 0
				tableColumns = 0
				tableItems = []string{}
				tempt = ""

			case bytes.HasPrefix(tokenBytes, OpenTRTag.Byte()):
				tableColumns++

			case bytes.HasPrefix(tokenBytes, OpenTHeadTag.Byte()):
				//Needed so that the <thead tag does not get
				//mistaken for a <th tag

			case bytes.HasPrefix(tokenBytes, OpenTHTag.Byte()):
				tableRows++

			case bytes.HasPrefix(tokenBytes, OpenHeadTag.Byte()):
				//Turn on ignore token
				ignoreToken = true

			case bytes.HasPrefix(tokenBytes, CloseHeadTag.Byte()):
				//Turn off ignore token
				ignoreToken = false

			case bytes.HasPrefix(tokenBytes, BreakTag.Byte()):
				if stack.Contains(OpenPTag, OpenH1Tag, OpenH2Tag, OpenH3Tag, OpenH4Tag, OpenH5Tag, OpenH6Tag) {
					tempt = fmt.Sprintf("%s\n", tempt)
				} else {
					result += fmt.Sprintf("%s\n\n", tempt)
				}

			default:
			}
		} else {
			if stack.Contains(OpenTableTag) {
				tableItems = append(tableItems, tokenString)
				tempt += tokenString
			} else if stack.Contains(OpenPTag, OpenDivTag, OpenLITag) {
				//Keep original spacing
				tempt += token.String()

			} else if ignoreToken {
				//Do nothing
			} else {
				tempt += tokenString
			}
		}
	}

	//Debug statements
	//fmt.Printf("ProcessToken -- token: %s\n\n", token)
	//fmt.Printf("ProcessToken -- tempt: %s\n\n", tempt)
	//fmt.Printf("Stack: ")
	//fmt.Print(stack)
	//fmt.Printf("Tempt: %s--[end]\n", tempt)
	//fmt.Printf("Result: %s\n", result)
	//fmt.Printf("TableRows: %d, TableColumns %d\n", tableRows, tableColumns)
	//fmt.Print("TableItems: ")
	//fmt.Println(tableItems)
	//fmt.Printf("%s\n\n", strings.Repeat("-", 20))
	return stack, tempt, result, links, listCount, tableRows, tableColumns, tableItems, ignoreToken, nil
}

//Translate -- translates html to text
func Translate(reader io.Reader) (string, []string, error) {
	var result string
	var tempt string
	var err error
	var links []string
	var tableRows int
	var tableColumns int
	var tableItems []string
	var ignoreToken bool // Used to ignore data in tags that are not rendered. For example, the style tag.

	listCount := 1 // Used for ol tag to count li
	stack := NewStack()

	data := html.NewTokenizer(reader)

	for {
		tokenType := data.Next()

		if tokenType == html.ErrorToken {
			if data.Err() == io.EOF {
				break
			}

			//Error case
			return result, links, fmt.Errorf("html Parser err token: %s", data.Err())
		}
		// Process the current token.
		token := data.Token()

		stack, tempt, result, links, listCount, tableRows, tableColumns, tableItems, ignoreToken, err = processToken(token, stack, tempt, result, links, listCount, tableRows, tableColumns, tableItems, ignoreToken)
		if err != nil {
			return result, links, fmt.Errorf("ProcessToken had an error: %s", err.Error())
		}
	}

	//Convert html string to readable text
	result = html.UnescapeString(result)

	return result, links, nil
}
