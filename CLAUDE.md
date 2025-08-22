# Claude Code - Recomendações para o Projeto BrxAgente-desafio4

## Visão Geral do Projeto

O projeto BrxAgente-desafio4 é um aplicativo desktop desenvolvido com Wails (Go no backend e React no frontend) que visa automatizar o processo mensal de compra de Vale Refeição (VR) para colaboradores, considerando diversos fatores como ausências, férias, admissões, desligamentos e calendário de feriados.

## Estrutura Atual do Projeto

```
BrxAgente-desafio4/
├── app.go              # Lógica principal do backend em Go (arquivo temporário para testes)
├── main.go             # Ponto de entrada da aplicação (versão de teste para desenvolvimento)
├── go.mod/go.sum       # Gerenciamento de dependências Go
├── wails.json          # Configuração do Wails
├── frontend/           # Código frontend React
├── files/              # Planilhas de dados Excel
├── build/              # Arquivos de build
└── internal/           # Implementação da lógica de negócio e processamento de dados
    ├── calculo/        # Funções de cálculo e consolidação de dados
    ├── excel/          # Funções de leitura e escrita de planilhas Excel
    ├── modelo/         # Estruturas de dados (structs) para representar colaboradores
    └── validacao/      # Funções de validação de dados
```

> ⚠️ **Nota**: O arquivo `app.go` foi criado temporariamente para testes e não contém a lógica principal da aplicação. A lógica de negócio foi implementada nos pacotes dentro do diretório `internal/`, seguindo as melhores práticas de organização de código Go.

## Recomendações Técnicas

### 1. Implementação da Lógica de Negócio

A lógica de cálculo do VR já foi parcialmente implementada nos pacotes internos:

```go
// A lógica principal foi implementada em internal/calculo/consolidar.go:
func ConsolidarBases(diretorioPlanilhas string) (map[string]*modelo.Colaborador, error) {
    // Implementar a lógica de processamento das planilhas
    // - Leitura dos arquivos Excel
    // - Consolidação das bases de dados
    // - Aplicação das regras de negócio
    // - Geração da planilha final
    return colaboradores, nil
}
```

> ⚠️ **Nota**: A lógica de negócio foi implementada seguindo as melhores práticas de organização de código Go, utilizando pacotes internos (`internal/`) em vez de colocar toda a lógica no arquivo `app.go`.

### 2. Manipulação de Arquivos Excel

Adicione dependências para manipulação de Excel:

```bash
go get github.com/xuri/excelize/v2
```

E implemente a leitura/escrita de planilhas:

```go
import "github.com/xuri/excelize/v2"

func (a *App) LerPlanilha(caminho string) (*excelize.File, error) {
    return excelize.OpenFile(caminho)
}
```

### 3. Interface do Usuário

Atualize o frontend para incluir:
- Seleção de diretório com as planilhas
- Botão para iniciar o processamento
- Área de visualização dos resultados
- Opções de configuração (chaves API, etc.)

### 4. Estrutura de Dados

Crie structs para representar os dados dos colaboradores:

```go
type Colaborador struct {
    Matricula       string
    Nome            string
    DataAdmissao    time.Time
    DataDesligamento *time.Time
    Sindicato       string
    Afastamentos    []Periodo
    Ferias          []Periodo
    // Outros campos relevantes
}

type Periodo struct {
    Inicio time.Time
    Fim    time.Time
}
```

## Recomendações de Funcionalidades

### 1. Processamento Automatizado

Implementar o cálculo automático com base nas regras especificadas:

- Cálculo de dias úteis por sindicato
- Exclusão de perfis não elegíveis (diretores, estagiários, aprendizes, etc.)
- Tratamento de datas quebradas (admissões/desligamentos no meio do mês)
- Aplicação das regras de desligamento (comunicado até dia 15)
- Consideração de feriados estaduais/municipais

### 2. Validações

Implementar as validações mencionadas:
- Consistência de datas
- Campos obrigatórios
- Formatação correta dos dados
- Relacionamento entre as planilhas (via matrícula)

### 3. Geração de Relatório

Gerar a planilha final conforme o modelo especificado:
- Valor total de VR por colaborador
- Rateio 80%/20% (empresa/colaborador)
- Formato compatível com o esperado pela operadora

## Melhorias na Arquitetura

### 1. Separação de Responsabilidades

Organize o código em pacotes distintos:
```
internal/
├── excel/       # Manipulação de planilhas
├── calculo/     # Lógica de cálculo do VR
├── modelo/      # Estruturas de dados
└── validacao/   # Regras de validação
```

> ⚠️ **Nota de Confidencialidade Importante**: Os nomes dos colaboradores não estão disponíveis nas planilhas por motivos de sigilo. Todas as referências aos colaboradores serão feitas exclusivamente através do campo **MATRICULA** como identificador único.

### 2. Tratamento de Erros

Implemente um sistema robusto de tratamento de erros:
```go
type ErroProcessamento struct {
    Mensagem string
    Arquivo  string
    Linha    int
}

func (e *ErroProcessamento) Error() string {
    return fmt.Sprintf("Erro no arquivo %s linha %d: %s", e.Arquivo, e.Linha, e.Mensagem)
}
```

### 3. Configuração

Adicione suporte a arquivos de configuração para:
- Caminhos padrão
- Regras de negócio parametrizáveis
- Credenciais de API

## Considerações Finais

1. **Testes**: Implemente testes unitários para as funções críticas de cálculo
2. **Documentação**: Mantenha o README.md atualizado com instruções de uso
3. **Performance**: Otimize o processamento para grandes volumes de dados
4. **Segurança**: Proteja chaves de API e dados sensíveis dos colaboradores
5. **Internacionalização**: Mantenha o suporte ao idioma Português Brasileiro

Essas recomendações fornecem uma base sólida para desenvolver um aplicativo completo que atenda às necessidades especificadas no desafio.