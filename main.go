package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BoomerangMessaging/notiongo/cache"
)

// Represents the parent object (in this case, a page)
type Parent struct {
	Type   string `json:"type"`
	PageID string `json:"page_id"`
}

// Represents the user who created/edited the block
type User struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

type Mention struct {
	Type string `json:"type"`
	Page struct {
		ID string `json:"id"`
	} `json:"page"`
}

// Represents a text object (for paragraphs, headings, etc.)
type RichText struct {
	Type        string      `json:"type"`
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text"`
	Mention     *Mention    `json:"mention,omitempty"`
	Href        *string     `json:"href,omitempty"`
}

// Represents the content of a text object
type Text struct {
	Content string  `json:"content"`
	Link    *string `json:"link,omitempty"`
}

// Represents text styling (annotations)
type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

// Represents a heading block
type Heading struct {
	RichText     []RichText `json:"rich_text"`
	IsToggleable bool       `json:"is_toggleable"`
	Color        string     `json:"color"`
}

// Represents a paragraph block
type Paragraph struct {
	RichText []RichText `json:"rich_text"`
	Color    string     `json:"color"`
}

// Represents a divider block
type Divider struct{}

// Represents a table block
type Table struct {
	TableWidth      int         `json:"table_width"`
	HasColumnHeader bool        `json:"has_column_header"`
	HasRowHeader    bool        `json:"has_row_header"`
	Rows            []*TableRow `json:"rows,omitempty"`
}

type TableRow struct {
	Cells [][]TableCell `json:"cells"`
}

type TableCell struct {
	Type        string      `json:"type"`
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text"`
	Href        string      `json:"href,omitempty"`
	Mention     *Mention    `json:"mention,omitempty"`
}
type NumberedListItem struct {
	RichText []RichText `json:"rich_text"`
	Color    string     `json:"color"`
}

type ToDoItem struct {
	RichText []RichText `json:"rich_text"`
	Checked  bool       `json:"checked"`
	Color    string     `json:"color"`
}

type CodeBlock struct {
	RichText []RichText `json:"rich_text"`
	Language string     `json:"language"`
}

type Quote struct {
	RichText []RichText `json:"rich_text"`
	Color    string     `json:"color"`
}

type Callout struct {
	RichText []RichText `json:"rich_text"`
	Icon     Icon       `json:"icon"`
	Color    string     `json:"color"`
}

type Icon struct {
	Type  string `json:"type"`
	Emoji string `json:"emoji,omitempty"`
}

type Image struct {
	Type     string `json:"type"`
	File     *File  `json:"file,omitempty"`
	External *Link  `json:"external,omitempty"`
}

type File struct {
	URL        string `json:"url"`
	ExpiryTime string `json:"expiry_time"`
}

type Link struct {
	URL string `json:"url"`
}

type LinkToPage struct {
	Type   string `json:"type"`
	PageID string `json:"page_id"`
}
type Bookmark struct {
	URL     string     `json:"url"`
	Caption []RichText `json:"caption"`
}

type NotionBlock struct {
	Object           string            `json:"object"`
	ID               string            `json:"id"`
	Parent           Parent            `json:"parent"`
	CreatedTime      string            `json:"created_time"`
	LastEditedTime   string            `json:"last_edited_time"`
	CreatedBy        User              `json:"created_by"`
	LastEditedBy     User              `json:"last_edited_by"`
	HasChildren      bool              `json:"has_children"`
	Archived         bool              `json:"archived"`
	InTrash          bool              `json:"in_trash"`
	Type             string            `json:"type"`
	Annotations      Annotations       `json:"annotations"`
	Heading1         *Heading          `json:"heading_1,omitempty"`
	Heading2         *Heading          `json:"heading_2,omitempty"`
	Heading3         *Heading          `json:"heading_3,omitempty"`
	BulltedListItem  *Paragraph        `json:"bulleted_list_item,omitempty"`
	Paragraph        *Paragraph        `json:"paragraph,omitempty"`
	Divider          *Divider          `json:"divider,omitempty"`
	Table            *Table            `json:"table,omitempty"`
	NumberedListItem *NumberedListItem `json:"numbered_list_item,omitempty"`
	ToDoItem         *ToDoItem         `json:"to_do,omitempty"`
	Code             *CodeBlock        `json:"code,omitempty"`
	Quote            *Quote            `json:"quote,omitempty"`
	Callout          *Callout          `json:"callout,omitempty"`
	Image            *Image            `json:"image,omitempty"`
	File             *File             `json:"file,omitempty"`
	Bookmark         *Bookmark         `json:"bookmark,omitempty"`
	TableRows        *TableRow         `json:"table_row,omitempty"`
	LinkToPage       *LinkToPage       `json:"link_to_page,omitempty"`
	ChildPage        *ChildPage        `json:"child_page,omitempty"`
}

type ChildPage struct {
	Title string `json:"title"`
}

type TableRowBlock struct {
	Object         string    `json:"object"`
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	CreatedTime    string    `json:"created_time"`
	LastEditedTime string    `json:"last_edited_time"`
	HasChildren    bool      `json:"has_children"`
	TableRow       *TableRow `json:"table_row"`
}

type NotionBlockChildrenResponse struct {
	Results    []NotionBlock `json:"results"`
	NextCursor string        `json:"next_cursor"`
	HasMore    bool          `json:"has_more"`
}

type NotionPage struct {
	Object     string                 `json:"object"`
	ID         string                 `json:"id"`
	Properties map[string]interface{} `json:"properties"`
}

type NotionQueryResponse struct {
	Results    []NotionPage `json:"results"`
	NextCursor string       `json:"next_cursor"`
	HasMore    bool         `json:"has_more"`
}

var conf struct {
	AssetsDir      string
	OutputDir      string
	APIToken       string
	DocsRoot       string
	slugRegistered []string
}

func stringExists(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// Fetch children blocks of a block (pages, databases, etc.)
func fetchChildren(token string, blockID string, cursor string) (NotionBlockChildrenResponse, error) {

	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children?page_size=100", blockID)

	cache := cache.NewCache()
	// Try to get the response from the cache first
	if cachedResponse, found := cache.Get(url); found {
		fmt.Println("Cache hit:", cachedResponse)
		return cachedResponse.(NotionBlockChildrenResponse), nil
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NotionBlockChildrenResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Notion-Version", "2022-06-28")

	if cursor != "" {
		req.URL.RawQuery = fmt.Sprintf("start_cursor=%s", cursor)
	}

	resp, err := client.Do(req)
	if err != nil {
		return NotionBlockChildrenResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NotionBlockChildrenResponse{}, err
	}

	if resp.StatusCode == 429 {
		fmt.Println("Rate limited. Waiting for retry...")
		time.Sleep(3 * time.Second)
		return fetchChildren(token, blockID, cursor)
	}

	var data NotionBlockChildrenResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return NotionBlockChildrenResponse{}, err
	}

	// Cache the response with a 5-second TTL
	cache.Set(url, data, 600*time.Second)

	return data, nil
}

// Fetch pages from a database
func fetchPagesFromDatabase(token string, databaseID string, cursor string) (NotionQueryResponse, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", databaseID)

	cache := cache.NewCache()
	// Try to get the response from the cache first
	if cachedResponse, found := cache.Get(url); found {
		fmt.Println("Cache hit:", cachedResponse)
		return cachedResponse.(NotionQueryResponse), nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return NotionQueryResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Notion-Version", "2022-06-28")

	if cursor != "" {
		reqBody := fmt.Sprintf(`{"start_cursor":"%s"}`, cursor)
		req.Body = io.NopCloser(strings.NewReader(reqBody))
	}

	resp, err := client.Do(req)
	if err != nil {
		return NotionQueryResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NotionQueryResponse{}, err
	}

	if resp.StatusCode == 429 {
		fmt.Println("Rate limited. Waiting for retry...")
		time.Sleep(3 * time.Second)
		return fetchPagesFromDatabase(token, databaseID, cursor)
	}

	var data NotionQueryResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return NotionQueryResponse{}, err
	}

	// Cache the response with a 5-second TTL
	cache.Set(url, data, 600*time.Second)

	return data, nil
}

// Check if either the directory or file exists
func namedDirOrFileExists(rootDir, name string) (bool, error) {
	// Construct paths for the directory and the file
	dirPath := filepath.Join(rootDir, name)        // /root/page1234
	filePath := filepath.Join(rootDir, name+".md") // /root/page1234.md

	// Check if the directory exists
	if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
		return true, nil // Directory exists
	}

	// Check if the file exists
	if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
		return true, nil // File exists
	}

	// Neither directory nor file exists
	return false, nil
}

// Fetch content of a page by retrieving its blocks
func fetchPage(token string, pageID string) (*NotionPage, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/pages/%s", pageID)

	cache := cache.NewCache()
	// Try to get the response from the cache first
	if cachedResponse, found := cache.Get(url); found {
		fmt.Println("Cache hit:", cachedResponse)
		return cachedResponse.(*NotionPage), nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Notion-Version", "2022-06-28")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 429 {
		fmt.Println("Rate limited. Waiting for retry...")
		time.Sleep(3 * time.Second)
		return fetchPage(token, pageID)
	}

	var response NotionPage
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Cache the response with a 5-second TTL
	cache.Set(url, response, 600*time.Second)

	return &response, nil
}

// Fetch content of a page by retrieving its blocks
func fetchPageContent(token string, pageID string) ([]NotionBlock, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children", pageID)

	cache := cache.NewCache()
	// Try to get the response from the cache first
	if cachedResponse, found := cache.Get(url); found {
		fmt.Println("Cache hit:", cachedResponse)
		return cachedResponse.([]NotionBlock), nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Notion-Version", "2022-06-28")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 429 {
		fmt.Println("Rate limited. Waiting for retry...")
		time.Sleep(3 * time.Second)
		return fetchPageContent(token, pageID)
	}

	var response NotionBlockChildrenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Cache the response with a 5-second TTL
	cache.Set(url, response.Results, 600*time.Second)

	return response.Results, nil
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func downloadImage(url string, filepath string) (string, error) {
	// Make the HTTP request to download the image
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}

	// Get the Content-Type header to determine the file extension
	contentType := resp.Header.Get("Content-Type")
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil || len(exts) == 0 {
		return "", fmt.Errorf("failed to determine file extension for content type: %s", contentType)
	}

	// Create a unique file name with random string and timestamp
	timestamp := time.Now().UnixNano()
	randomStr := randomString(8)
	filename := fmt.Sprintf("%s_%d%s", randomStr, timestamp, exts[0])

	img := fmt.Sprintf("%s/%s", filepath, filename)

	// Create the file with the appropriate extension
	out, err := os.Create(img)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copy the image data to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func formatBlockHTML(rt RichText) string {

	rt.PlainText = strings.ReplaceAll(rt.PlainText, "·", "-")
	rt.PlainText = strings.ReplaceAll(rt.PlainText, "\n", "<br />")

	// format text
	if rt.Annotations.Bold {
		rt.PlainText = "<strong>" + rt.PlainText + "</strong>"
	}
	if rt.Annotations.Italic {
		rt.PlainText = "<em>" + rt.PlainText + "</em>"
	}
	if rt.Annotations.Underline {
		rt.PlainText = "<u>" + rt.PlainText + "</u>"
	}
	if rt.Annotations.Strikethrough {
		rt.PlainText = "<del>" + rt.PlainText + "</del>"
	}
	if rt.Annotations.Code {
		rt.PlainText = "<code>" + rt.PlainText + "</code>"
	}

	return rt.PlainText

}

// Convert Notion blocks to Markdown content
func blocksToMarkdown(token string, blocks []NotionBlock, isChildren bool) string {
	var markdownBuilder strings.Builder

	for _, block := range blocks {
		var plainText string
		switch block.Type {
		case "paragraph":
			for _, t := range block.Paragraph.RichText {
				if t.Type == "mention" {
					page, err := fetchPage(token, t.Mention.Page.ID)
					if err != nil {
						log.Printf("[ERROR] while fetching mention_to_page %v", err)
						continue
					} else {
						title, slug, _ := extractPageProperties(*page)
						if len(slug) > 0 && slug[0:1] == "/" {
							slug = conf.DocsRoot + slug
						}
						plainText += fmt.Sprintf("[%s](%s)", title, slug)
					}
				} else {
					plainText += formatBlockHTML(t)
				}
			}
			markdownBuilder.WriteString(plainText + "  \n")
		case "heading_1":
			for _, t := range block.Heading1.RichText {
				plainText += t.PlainText
			}
			markdownBuilder.WriteString("# " + plainText + "  \n")
		case "heading_2":
			for _, t := range block.Heading2.RichText {
				plainText += t.PlainText
			}
			markdownBuilder.WriteString("## " + plainText + "  \n")
		case "heading_3":
			for _, t := range block.Heading3.RichText {
				plainText += t.PlainText
			}
			markdownBuilder.WriteString("### " + plainText + "  \n")
		case "bulleted_list_item":
			var PlainText string
			for _, t := range block.BulltedListItem.RichText {
				PlainText += formatBlockHTML(t)
			}
			content := "- " + PlainText + "  \n"
			if isChildren {
				content = fmt.Sprintf("\t%s", content)
			}
			markdownBuilder.WriteString(content)

		case "table":
			tableRows, err := fetchTableContent(token, block.ID)
			if err != nil {
				log.Printf("Error fetching table content: %v", err)
				markdownBuilder.WriteString("[Error: Could not fetch table content]\n")
			} else {
				markdownBuilder.WriteString(renderTable(block.Table, tableRows) + "  \n")
			}
		case "table_row":
			var allRows []TableRow
			allRows = append(allRows, *block.TableRows)
			markdownBuilder.WriteString(renderTable(block.Table, allRows) + "  \n")

		case "divider":
			markdownBuilder.WriteString("\n--- \n")
		case "numbered_list_item":
			for _, t := range block.NumberedListItem.RichText {
				plainText += formatBlockHTML(t)
			}

			content := "1. " + plainText + "  \n"
			if isChildren {
				content = fmt.Sprintf("\t%s", content)
			}
			markdownBuilder.WriteString(content)

		case "to_do":
			checkbox := "[ ]"
			if block.ToDoItem.Checked {
				checkbox = "[x]"
			}
			for _, t := range block.ToDoItem.RichText {
				plainText += formatBlockHTML(t)
			}
			markdownBuilder.WriteString(checkbox + " " + plainText + "  \n")
		case "code":
			for _, t := range block.Code.RichText {
				plainText += t.PlainText
			}
			markdownBuilder.WriteString("```" + block.Code.Language + "  \n" + plainText + "  \n```\n")
		case "quote":
			for _, t := range block.Quote.RichText {
				plainText += formatBlockHTML(t)
			}
			markdownBuilder.WriteString("> " + plainText + "  \n")
		case "callout":
			icon := ""
			if block.Callout.Icon.Type == "emoji" {
				icon = block.Callout.Icon.Emoji + " "
			}
			for _, t := range block.Callout.RichText {
				plainText += formatBlockHTML(t)
			}
			markdownBuilder.WriteString("> " + icon + plainText + "  \n")
		case "image":
			caption := ""
			url := ""
			if block.Image.Type == "file" {
				url = block.Image.File.URL
			} else if block.Image.Type == "external" {
				url = block.Image.External.URL
			}
			staticDir := fmt.Sprintf("%s/docs-images", conf.AssetsDir)

			if _, err := os.Stat(staticDir); os.IsNotExist(err) {
				if err := os.MkdirAll(staticDir, os.ModePerm); err != nil {
					log.Println("failed to create subdirectory ", err)
				}
			}

			filename, err := downloadImage(url, staticDir)
			if err != nil {
				log.Println("error occured while downloading image", err)
			} else {
				caption = filename
				url = fmt.Sprintf("/docs-images/%s", filename)
			}

			markdownBuilder.WriteString(fmt.Sprintf("![%s](%s)\n\n", caption, url))
		case "file":
			markdownBuilder.WriteString(fmt.Sprintf("[File](%s)  \n", block.File.URL))
		case "bookmark":
			caption := ""
			for _, t := range block.Bookmark.Caption {
				caption += t.PlainText
			}
			markdownBuilder.WriteString(fmt.Sprintf("[%s](%s)  \n", caption, block.Bookmark.URL))
		case "link_to_page":
			page, err := fetchPage(token, block.LinkToPage.PageID)
			if err != nil {
				log.Printf("[ERROR] while fetching link_to_page %v", err)
				continue
			} else {
				title, slug, _ := extractPageProperties(*page)

				if len(slug) > 0 && slug[0:1] == "/" {
					slug = conf.DocsRoot + slug
				}

				markdownBuilder.WriteString(fmt.Sprintf("[%s](%s)<br/>", title, slug))
			}
		case "unsupported":
		default:
			markdownBuilder.WriteString(fmt.Sprintf("[Unsupported block type: %s]  \n", block.Type))
		}

		if block.HasChildren {
			blocks, err := fetchPageContent(token, block.ID)
			if err != nil {
				log.Println("[ERROR] failed to fetch ")
			} else {
				// Convert blocks to markdown content
				contentMarkdown := blocksToMarkdown(token, blocks, true)
				markdownBuilder.WriteString(contentMarkdown)
			}

		}
	}

	return markdownBuilder.String()
}

func fetchTableContent(token string, tableBlockID string) ([]TableRow, error) {
	var allRows []TableRow
	var nextCursor string
	hasMore := true

	for hasMore {
		response, err := fetchChildren(token, tableBlockID, nextCursor)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch table content: %v", err)
		}

		for _, result := range response.Results {
			if result.Type == "table_row" {
				allRows = append(allRows, *result.TableRows)
			}
		}

		hasMore = response.HasMore
		nextCursor = response.NextCursor
		time.Sleep(1 * time.Second) // Add delay to respect rate limits
	}

	return allRows, nil
}

func renderTable(table *Table, rows []TableRow) string {
	if table == nil || len(rows) == 0 {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(`<table>`)
	// var tableSep string
	for i, row := range rows {
		sb.WriteString("<tr>")
		for _, cell := range row.Cells {
			if i == 0 {
				sb.WriteString("<th>" + renderTableCell(cell) + "</th>")
			} else {
				sb.WriteString("<td>" + renderTableCell(cell) + "</td>")
			}
		}
		sb.WriteString("</tr>")
	}
	sb.WriteString("</table>")

	return sb.String()
}

// Helper function to render a table cell
func renderTableCell(cell []TableCell) string {
	var cellContent string
	for _, rt := range cell {

		if rt.Type == "mention" {
			page, err := fetchPage(conf.APIToken, rt.Mention.Page.ID)
			if err != nil {
				log.Printf("[ERROR] while fetching mention_to_page %v", err)
				continue
			} else {
				title, slug, _ := extractPageProperties(*page)
				if len(slug) > 0 && slug[0:1] == "/" {
					slug = conf.DocsRoot + slug
				}
				rt.PlainText = fmt.Sprintf("[%s](%s)", title, slug)
			}
		}

		rt.PlainText = strings.ReplaceAll(rt.PlainText, "·", "-")
		rt.PlainText = strings.ReplaceAll(rt.PlainText, "\n", "<br />")

		// format text
		if rt.Annotations.Bold {
			rt.PlainText = "<strong>" + rt.PlainText + "</strong>"
		}
		if rt.Annotations.Italic {
			rt.PlainText = "<em>" + rt.PlainText + "</em>"
		}
		if rt.Annotations.Underline {
			rt.PlainText = "<u>" + rt.PlainText + "</u>"
		}
		if rt.Annotations.Strikethrough {
			rt.PlainText = "<del>" + rt.PlainText + "</del>"
		}
		if rt.Annotations.Code {
			rt.PlainText = "<code>" + rt.PlainText + "</code>"
		}

		cellContent += rt.PlainText
		// if len(cell) > 1 && (len(cell)-1) != cellIndex {
		// 	cellContent += "<br/>"
		// }

	}
	cellContent = strings.TrimSuffix(cellContent, "\n")
	// Escape pipe characters in cell content
	cellContent = strings.ReplaceAll(cellContent, "|", "&#124;")
	return cellContent
}

// Helper function to extract property values from a page
func extractPageProperties(page NotionPage) (title string, slug string, keywords string) {
	// Title
	if titleProp, ok := page.Properties["Name"].(map[string]interface{}); ok {
		title = titleProp["title"].([]interface{})[0].(map[string]interface{})["plain_text"].(string)
	}

	// Slug
	if slugProp, ok := page.Properties["Slug"].(map[string]interface{}); ok {
		slug = title
		if len(slugProp["rich_text"].([]interface{})) > 0 {
			slug = slugProp["rich_text"].([]interface{})[0].(map[string]interface{})["plain_text"].(string)

		}
		slug = strings.ReplaceAll(slug, " ", "-")
	}

	// Keywords
	if keywordsProp, ok := page.Properties["Keywords"].(map[string]interface{}); ok {
		keywords = ""
		if len(keywordsProp["rich_text"].([]interface{})) > 0 {
			keywords = keywordsProp["rich_text"].([]interface{})[0].(map[string]interface{})["plain_text"].(string)
		}
	}

	return title, slug, keywords
}

func extractPageRelations(page NotionPage) (parentId string, childPages []string) {

	if parent, ok := page.Properties["Parent"].(map[string]interface{}); ok {
		if len(parent["relation"].([]interface{})) > 0 {
			parentId = parent["relation"].([]interface{})[0].(map[string]interface{})["id"].(string)
		}
	}
	if subItems, ok := page.Properties["Sub-Items"].(map[string]interface{}); ok {
		if len(subItems["relation"].([]interface{})) > 0 {
			for _, item := range subItems["relation"].([]interface{}) {
				childPages = append(childPages, item.(map[string]interface{})["id"].(string))
			}
		}
	}

	return parentId, childPages
}

// Convert a page to markdown, including content
func pageToMarkdown(token string, page NotionPage, position int) (string, error) {
	title, slug, keywords := extractPageProperties(page)

	// Fetch page content (blocks)
	blocks, err := fetchPageContent(token, page.ID)
	if err != nil {
		return "", err
	}

	// Convert blocks to markdown content
	contentMarkdown := blocksToMarkdown(token, blocks, false)

	// Format keywords for markdown
	keywordString := "[" + keywords + "]"

	slug = strings.ReplaceAll(slug, "(", "")
	slug = strings.ReplaceAll(slug, ")", "")

	if stringExists(conf.slugRegistered, slug) {
		slug += "-dup"
	}
	conf.slugRegistered = append(conf.slugRegistered, slug)

	// Template for markdown output
	return fmt.Sprintf(`---
title: %s
slug: %s
tags: %s
sidebar_position: %d
---

%s
`, title, slug, keywordString, position, contentMarkdown), nil
}

// Write markdown to file
func writeMarkdown(outputDir string, token string, page NotionPage, position int) error {
	markdown, err := pageToMarkdown(token, page, position)
	if err != nil {
		return err
	}
	// title, _, _ := extractPageProperties(page)

	// sub := strings.Split(slug, "/")
	dir := outputDir

	_, childPages := extractPageRelations(page)
	// if parent != "" {
	// 	dir = fmt.Sprintf("%s/%s", dir, parent)
	// 	fmt.Printf("directory to be created %s\n", dir)
	// 	if _, err := os.Stat(dir); os.IsNotExist(err) {
	// 		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
	// 			log.Println("failed to create subdirectory ", err)
	// 		}
	// 	}
	// }

	HasChildren := false
	if len(childPages) > 0 {
		dir = fmt.Sprintf("%s/%s", dir, page.ID)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				log.Println("failed to create subdirectory ", err)
			}
		}
		for cPageIndex, child := range childPages {
			childPage, err := fetchPage(token, child)
			if err != nil {
				fmt.Printf("failed to fetch child id %s", child)
			} else {
				if err := writeMarkdown(dir, token, *childPage, cPageIndex); err != nil {
					fmt.Printf("failed to write markdown for child page %s", childPage.ID)
					continue
				}
				HasChildren = true
			}
		}
	}

	title, _, _ := extractPageProperties(page)
	title = strings.ReplaceAll(title, `\`, `\\`)
	title = strings.ReplaceAll(title, `"`, `\"`)
	title = strings.ReplaceAll(title, "\n", `\n`)
	title = strings.ReplaceAll(title, "\t", `\t`)
	title = strings.ReplaceAll(title, "\r", `\r`)
	title = strings.ReplaceAll(title, "\b", `\b`)
	title = strings.ReplaceAll(title, "\f", `\f`)

	categoryJson := fmt.Sprintf(`{
	"label": "%s",
	"position": %d
}`, title, position)

	pagename := page.ID
	if HasChildren {
		pagename = "index"
		os.WriteFile(fmt.Sprintf("%s/_category_.json", dir), []byte(categoryJson), 0644)
	}

	filePath := fmt.Sprintf("%s/%s.md", dir, pagename)

	return os.WriteFile(filePath, []byte(markdown), 0644)
}

// Process blocks recursively
func processBlocks(token string, blockID string, outputDir string) {
	var nextCursor string
	hasMore := true

	for hasMore {
		response, err := fetchChildren(token, blockID, nextCursor)
		if err != nil {
			log.Fatalf("Failed to fetch children: %v", err)
		}

		for index, block := range response.Results {

			switch block.Type {
			case "link_to_page":
				page, err := fetchPage(token, block.LinkToPage.PageID)
				if err != nil {
					log.Println("failed to fetch page id: ", block.LinkToPage.PageID)
					continue
				}
				log.Println("FETCHING PAGE", page.ID)
				writeMarkdown(outputDir, token, *page, index)
			case "child_page":
				if block.HasChildren {
					subOutput := outputDir + "/" + block.ChildPage.Title
					os.MkdirAll(subOutput, 0755)
					processBlocks(token, block.ID, subOutput)

				}
			}

			// if block.Type == "child_page" {
			// 	fmt.Printf("Processing child page: %s\n", block.ID)
			// 	// Write the markdown for the child page
			// 	_, err := pageToMarkdown(token, NotionPage{ID: block.ID}) // Simulate a page
			// 	if err == nil {
			// 		writeMarkdown(outputDir, token, NotionPage{ID: block.ID}) // Write to file
			// 	}
			// } else if block.Type == "child_database" {
			// 	fmt.Printf("Processing child database: %s\n", block.ID)
			// 	// Fetch and process databases for pages
			// 	processDatabases(token, block.ID, outputDir)
			// }
		}

		hasMore = response.HasMore
		nextCursor = response.NextCursor
		time.Sleep(1 * time.Second)
	}
}

// Process pages in a database
func processDatabases(token string, databaseID string, outputDir string) {
	var nextCursor string
	hasMore := true

	for hasMore {
		response, err := fetchPagesFromDatabase(token, databaseID, nextCursor)
		if err != nil {
			log.Fatalf("Failed to fetch pages from database: %v", err)
		}

		for index, page := range response.Results {
			fmt.Printf("Writing markdown for page: %s\n", page.ID)
			if err := writeMarkdown(outputDir, token, page, index); err != nil {
				log.Printf("Failed to write markdown for page %s: %v", page.ID, err)
			}
		}

		hasMore = response.HasMore
		nextCursor = response.NextCursor
		time.Sleep(1 * time.Second)
	}
}

func main() {
	token := flag.String("t", "", "Notion API token")
	rootID := flag.String("r", "", "Root block ID (page or database)")
	outputDir := flag.String("o", "./output", "Output directory for markdown files")
	DocsRoot := flag.String("docs", "/docs", "root docs directory")
	AssetsRoot := flag.String("assets", "./static", "root docs directory")

	flag.Parse()

	if *token == "" || *rootID == "" {
		log.Fatal("Notion API token and root ID are required")
	}

	if _, err := os.Stat(*outputDir); os.IsNotExist(err) {
		os.MkdirAll(*outputDir, os.ModePerm)
	}

	conf.DocsRoot = *DocsRoot
	conf.AssetsDir = *AssetsRoot
	conf.APIToken = *token

	processBlocks(*token, *rootID, *outputDir)

	fmt.Println("Export completed successfully.")
}
