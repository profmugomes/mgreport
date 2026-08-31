# MGReport

[![License](https://img.shields.io/badge/license-PolyForm%20Perimeter%201.0.1-5351FB)](LICENSE.md)

> **Aviso importante:** Este componente está em desenvolvimento!

MGReport é uma biblioteca de relatórios para **Fyne**. Ela permite criar relatórios visuais com seções (bandas), vinculação a dados, zoom, impressão e exportação PDF.

---

## 🧩 Compatibilidade

* Go 1.26.5+
* Fyne 2.8.0

---

## 📦 Instalação

```bash
go get github.com/profmugomes/mgreport
```

---

## 🚀 Início Rápido

```go
package main

import (
    "fyne.io/fyne/v2/app"
    "github.com/profmugomes/mgreport"
)

func main() {
    report := mgreport.NewReport()
    report.SetTitle("Meu Relatório")

    // Cabeçalho do relatório (aparece 1 vez)
    report.Section(mgreport.ReportHeader, 60).
        AddLabel("Empresa XYZ", mgreport.Position{X: 10, Y: 10},
            mgreport.WithFontSize(18), mgreport.WithBold(true)).
        AddLine(10, 50, 585, 50)

    // Cabeçalho da página (repete em cada página)
    report.Section(mgreport.PageHeader, 25).
        AddLabel("Nome", mgreport.Position{X: 10, Y: 8}, mgreport.WithBold(true)).
        AddLabel("Valor", mgreport.Position{X: 200, Y: 8}, mgreport.WithBold(true))

    // Detalhe (repete para cada registro)
    report.Section(mgreport.Detail, 18).
        AddField("nome", mgreport.Position{X: 10, Y: 4}).
        AddField("valor", mgreport.Position{X: 200, Y: 4},
            mgreport.WithFieldFormat(mgreport.FormatCurrency))

    // Rodapé da página
    report.Section(mgreport.PageFooter, 25).
        AddLabel("Página", mgreport.Position{X: 250, Y: 8}).
        AddPageNumber(mgreport.Position{X: 290, Y: 8})

    // Dados
    report.SetDataSource([]map[string]interface{}{
        {"nome": "João", "valor": 1500.00},
        {"nome": "Maria", "valor": 2300.00},
    })

    // Visualizar
    app := app.New()
    mgreport.ShowPreview(app, report)
    app.Run()
}
```

---

## 📐 Seções do Relatório

| Seção | Descrição |
|-------|-----------|
| `ReportHeader` | Cabeçalho do relatório. Aparece apenas **uma vez** no início. |
| `PageHeader` | Cabeçalho da página. **Repete** no topo de cada página. |
| `Detail` | Detalhe/Linha. **Repete** para cada registro no DataSource. |
| `ReportFooter` | Rodapé do relatório. Aparece apenas **uma vez** no final. |
| `PageFooter` | Rodapé da página. **Repete** na base de cada página. |

---

## 🎨 Itens Disponíveis

### Label (Texto Estático)
```go
section.AddLabel("Texto", mgreport.Position{X: 10, Y: 10},
    mgreport.WithFontSize(14),
    mgreport.WithFontFace("Arial"),
    mgreport.WithColor(mgreport.ColorBlack),
    mgreport.WithBold(true),
    mgreport.WithAlign(mgreport.AlignCenter),
)
```

### Field (Campo de Dados)
```go
section.AddField("nome_do_campo", mgreport.Position{X: 10, Y: 10},
    mgreport.WithFieldFontSize(12),
    mgreport.WithFieldFormat(mgreport.FormatCurrency),
    mgreport.WithFieldColor(mgreport.ColorBlue),
    mgreport.WithFieldAlign(mgreport.AlignRight),
    mgreport.WithFieldBold(true),
    mgreport.WithFieldFontFace("Arial"),
)
```

### Line (Linha)
```go
section.AddLine(10, 50, 585, 50,
    mgreport.WithLineColor(mgreport.ColorGray),
    mgreport.WithLineWidth(0.5),
)
```

### PageNumber (Número de Página)
```go
section.AddPageNumber(mgreport.Position{X: 290, Y: 8},
    mgreport.WithFontSize(10),
    mgreport.WithColor(mgreport.ColorGray),
)
```

---

## 📊 Formatadores de Dados

| Formatador | Descrição | Exemplo de Saída |
|-----------|-----------|-----------------|
| `FormatCurrency` | Moeda brasileira (R$) | `R$ 1500.00` |
| `FormatInteger` | Número inteiro | `1500` |
| `FormatPercent` | Porcentagem | `15.00%` |
| `FormatDateBR` | Data brasileira (DD/MM/AAAA) | `25/12/2026` |

Criar formatador personalizado:
```go
func FormatCustom(v interface{}) string {
    switch val := v.(type) {
    case float64:
        return fmt.Sprintf("%.1f kg", val)
    default:
        return fmt.Sprintf("%v", v)
    }
}

section.AddField("peso", pos, mgreport.WithFieldFormat(FormatCustom))
```

---

## 🔧 Tamanhos de Papel

| Constante | Dimensões (pontos) | Dimensões (mm) |
|-----------|-------------------|----------------|
| `mgreport.A4` | 595 × 842 | 210 × 297 |
| `mgreport.Letter` | 612 × 792 | 216 × 279 |
| `mgreport.Legal` | 612 × 1008 | 216 × 356 |
| `mgreport.A5` | 420 × 595 | 148 × 210 |

---

## 🖥️ Visualizador (Viewer)

O `ShowPreview` abre uma janela completa com:

- **Navegação**: Primeira, Anterior, Próxima, Última página
- **Zoom**: +/- 10%, "Ajustar à Janela"
- **Exportação**: Botão para salvar PDF
- **Impressão**: Abre o PDF no visualizador do sistema
- **Atalhos de Teclado**:
  - `←` / `→` — Navegar entre páginas
  - `+` / `-` — Zoom in/out
  - `Home` / `End` — Primeira/Última página

---

## 📁 Estrutura do Projeto

```
mgreport/
├── types.go          # Tipos fundamentais (PaperSize, Position, Margins, etc.)
├── report.go         # Report + Section + SectionBuilder (API Builder)
├── items.go          # Label, Field, Line, PageNumber + Formatadores
├── renderer.go       # Motor de renderização Report → Pages
├── viewer.go         # ReportViewer widget + ShowPreview
├── export_pdf.go     # Exportação PDF via fpdf
├── print.go          # Impressão via sistema
├── mgreport.go       # Documentação do pacote
├── examples/
│   └── demo/
│       └── main.go   # Projeto demo completo
├── go.mod
├── go.sum
├── README.md
└── LICENSE.md
```

---

## 👤 Autor

**Murilo Gomes Julio**

🔗 [https://www.profmugomes.com.br](https://www.profmugomes.com.br)

📺 [https://youtube.com/@profmugomes](https://youtube.com/@profmugomes)

---

## License

Copyright (c) 2026 Murilo Gomes Julio. All Rights Reserved.

This project is licensed under the PolyForm Perimeter License 1.0.1.

### Summary

This software is available for commercial and noncommercial use, subject to the terms of the PolyForm Perimeter License 1.0.1.

You may:

* ✔ Use the software for commercial and noncommercial purposes.
* ✔ Inspect and study the source code.
* ✔ Modify the software.
* ✔ Create derivative works based on the software.
* ✔ Redistribute the software and permitted modifications.

You may not:

* ✖ Provide a product that competes with the software.

See the full license terms at LICENSE.md.

This summary is provided for convenience only and does not replace or modify the full license terms.
