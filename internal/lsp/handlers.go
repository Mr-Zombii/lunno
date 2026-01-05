package lsp

import (
	"encoding/json"
	"lunno/internal/lexer"
	"lunno/internal/parser"
	"strings"
)

func handleInitialize(req RequestMessage) {
	result := InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync: 1,
			CompletionProvider: &CompletionOptions{
				ResolveProvider:   false,
				TriggerCharacters: []string{"."},
			},
		},
	}
	send(ResponseMessage{
		Jsonrpc: "2.0",
		ID:      req.ID,
		Result:  result,
	})
}

func handleDidOpen(params json.RawMessage) {
	var p DidOpenTextDocumentParams
	err := json.Unmarshal(params, &p)
	if err != nil {
		return
	}
	docs[p.TextDocument.URI] = p.TextDocument.Text
	updateSymbols(p.TextDocument.URI, p.TextDocument.Text, lx)
	runDiagnostics(p.TextDocument.URI, lx)
}

func handleDidChange(params json.RawMessage) {
	var p DidChangeTextDocumentParams
	err := json.Unmarshal(params, &p)
	if err != nil {
		return
	}
	if len(p.ContentChanges) > 0 {
		docs[p.TextDocument.URI] = p.ContentChanges[0].Text
		updateSymbols(p.TextDocument.URI, p.ContentChanges[0].Text, lx)
		runDiagnostics(p.TextDocument.URI, lx)
	}
}

func handleShutdown(req RequestMessage) {
	send(ResponseMessage{
		Jsonrpc: "2.0",
		ID:      req.ID,
		Result:  nil,
	})
}

func runDiagnostics(uri string, lx *lexer.Lexer) {
	text := docs[uri]
	_, tokens, err := lexer.Tokenize(text, uri)
	if err != nil {
		send(map[string]any{
			"jsonrpc": "2.0",
			"method":  "textDocument/publishDiagnostics",
			"params": PublishDiagnosticsParams{
				URI: uri,
				Diagnostics: []Diagnostic{{
					Severity: 1,
					Message:  "Lexing error: " + err.Error(),
					Range: Range{
						Start: Position{Line: 0, Character: 0},
						End:   Position{Line: 0, Character: 1},
					},
				}},
			},
		})
		return
	}

	_, parseErrs := parser.ParseProgram(tokens, lx)

	var diagnostics []Diagnostic
	for _, e := range parseErrs {
		diagnostics = append(diagnostics, Diagnostic{
			Severity: 1,
			Message:  e,
			Range: Range{
				Start: Position{Line: 0, Character: 0},
				End:   Position{Line: 0, Character: 1},
			},
		})
	}

	send(map[string]any{
		"jsonrpc": "2.0",
		"method":  "textDocument/publishDiagnostics",
		"params": PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: diagnostics,
		},
	})
}

func handleCompletion(req RequestMessage) {
	var params CompletionParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return
	}
	text, ok := docs[params.TextDocument.URI]
	if !ok {
		return
	}
	lines := strings.Split(text, "\n")
	if params.Position.Line >= len(lines) {
		return
	}
	line := lines[params.Position.Line]
	charIndex := params.Position.Character
	if charIndex > len(line) {
		charIndex = len(line)
	}
	prefix := ""
	if charIndex > 0 {
		fields := strings.Fields(line[:charIndex])
		if len(fields) > 0 {
			prefix = fields[len(fields)-1]
		}
	}
	var suggestions []CompletionItem
	for kw := range lexer.Keywords() {
		if strings.HasPrefix(kw, prefix) {
			suggestions = append(suggestions, CompletionItem{
				Label: kw,
				Kind:  14,
			})
		}
	}
	for _, sym := range symbols[params.TextDocument.URI] {
		if strings.HasPrefix(sym.Name, prefix) {
			kind := 6
			if sym.Kind == FunctionSymbol {
				kind = 3
			}
			suggestions = append(suggestions, CompletionItem{
				Label: sym.Name,
				Kind:  kind,
			})
		}
	}
	send(ResponseMessage{
		Jsonrpc: "2.0",
		ID:      req.ID,
		Result:  suggestions,
	})
}
