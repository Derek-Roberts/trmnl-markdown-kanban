package kanban

import (
    "io/ioutil"
    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/ast"
    "github.com/yuin/goldmark/text"
)

// Card represents a single Kanban card.
type Card struct {
    Title   string
    Content string
}

// Column represents one column in the Kanban board.
type Column struct {
    Title string
    Cards []Card
}

// Board holds all columns parsed from Markdown.
type Board struct {
    Columns []Column
}

// LoadBoard reads the Markdown at 'path', parses it, and builds a Board.
func LoadBoard(path string) (*Board, error) {
    // 1. Read file
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    // 2. Parse into AST
    md := goldmark.New()
    reader := text.NewReader(data)
    doc := md.Parser().Parse(reader)

    var board Board
    var currentCol *Column

    // 3. Walk the AST
    ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
        if !entering {
            return ast.WalkContinue, nil
        }
        switch n.Kind() {
        case ast.KindHeading:
            heading := n.(*ast.Heading)
            // Only treat level-2 headings as columns (adjust level as needed)
            if heading.Level == 2 {
                // Extract the heading text
                text := string(heading.Text(data))
                board.Columns = append(board.Columns, Column{Title: text})
                currentCol = &board.Columns[len(board.Columns)-1]
            }

        case ast.KindListItem:
            if currentCol != nil {
                // First child of a list item is typically a paragraph
                if first := n.FirstChild(); first != nil && first.Kind() == ast.KindParagraph {
                    title := string(first.Text(data))
                    // Collect any following siblings as the card’s content
                    var content string
                    for sib := first.NextSibling(); sib != nil; sib = sib.NextSibling() {
                        content += string(sib.Text(data)) + "\n"
                    }
                    currentCol.Cards = append(currentCol.Cards, Card{
                        Title:   title,
                        Content: content,
                    })
                }
            }
        }
        return ast.WalkContinue, nil
    })

    return &board, nil
}
