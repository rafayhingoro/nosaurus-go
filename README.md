# 🦕 Nosaurus-Go: Notion to Docusaurus Converter

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

**Transform your Notion workspace into beautiful Docusaurus documentation with full pagination support and rich content preservation.**

Nosaurus-Go is a powerful command-line tool that seamlessly converts Notion pages, databases, and blocks into Markdown files optimized for Docusaurus. Perfect for documentation sites, knowledge bases, and content migration from Notion to static site generators.

## ✨ Key Features

- 🚀 **Full Notion API Integration** - Complete support for Notion's API with proper pagination handling
- 📄 **Rich Content Support** - Tables, images, code blocks, callouts, lists, and more
- 🔗 **Smart Link Resolution** - Automatically converts Notion internal links to proper markdown links
- 🖼️ **Image Management** - Downloads and organizes images with unique naming
- 📁 **Hierarchical Structure** - Maintains parent-child relationships and nested content
- ⚡ **Pagination Handling** - Fetches all content beyond Notion's 100-block API limit
- 🐛 **Debug Mode** - Optional detailed logging and API response saving for troubleshooting
- 🔧 **Environment Variables** - Secure token management with environment variable support

## 🎯 Use Cases

- **Documentation Sites** - Convert Notion workspaces to Docusaurus documentation
- **Knowledge Bases** - Migrate internal wikis from Notion to static sites
- **Blog Content** - Transform Notion pages into markdown blog posts
- **API Documentation** - Convert Notion-based API docs to developer-friendly formats
- **Team Wikis** - Export team knowledge from Notion to searchable documentation

## 📋 Prerequisites

- **Go 1.21+** installed on your machine
- **Notion API Token** - [Create an integration](https://www.notion.so/my-integrations)
- **Root Block/Page ID** - The starting point for your content export

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/rafayhingoro/nosaurus-go.git
cd nosaurus-go

# Build the application
go build .

# Make it executable (Linux/Mac)
chmod +x nosaurus-go
```

### Basic Usage

```bash
# Basic conversion
./nosaurus-go -t your_notion_token -r page_or_database_id

# With custom output directory
./nosaurus-go -t your_notion_token -r page_or_database_id -o ./my-docs

# With debug mode for troubleshooting
./nosaurus-go -t your_notion_token -r page_or_database_id --debug
```

### Environment Variables (Recommended)

```bash
# Set environment variables for security
export NOTION_API_TOKEN="your_notion_integration_token"
export NOTION_ROOT_ID="your_root_page_or_database_id"

# Run without exposing tokens in command history
./nosaurus-go
```

## 🛠️ Configuration Options

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `-t` | `NOTION_API_TOKEN` | - | Notion API integration token (required) |
| `-r` | `NOTION_ROOT_ID` | - | Root page or database ID (required) |
| `-o` | - | `./output` | Output directory for markdown files |
| `-docs` | - | `/docs` | Root docs directory path for internal links |
| `-assets` | - | `./static` | Assets directory for images and files |
| `--debug` | - | `false` | Enable debug mode for detailed logging |

## 📖 Detailed Usage Examples

### Convert a Notion Page

```bash
# Convert a single page and its subpages
./nosaurus-go -t secret_abc123 -r 1a2b3c4d-5e6f-7890-abcd-ef1234567890
```

### Convert a Notion Database

```bash
# Convert all pages in a database
./nosaurus-go -t secret_abc123 -r database-id-here -o ./docs
```

### Production Setup with Environment Variables

```bash
# 1. Create a .env file (don't commit this!)
echo "NOTION_API_TOKEN=secret_abc123" > .env
echo "NOTION_ROOT_ID=1a2b3c4d-5e6f-7890-abcd-ef1234567890" >> .env

# 2. Load environment variables
source .env

# 3. Run the converter
./nosaurus-go -o ./docusaurus/docs --debug
```

## 🔧 Notion API Setup

### 1. Create a Notion Integration

1. Go to [Notion Integrations](https://www.notion.so/my-integrations)
2. Click **"+ New integration"**
3. Name your integration (e.g., "Docusaurus Converter")
4. Select the workspace
5. Copy the **Internal Integration Token**

### 2. Share Your Content with the Integration

1. Open your Notion page/database
2. Click **"Share"** in the top right
3. Click **"Invite"**
4. Search for your integration name
5. Select it and grant access

### 3. Get the Page/Database ID

```
From URL: https://notion.so/My-Page-1a2b3c4d5e6f7890abcdef1234567890
Page ID:  1a2b3c4d-5e6f-7890-abcd-ef1234567890
```

## 📁 Output Structure

Nosaurus-Go creates a well-organized directory structure:

```
output/
├── page-title.md                 # Individual pages
├── database-pages/               # Database contents
│   ├── index.md                 # Database overview
│   ├── _category_.json          # Docusaurus category config
│   └── page-1.md               # Database entries
├── nested-content/              # Hierarchical content
│   ├── parent-page/
│   │   ├── index.md
│   │   └── child-page.md
└── static/
    └── docs-images/             # Downloaded images
        ├── image1.png
        └── image2.jpg
```

## 🐛 Troubleshooting

### Debug Mode

Enable debug mode to see detailed processing information:

```bash
./nosaurus-go -t your_token -r your_id --debug
```

This creates a `debug_responses/` folder with:
- Raw API responses as JSON files
- Detailed pagination information
- Processing logs for each step

### Common Issues

**"Failed to fetch page content"**
```bash
# Check your integration has access to the page
# Verify the page ID is correct
# Ensure the integration token is valid
```

**"Rate limited"**
```bash
# The tool automatically handles rate limiting
# Large exports may take time due to API limits
# Use --debug to monitor progress
```

**"Missing content after certain sections"**
```bash
# This was a pagination issue - now fixed!
# The tool automatically fetches all content beyond 100 blocks
# Use --debug to verify all pages are fetched
```

## 🔄 Migration Workflow

### For Existing Docusaurus Sites

```bash
# 1. Backup your existing docs
cp -r docs docs.backup

# 2. Run the converter
./nosaurus-go -o ./docs-notion

# 3. Review and merge content
# 4. Update docusaurus.config.js if needed
# 5. Test your site locally
npm start
```

### Content Mapping

| Notion Element | Markdown Output |
|---------------|----------------|
| Headings (H1-H3) | `# ## ###` |
| Paragraphs | Standard markdown |
| **Bold**, *Italic* | `**bold**`, `*italic*` |
| `Code` | Inline code |
| Code blocks | Fenced code blocks |
| Bulleted lists | `- item` |
| Numbered lists | `1. item` |
| Tables | HTML tables (Docusaurus compatible) |
| Images | Downloaded and linked |
| Internal links | Converted to relative paths |
| Callouts | Blockquotes with emoji |

## 🤝 Contributing

We welcome contributions! Here's how to get started:

```bash
# 1. Fork the repository
# 2. Clone your fork
git clone https://github.com/yourusername/nosaurus-go.git

# 3. Create a feature branch
git checkout -b feature/your-feature-name

# 4. Make your changes and test
go test ./...

# 5. Submit a pull request
```

### Development Setup

```bash
# Install dependencies
go mod tidy

# Run tests
go test -v ./...

# Build for development
go build -o nosaurus-go-dev .
```

## 📊 Performance

- **Small sites** (< 50 pages): ~1-2 minutes
- **Medium sites** (50-200 pages): ~5-10 minutes  
- **Large sites** (200+ pages): ~15-30 minutes

*Times vary based on content complexity and API rate limits*

## 🔐 Security Best Practices

1. **Never commit API tokens** to version control
2. **Use environment variables** for tokens
3. **Rotate tokens regularly** in Notion settings
4. **Limit integration permissions** to necessary workspaces

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with ❤️ using Go
- Inspired by the need to bridge Notion and Docusaurus
- Thanks to the Notion API team for comprehensive documentation

---

**Made with 🦕 by [Rafay Hingoro](https://github.com/rafayhingoro)**

*Star ⭐ this repo if it helped you migrate from Notion to Docusaurus!*
