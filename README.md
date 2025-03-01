# nosaurus-go
Render from Notion to Docusaurus (converts files to markdown)

### Description
nosaurus-go is a tool designed to convert content from Notion into markdown files that can be used with Docusaurus. This makes it easy to integrate Notion content into a Docusaurus-based documentation site.

### Features
- Convert Notion pages to markdown
- Easily integrate with Docusaurus
- Simple command-line interface

# Installation
To install nosaurus-go, ensure that you have Go installed on your machine, then run:

```bash
go get github.com/rafayhingoro/nosaurus-go
```

### Usage
#### Prerequisites
- Notion API Token: You will need a Notion API token which you can generate from the Notion Integrations page.
- Root Block ID: The ID of the root block (page or database) from which you want to start the conversion.

#### Running nosaurus-go
Use the following command to convert the Notion content to markdown:

```bash
nosaurus-go -t <Notion_API_Token> -r <Root_Block_ID> -o <Output_Directory>
```
Replace <Notion_API_Token> with your Notion API token, <Root_Block_ID> with the ID of the root block, and <Output_Directory> with the desired output directory for the markdown files.

#### Command-Line Options
-t: Notion API token (required)
-r: Root block ID (page or database) (required)
-o: Output directory for markdown files (default: ./output)
-docs: Root docs directory (default: /docs)
-assets: Root assets directory (default: ./static)

### Configuration
nosaurus-go supports various configuration options to customize the conversion process. You can specify the root docs directory and root assets directory using the command-line options -docs and -assets.

### Contributing
Contributions are welcome! Please fork the repository and submit pull requests.

### License
This project is licensed under the MIT License.

Feel free to update or expand this README as the project evolves and more features or configuration options are added.
