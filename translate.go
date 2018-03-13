package htmlToText

import (
	"bytes"
	"fmt"
	"io"
	"log"
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

func processToken(token html.Token, stack Stack, tempt, result string, links []string, parentTag Tag, listCount int, tableRows int, tableColumns int, tableItems []string, ignoreToken bool) (Stack, string, string, []string, Tag, int, int, int, []string, bool, error) {
	log.Printf("ProcessToken -- token: %s\n\n", token)
	log.Printf("ProcessToken -- tempt: %s\n\n", tempt)
	log.Printf("ProcessToken -- result: %s\n\n", result)
	log.Println(stack)

	var err error
	tokenString := strings.TrimSpace(token.String())
	tokenBytes := []byte(tokenString)

	if len(tokenString) > 2 {

		//Check for Tag
		if bytes.HasPrefix(tokenBytes, []byte("<")) && bytes.HasSuffix(tokenBytes, []byte(">")) {
			switch {
			case bytes.HasPrefix(tokenBytes, OpenH1Tag.Byte()):
				stack = stack.Push(OpenH1Tag)
				tempt = fmt.Sprintf("%s", tempt)

			case bytes.HasPrefix(tokenBytes, CloseH1Tag.Byte()):
				stack, _, err = stack.Pop()
				if err != nil {
					return stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken, err
				}
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s\n", result, tempt)
					tempt = ""
				}

			case bytes.HasPrefix(tokenBytes, OpenH2Tag.Byte()):
				stack = stack.Push(OpenH2Tag)
				tempt = fmt.Sprintf("%s", tempt)

			case bytes.HasPrefix(tokenBytes, CloseH2Tag.Byte()):
				stack, _, err = stack.Pop()
				if err != nil {
					return stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken, err
				}
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s\n", result, tempt)
					tempt = ""
				}

			case bytes.HasPrefix(tokenBytes, OpenH3Tag.Byte()):
				stack = stack.Push(OpenH3Tag)
				tempt = fmt.Sprintf("%s", tempt)

			case bytes.HasPrefix(tokenBytes, CloseH3Tag.Byte()):
				stack, _, err = stack.Pop()
				if err != nil {
					return stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken, err
				}
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s\n", result, tempt)
					tempt = ""
				}

			case bytes.HasPrefix(tokenBytes, OpenH4Tag.Byte()):
				stack = stack.Push(OpenH4Tag)
				tempt = fmt.Sprintf("%s", tempt)

			case bytes.HasPrefix(tokenBytes, CloseH4Tag.Byte()):
				stack, _, err = stack.Pop()
				if err != nil {
					return stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken, err
				}
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s\n", result, tempt)
					tempt = ""
				}

			case bytes.HasPrefix(tokenBytes, OpenH5Tag.Byte()):
				stack = stack.Push(OpenH5Tag)
				tempt = fmt.Sprintf("%s", tempt)

			case bytes.HasPrefix(tokenBytes, CloseH5Tag.Byte()):
				stack, _, err = stack.Pop()
				if err != nil {
					return stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken, err
				}
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s\n", result, tempt)
					tempt = ""
				}

			case bytes.HasPrefix(tokenBytes, OpenH6Tag.Byte()):
				stack = stack.Push(OpenH6Tag)
				tempt = fmt.Sprintf("%s", tempt)

			case bytes.HasPrefix(tokenBytes, CloseH6Tag.Byte()):
				stack, _, err = stack.Pop()
				if err != nil {
					return stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken, err
				}
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s\n", result, tempt)
					tempt = ""
				}

			case bytes.HasPrefix(tokenBytes, OpenPTag.Byte()):
				log.Println("Case: <p>")
				stack = stack.Push(OpenPTag)

				//Set Parent Tag
				parentTag = OpenPTag

			case bytes.HasPrefix(tokenBytes, ClosePTag.Byte()):
				log.Println("Case: </p>")

				//Unset Parent Tag
				parentTag = Tag("")

				stack, _, err = stack.Pop()
				if err != nil {
					return stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken, err
				}
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s\n\n", result, tempt)
					tempt = ""
				}

			case bytes.HasPrefix(tokenBytes, OpenOLTag.Byte()):
				log.Println("Case: <ol")
				log.Print("Attributes: ")
				log.Print(token.Attr)

				//Set Parent Tag
				parentTag = OpenOLTag

				for _, attr := range token.Attr {
					if strings.EqualFold(attr.Key, "start") {
						listCount, err = strconv.Atoi(attr.Val)
						if err != nil {
							return stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken, err
						}
					}
				}
				stack = stack.Push(OpenOLTag)

			case bytes.HasPrefix(tokenBytes, CloseOLTag.Byte()):
				log.Println("Case: </ol>")

				//UnSet Parent Tag
				parentTag = Tag("")

				stack, _, err = stack.Pop()
				if err != nil {
					return stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken, err
				}
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s\n", result, tempt)
					tempt = ""
				}

				//Reset variable
				listCount = 1

			case bytes.HasPrefix(tokenBytes, OpenULTag.Byte()):
				log.Println("Case: <ul>")

				//Set Parent Tag
				parentTag = OpenULTag

				stack = stack.Push(OpenULTag)

			case bytes.HasPrefix(tokenBytes, CloseULTag.Byte()):
				log.Println("Case: </ul>")

				//UnSet Parent Tag
				parentTag = Tag("")

				stack, _, err = stack.Pop()
				if err != nil {
					return stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken, err
				}
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s\n", result, tempt)
					tempt = ""
				}

			case bytes.HasPrefix(tokenBytes, OpenLITag.Byte()):
				log.Println("Case: <li>")
				if strings.EqualFold(parentTag.String(), OpenOLTag.String()) {
					tempt = fmt.Sprintf("%s  %d. ", tempt, listCount)
					listCount++
				} else if strings.EqualFold(parentTag.String(), OpenULTag.String()) {
					tempt = fmt.Sprintf("%s  * ", tempt)
				}

			case bytes.HasPrefix(tokenBytes, CloseLITag.Byte()):
				log.Println("Case: </li>")
				tempt = fmt.Sprintf("%s\n", tempt)

			case bytes.HasPrefix(tokenBytes, OpenATag.Byte()):
				log.Println("Case: <a")
				log.Print("Attributes: ")
				log.Print(token.Attr)

				stack = stack.Push(OpenATag)

				//Capture the href
				for _, attr := range token.Attr {
					if strings.EqualFold(attr.Key, "href") {
						links = append(links, attr.Val)
					}
				}

				//Checking if I need to add a space
				if strings.EqualFold(parentTag.String(), OpenPTag.String()) {
					tempt = fmt.Sprintf("%s ", tempt)
				}

			case bytes.HasPrefix(tokenBytes, CloseATag.Byte()):
				log.Println("Case: </a>")

				//Checking if I need to add a space
				if strings.EqualFold(parentTag.String(), OpenPTag.String()) {
					tempt = fmt.Sprintf("%s[%d] ", tempt, len(links))
				}

				stack, _, err = stack.Pop()
				if err != nil {
					return stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken, err
				}
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s[%d]", result, tempt, len(links))
					tempt = ""
				}

			case bytes.HasPrefix(tokenBytes, OpenTableTag.Byte()):
				log.Println("Case: <table>")
				stack = stack.Push(OpenTableTag)

				//Set Parent Tag
				parentTag = OpenTableTag

			case bytes.HasPrefix(tokenBytes, CloseTableTag.Byte()):
				log.Println("Case: </table>")

				//Unset Parent Tag
				parentTag = Tag("")

				stack, _, err = stack.Pop()
				if err != nil {
					return stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken, err
				}

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
							tableString = fmt.Sprintf("%s|%s", tableString)
						}
						count++
					}
					//I Need this to stop one iterval prior to the end
					if y < tableColumns {
						tableString = fmt.Sprintf("%s|\n", tableString)
					}
				}

				//Step 3: add table to result
				result += fmt.Sprintf("%s\n", tableString)

			case bytes.HasPrefix(tokenBytes, OpenTRTag.Byte()):
				log.Println("Case: <tr>")
				tableColumns++

			case bytes.HasPrefix(tokenBytes, CloseTRTag.Byte()):
				log.Println("Case: </tr>")

			case bytes.HasPrefix(tokenBytes, OpenTHTag.Byte()):
				log.Println("Case: <th>")
				tableRows++

			case bytes.HasPrefix(tokenBytes, CloseTHTag.Byte()):
				log.Println("Case: </th>")

			case bytes.HasPrefix(tokenBytes, OpenTDTag.Byte()):
				log.Println("Case: <td>")

			case bytes.HasPrefix(tokenBytes, CloseTDTag.Byte()):
				log.Println("Case: </td>")

			case bytes.HasPrefix(tokenBytes, OpenStyleTag.Byte()):
				log.Println("Case: <style>")

			case bytes.HasPrefix(tokenBytes, CloseStyleTag.Byte()):
				log.Println("Case: </style>")

			case bytes.HasPrefix(tokenBytes, OpenHeadTag.Byte()):
				log.Println("Case: <head>")
				//Turn on ignore token
				ignoreToken = true

			case bytes.HasPrefix(tokenBytes, CloseHeadTag.Byte()):
				log.Println("Case: </head>")
				//Turn off ignore token
				ignoreToken = false

			default:
			}
		} else {
			if strings.EqualFold(parentTag.String(), OpenTableTag.String()) {
				tableItems = append(tableItems, tokenString)
			} else if ignoreToken {
				//Do nothing
			} else {
				tempt += tokenString
			}
		}
	}

	log.Println("Returned:")
	log.Printf("Stack: ")
	log.Print(stack)
	log.Printf("Tempt: %s\n", tempt)
	log.Printf("Result: %s\n", result)
	log.Printf("TableRows: %d, TableColumns %d\n", tableRows, tableColumns)
	log.Print("TableItems: ")
	log.Println(tableItems)
	return stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken, nil
}

//Translate -- translates html to text
func Translate(reader io.Reader) (string, []string, error) {
	var result string
	var tempt string
	var err error
	var parentTag Tag
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
			log.Fatalf("html Parser err token: %s", data.Err())
		}
		// Process the current token.
		token := data.Token()

		stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken, err = processToken(token, stack, tempt, result, links, parentTag, listCount, tableRows, tableColumns, tableItems, ignoreToken)
		if err != nil {
			log.Fatalf("ProcessToken had an error: %s", err.Error())
		}
	}

	return result, links, nil
}
