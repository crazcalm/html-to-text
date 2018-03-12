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

func processToken(token html.Token, stack Stack, tempt, result string, links []string, parentTag Tag, listCount int) (Stack, string, string, []string, Tag, int, error) {
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
					return stack, tempt, result, links, parentTag, listCount, err
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
					return stack, tempt, result, links, parentTag, listCount, err
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
					return stack, tempt, result, links, parentTag, listCount, err
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
					return stack, tempt, result, links, parentTag, listCount, err
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
					return stack, tempt, result, links, parentTag, listCount, err
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
					return stack, tempt, result, links, parentTag, listCount, err
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
					return stack, tempt, result, links, parentTag, listCount, err
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
							return stack, tempt, result, links, parentTag, listCount, err
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
					return stack, tempt, result, links, parentTag, listCount, err
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
					return stack, tempt, result, links, parentTag, listCount, err
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
					return stack, tempt, result, links, parentTag, listCount, err
				}
				if len(stack) == 0 {
					result = fmt.Sprintf("%s%s[%d]", result, tempt, len(links))
					tempt = ""
				}

			default:
			}
		} else {
			tempt += tokenString
		}
	}

	log.Println("Returned:")
	log.Printf("Stack: ")
	log.Print(stack)
	log.Printf("Tempt: %s\n", tempt)
	log.Printf("Result: %s\n", result)
	return stack, tempt, result, links, parentTag, listCount, nil
}

//Translate -- translates html to text
func Translate(reader io.Reader) (string, error) {
	var result string
	var tempt string
	var err error
	var parentTag Tag
	var links []string

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

		stack, tempt, result, links, parentTag, listCount, err = processToken(token, stack, tempt, result, links, parentTag, listCount)
		if err != nil {
			log.Fatalf("ProcessToken had an error: %s", err.Error())
		}
	}

	return result, nil
}
