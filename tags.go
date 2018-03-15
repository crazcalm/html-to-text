package htmltotext

//Tag -- A struct for HTML tags
type Tag string

func (t Tag) String() string {
	return string(t)
}

//Byte -- returns the tag in a slice of bytes
func (t Tag) Byte() []byte {
	return []byte(t.String())
}

const (
	//OpenDivTag -- Open html div tag
	OpenDivTag = Tag("<div")
	//CloseDivTag -- Close html div tag
	CloseDivTag = Tag("</div>")

	//OpenH1Tag -- Open html h1 tag
	OpenH1Tag = Tag("<h1>")
	//CloseH1Tag -- Closed html h1 tag
	CloseH1Tag = Tag("</h1>")

	//OpenH2Tag -- Open html h2 tag
	OpenH2Tag = Tag("<h2>")
	//CloseH2Tag -- Close html h2 tag
	CloseH2Tag = Tag("</h2>")

	//OpenH3Tag -- Open html h3 tag
	OpenH3Tag = Tag("<h3>")
	//CloseH3Tag -- Close html h3 tag
	CloseH3Tag = Tag("</h3>")

	//OpenH4Tag -- Open html h4 tag
	OpenH4Tag = Tag("<h4>")
	//CloseH4Tag -- Close html h4 tag
	CloseH4Tag = Tag("</h4>")

	//OpenH5Tag -- Open html h5 tag
	OpenH5Tag = Tag("<h5>")
	//CloseH5Tag -- Close html h5 tag
	CloseH5Tag = Tag("</h5>")

	//OpenH6Tag -- Open html h6 tag
	OpenH6Tag = Tag("<h6>")
	//CloseH6Tag -- Close html h6 tag
	CloseH6Tag = Tag("</h6>")

	//OpenPTag -- Open html p tag
	OpenPTag = Tag("<p>")
	//ClosePTag -- Closed html p tag
	ClosePTag = Tag("</p>")

	//OpenOLTag -- Open html ordered list tag
	OpenOLTag = Tag("<ol")
	//CloseOLTag -- Close html ordered list tag
	CloseOLTag = Tag("</ol>")

	//OpenULTag -- Open html unordered list tag
	OpenULTag = Tag("<ul>")
	//CloseULTag -- Closed html unordered list tag
	CloseULTag = Tag("</ul>")

	//OpenLITag -- Open html list tag
	OpenLITag = Tag("<li>")
	//CloseLITag -- Close html list tag
	CloseLITag = Tag("</li>")

	//OpenATag -- Open html a tag
	OpenATag = Tag("<a")
	//CloseATag -- Close html a tag
	CloseATag = Tag("</a>")

	//OpenTableTag -- Open html table tag
	OpenTableTag = Tag("<table")
	//CloseTableTag -- Close html table tag
	CloseTableTag = Tag("</table>")

	//OpenTRTag -- Open html tr tag
	OpenTRTag = Tag("<tr>")
	//CloseTRTag -- Close html tr tag
	CloseTRTag = Tag("</tr>")

	//OpenTHeadTag -- Open html thead tag
	OpenTHeadTag = Tag("<thead")
	//CloseTHeadTag -- Close html thead tag
	CloseTHeadTag = Tag("</thead>")

	//OpenTHTag -- Open html th tag
	OpenTHTag = Tag("<th")
	//CloseTHTag -- Close html th tag
	CloseTHTag = Tag("</th>")

	//OpenTDTag -- Open html td tag
	OpenTDTag = Tag("<td>")
	//CloseTDTag -- Close html td tag
	CloseTDTag = Tag("</td>")

	//OpenStyleTag -- Open html style tag
	OpenStyleTag = Tag("<style>")
	//CloseStyleTag -- Close html style tag
	CloseStyleTag = Tag("</style>")

	//OpenHeadTag -- Open html head tag
	OpenHeadTag = Tag("<head>")
	//CloseHeadTag -- Close html head tag
	CloseHeadTag = Tag("</head>")

	//BreakTag -- html br tag
	BreakTag = Tag("<br")
)
