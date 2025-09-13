# BrxAgente-desafio4

## Sobre

Este aplicativo automatiza o processo mensal de compra de VR (Vale Refeição), garantindo que cada colaborador receba o valor correto, considerando ausências, férias, datas de admissão ou desligamento e calendário de feriados.

O aplicativo processa planilhas Excel contendo dados de colaboradores e gera uma planilha de resultado com os valores de VR calculados conforme as regras de negócio definidas pelos sindicatos.

## Índice

- [Requisitos](#requisitos)
- [Instalação](#instalação)
- [Como Usar](#como-usar)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Desenvolvimento](#desenvolvimento)
- [Contribuição](#contribuição)
- [Licença](#licença)
- [Troubleshooting](#troubleshooting)

## Requisitos

Antes de começar, certifique-se de ter instalado em seu sistema:

- [Go](https://golang.org/dl/) versão 1.21 ou superior
- [Node.js](https://nodejs.org/) versão 16 ou superior
- [Wails](https://wails.io/docs/gettingstarted/installation) v2
- [Git](https://git-scm.com/)

### Requisitos Opcionais (para funcionalidades de IA)

- Chave de API do OpenAI (para usar modelos GPT)
- [Ollama](https://ollama.ai/) (para usar modelos locais)

## 🤖 Inteligência Artificial

### Chat Avançado
O BrxAgente possui um sistema de chat inteligente que combina:
- **Base de conhecimento estruturada** sobre políticas de VR
- **Dados processados em tempo real** do último processamento (1794+ colaboradores)
- **Classificação automática** de perguntas
- **Duas modalidades de análise**: Rápida (⚡) ou Detalhada (🐌)

### 🎯 Modalidades de Análise

#### **🐌 Análise Detalhada (Padrão)**
- ✅ **Cache desabilitado** por padrão para máxima precisão
- ✅ **Dados completos**: Envia todos os 1794+ colaboradores individuais (216KB+)
- ✅ **Respostas brutas da LLM**: Sem formatação pré-definida, análises personalizadas
- ⏱️ **Tempo**: ~13 segundos para análises profundas
- 🎯 **Ideal para**: Análises estatísticas, detecção de anomalias, relatórios customizados

#### **⚡ Análise Rápida (Opcional)**
- 🚀 **Cache habilitado** para velocidade máxima
- 📊 **Dados resumidos**: Informações agregadas otimizadas (702 chars)
- 🎨 **Respostas formatadas**: Interface padronizada e consistente
- ⏱️ **Tempo**: ~2 segundos para consultas diretas
- 🎯 **Ideal para**: Consultas simples, políticas, validações rápidas

### Exemplos de Análises Avançadas
- **Detecção de Anomalias:** "Identifique colaboradores com VR 20% acima da média do sindicato"
- **Análise Comparativa:** "Compare eficiência de VR entre sindicatos considerando variabilidade"
- **Validação de Conformidade:** "Analise colaboradores com valores incorretos e sugira correções"

### Capacidades do Sistema
- 📋 **Consultor de Políticas**: Respostas baseadas em regulamentações oficiais
- 🧮 **Explicação de Cálculos**: Como aplicar regras para cenários específicos
- 📊 **Análise Estatística Profunda**: Processamento de dados completos com LLM
- 🔍 **Detecção de Padrões**: Identificação automática de anomalias e inconsistências
- 🎯 **Confiança das Respostas**: Indica nível de certeza (Baixa/Média/Alta/Muito Alta)

## Instalação

### 1. Clonar o Repositório

```bash
git clone https://github.com/vstram/BrxAgente-desafio4.git
cd BrxAgente-desafio4
```

### 2. Instalar Dependências

```bash
# Instalar dependências do Go
go mod tidy

# Instalar dependências do frontend
cd frontend
npm install
cd ..
```

### 3. Construir o Aplicativo

```bash
# Construir para produção
wails build

# Ou executar em modo de desenvolvimento
wails dev
```

O aplicativo será construído na pasta `build/bin/`.

## Como Usar

### 1. Preparar as Planilhas

Certifique-se de ter as seguintes planilhas no formato `.xlsx` em uma pasta:

1. `ATIVOS.xlsx` - Lista de colaboradores ativos
2. `FÉRIAS.xlsx` - Períodos de férias dos colaboradores
3. `DESLIGADOS.xlsx` - Colaboradores desligados
4. `Base sindicato x valor.xlsx` - Valores de VR por sindicato
5. `Base dias uteis.xlsx` - Dias úteis por sindicato

### 2. Executar o Aplicativo

- Se construiu o aplicativo: Execute o binário gerado em `build/bin/`
- Se estiver em modo de desenvolvimento: Execute `wails dev`

### 3. Usar a Interface

1. Clique no botão "Selecionar Diretório" e escolha a pasta com as planilhas
2. Clique no botão "Fazer Análise" para iniciar o processamento
3. A planilha de resultado será gerada na pasta de Downloads do usuário

### 4. Usar Chat Avançado

**🐌 Modo Padrão (Análise Detalhada):**
- O chat abre automaticamente em modo detalhado para análises profundas
- Faça perguntas complexas sobre os dados processados
- Sistema processa todos os colaboradores individuais para análises precisas

**⚡ Modo Rápido (Opcional):**
- Clique no botão do cache (🐌→⚡) para ativar respostas rápidas
- Ideal para consultas simples e validações rápidas
- Usa dados agregados para velocidade máxima

**Exemplos de perguntas por modalidade:**
- **Detalhada:** "Compare a distribuição de VR entre sindicatos e identifique outliers"
- **Rápida:** "Quantos colaboradores foram processados neste mês?"

### 5. Configuração (Opcional)

Para usar funcionalidades de IA:
1. Clique no botão "Configurações"
2. Insira sua chave de API do OpenAI ou configure o acesso ao Ollama
3. Salve as configurações

## Estrutura do Projeto

```
BrxAgente-desafio4/
├── app.go                 # Ponto de entrada da aplicação
├── main.go                # Configuração do Wails
├── wails.json             # Configuração do Wails
├── go.mod                 # Dependências do Go
├── go.sum                 # Soma de verificação das dependências
├── frontend/              # Código do frontend React
│   ├── src/
│   │   ├── App.tsx        # Componente principal
│   │   └── ...
│   └── ...
├── internal/              # Pacotes internos
│   ├── calculo/           # Lógica de cálculo de VR
│   ├── excel/             # Manipulação de planilhas Excel
│   ├── modelo/            # Estruturas de dados
│   ├── validacao/         # Validação de dados
│   └── config/           # Configuração da aplicação
├── files/                 # Arquivos de exemplo e dados
└── build/                 # Arquivos de construção
```

### Pacotes Internos

- **`internal/calculo`**: Contém a lógica de cálculo do VR, incluindo regras de negócio como datas quebradas, regra de desligamento, cálculo de dias úteis por sindicato, etc.

- **`internal/excel`**: Responsável pela leitura e escrita de planilhas Excel, incluindo a consolidação das bases de dados.

- **`internal/modelo`**: Define as estruturas de dados utilizadas no projeto, como `Colaborador` e `Periodo`.

- **`internal/validacao`**: Implementa as regras de validação dos dados, incluindo a verificação de integridade das planilhas.

- **`internal/config`**: Gerencia as configurações da aplicação, incluindo chaves de API e configurações do Ollama.

- **`internal/feriados`**: Gerencia os feriados nacionais, estaduais e municipais para cálculo de dias úteis.

- **`internal/chat`**: Implementa a funcionalidade de chat com IA.

## Desenvolvimento

### Ambiente de Desenvolvimento

Para configurar o ambiente de desenvolvimento:

1. Instale todas as dependências conforme a seção [Instalação](#instalação)
2. Execute `wails dev` para iniciar o modo de desenvolvimento com hot-reload

### Estrutura de Código

O projeto segue uma arquitetura modular com os seguintes princípios:

1. **Separação de Responsabilidades**: Cada pacote tem uma responsabilidade específica
2. **Testabilidade**: Código escrito com testes em mente
3. **Manutenibilidade**: Código limpo e bem documentado

### Executando Testes

```bash
# Executar todos os testes
go test ./...

# Executar testes com cobertura
go test -cover ./...
```

### Construção

```bash
# Construir para produção
wails build

# Construir para uma plataforma específica
wails build -platform linux/amd64
```

## Contribuição

Contribuições são bem-vindas! Para contribuir:

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/NovaFeature`)
3. Commit suas mudanças (`git commit -m 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/NovaFeature`)
5. Abra um Pull Request

### Diretrizes de Contribuição

- Siga o estilo de código existente
- Adicione testes para novas funcionalidades
- Atualize a documentação conforme necessário
- Mantenha mensagens de commit claras e descritivas

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## Troubleshooting

### Problemas Comuns

#### 1. Erro ao construir o aplicativo

**Problema**: `wails build` falha com erros de dependência

**Solução**:
```bash
go mod tidy
cd frontend && npm install && cd ..
wails build
```

#### 2. Erro ao executar em modo de desenvolvimento

**Problema**: `wails dev` falha com erros do frontend

**Solução**:
```bash
cd frontend
npm install
npm run build
cd ..
wails dev
```

#### 3. Planilha de resultado não é gerada

**Problema**: O aplicativo processa as planilhas mas não gera o arquivo de resultado

**Solução**:
- Verifique se há permissões para escrever na pasta de Downloads
- Certifique-se de que as planilhas de entrada estão no formato correto
- Verifique o console do aplicativo para mensagens de erro

#### 4. Erro ao carregar planilhas

**Problema**: Mensagem de erro "Erro ao abrir o arquivo Excel"

**Solução**:
- Verifique se os arquivos estão no formato `.xlsx`
- Certifique-se de que os arquivos não estão corrompidos
- Verifique se os nomes dos arquivos correspondem exatamente aos esperados

#### 5. Problemas com configuração de IA

**Problema**: Erros ao usar funcionalidades de IA

**Solução**:
- Verifique se a chave de API do OpenAI está correta
- Se estiver usando Ollama, certifique-se de que o serviço está em execução
- Verifique as configurações de rede e firewall

### Suporte

Se você encontrar problemas que não conseguem ser resolvidos com estas soluções, por favor:

1. Verifique se há issues similares no repositório
2. Abra uma nova issue descrevendo o problema em detalhes
3. Inclua informações sobre seu ambiente (sistema operacional, versões, etc.)
4. Inclua logs de erro relevantes