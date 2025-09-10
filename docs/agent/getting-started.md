# 🚀 Guia de Início Rápido - Agente de IA

Este guia te levará de zero a produtivo com o Agente de IA BrxAgente em menos de 15 minutos.

## ⚡ Instalação Rápida

### 1. **Pré-requisitos**
```bash
# Verifique se tem Go instalado
go version  # Deve ser >= 1.19

# Verifique se tem Node.js instalado (para frontend)
node --version  # Deve ser >= 16

# Clone o projeto (se ainda não fez)
git clone https://github.com/vstram/BrxAgente-desafio4.git
cd BrxAgente-desafio4
```

### 2. **Build da Aplicação Desktop**
```bash
# Instale o Wails CLI se não tiver
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Compile a aplicação desktop
wails build

# Ou para desenvolvimento com hot-reload
wails dev
```

### 3. **Primeira Execução**
```bash
# Inicie a aplicação em modo desenvolvimento
wails dev

# Ou execute o build final
./build/bin/BrxAgente-desafio4  # Linux/Mac
./build/bin/BrxAgente-desafio4.exe  # Windows
```

A aplicação abrirá uma janela desktop com interface gráfica moderna.

## 🎯 Primeiro Uso - Interface Gráfica

### **Tela Principal**
Ao abrir a aplicação, você verá:
- 🏷️ **Cabeçalho**: Logo e título da aplicação
- 📁 **Seção de Seleção**: Botão para escolher pasta com planilhas Excel
- ▶️ **Seção de Processamento**: Botão para iniciar processamento (aparece após seleção válida)
- 📊 **Área de Resultados**: Exibe progresso e resultados do processamento
- ⚙️ **Configurações**: Modal acessível via botão no rodapé
- 💬 **Chat Inteligente**: Painel expansível na parte inferior da tela

### **Cenário 1: Configuração Inicial**
*"Primeira vez usando o sistema"*

1. **Clique em "Configurações"** na interface
2. **Configure sua API Key**:
   - OpenAI: Cole sua chave da API
   - Ou Ollama: Configure URL local
3. **Selecione Diretório das Planilhas**:
   - Clique no botão "Selecionar Pasta"
   - Escolha a pasta com seus arquivos .xlsx
4. **Teste a Conexão**: Clique em "Testar API"

### **Cenário 2: Processamento de VR**
*"Preciso processar o VR de outubro/2025"*

1. **Organize suas planilhas** em uma pasta:
   - Arquivo de colaboradores (.xlsx)
   - Arquivo de afastamentos (.xlsx)  
   - Arquivo de feriados (.xlsx)
   - Arquivo de desligamentos (.xlsx)

2. **Na tela principal, clique em "Selecionar Diretório"**

3. **Escolha a pasta** com suas planilhas organizadas

4. **Aguarde a validação** - aparecerá "Válido" se a pasta contém arquivos Excel

5. **Clique em "Iniciar Processamento"** quando disponível

6. **Acompanhe o progresso** na área de resultados:
   - Status: "Análise em andamento, por favor aguarde..."
   - Quando concluído: mostra número de colaboradores processados
   - Arquivo é salvo automaticamente na pasta Downloads

### **Cenário 3: Chat Inteligente**
*"Tenho dúvidas sobre dados processados"*

1. **Processe os dados primeiro** (Cenário 2)
2. **Localize o chat** na parte inferior da tela
3. **Clique para expandir** o painel do chat
4. **Digite suas perguntas** na caixa de texto
5. **Aguarde a resposta** do agente de IA

**Exemplos de perguntas:**
```
Você: "Quantos colaboradores foram processados?"
Agente: "Foram processados X colaboradores no último processamento."

Você: "Qual o valor total de VR calculado?"
Agente: "O valor total calculado foi de R$ X,XX baseado nos 
        dados consolidados das planilhas."

Você: "Explique como foi feito o cálculo?"
Agente: "O cálculo foi realizado considerando os dias úteis 
        do mês, descontando afastamentos e feriados conforme 
        as regras de cada sindicato."
```

**Nota:** O chat precisa de dados processados para funcionar com contexto completo.

### **Cenário 4: Análise de Anomalias**
*"Suspeito que há problemas nos dados deste mês"*

1. **Clique na aba "🔄 Workflows"** na interface
2. **Selecione "Análise de Anomalias"** no dropdown
3. **Clique em "▶️ Executar"**
4. **Aguarde a análise** (1-2 minutos)

**Ou use o chat:**
1. **Vá para a aba "💬 Chat"**
2. **Pergunte: "Verifique se há algo estranho nos dados de outubro"**

**O agente pode retornar:**
```
🚨 3 anomalias detectadas:
1. Colaborador MAT001: VR 340% acima da média do SINDAC
2. Colaborador MAT445: Possível duplicata de matrícula  
3. Colaborador MAT892: Data de admissão posterior à data de demissão

📊 Relatório detalhado disponível para download
💡 Recomendação: Revisar manualmente os 3 casos antes do processamento final
```

## 🎛️ Dashboard Integrado

### **O que você vê na Interface Principal:**
- 📊 **Status em tempo real** dos workflows em execução
- 📈 **Métricas de performance** (throughput, tempo de execução)
- 🔍 **Anomalias detectadas** com opções de ação
- 📋 **Histórico de processamentos** anteriores
- ⚙️ **Acesso rápido às configurações**

## 🔧 Configurações Básicas

### **Arquivo de Configuração**
O agente cria `config/agent.yaml` com suas configurações:

```yaml
# Configurações principais
llm:
  provider: "openai"  # ou "claude", "local"
  api_key: "sua-chave-aqui"
  model: "gpt-4"

files:
  input_dir: "./files"
  output_dir: "./output"
  backup_enabled: true

cache:
  enabled: true
  ttl_hours: 24
  max_size: 1000

notifications:
  email:
    enabled: true
    smtp_host: "smtp.gmail.com"
    recipients: ["gestor@empresa.com"]
```

### **Configurações Avançadas**
```bash
# Ver todas as configurações disponíveis
./brx-agente config --show-all

# Configurar apenas cache
./brx-agente config cache

# Testar configurações
./brx-agente config test
```

## 📂 Estrutura de Arquivos

Após a configuração inicial, você terá:

```
BrxAgente-desafio4/
├── brx-agente          # Executável principal
├── config/
│   ├── agent.yaml      # Configurações do agente
│   └── performance.yaml # Configurações de performance
├── files/              # 📁 SUAS PLANILHAS AQUI
│   ├── colaboradores.xlsx
│   ├── afastamentos.xlsx
│   └── feriados.xlsx
├── output/             # Resultados gerados
│   ├── vr-outubro-2025.xlsx
│   └── relatorio-anomalias.xlsx
├── logs/               # Logs detalhados
└── cache/              # Cache do agente (automático)
```

## ⚡ Funcionalidades Essenciais

### **Workflows Principais (Aba "🔄 Workflows")**
- **"Processamento Completo de VR"**: Executa todo o fluxo automaticamente
- **"Validação de Planilhas"**: Apenas valida a estrutura e consistência
- **"Análise de Anomalias"**: Detecta problemas nos dados
- **"Geração de Relatórios"**: Cria relatórios personalizados

### **Monitoramento (Dashboard Principal)**
- **Status do Agente**: Visível no topo da interface
- **Métricas**: Painel lateral com estatísticas
- **Logs**: Acessíveis via botão "📋 Ver Logs"
- **Cache**: Gerenciado automaticamente (botão "🧹 Limpar Cache")

### **Chat e Consultas (Aba "💬 Chat")**
- **Digite perguntas** na caixa de texto
- **Exemplos úteis**:
  - "Quantos colaboradores temos?"
  - "Qual o VR total do SINDPD?"
  - "Há anomalias nos dados atuais?"

## 🚨 Primeiros Passos - Checklist

### ✅ **Antes do Primeiro Uso**
- [ ] Go 1.19+ e Node.js 16+ instalados
- [ ] Projeto clonado e compilado com `wails build`
- [ ] Aplicação iniciada com sucesso
- [ ] Planilhas organizadas em pasta dedicada
- [ ] Chave da API configurada na interface
- [ ] Teste de conexão realizado

### ✅ **Primeiro Processamento**
- [ ] Backup das planilhas originais
- [ ] Validação prévia via interface (Workflows → Validação)
- [ ] Análise de anomalias executada
- [ ] Processamento em dados de teste primeiro
- [ ] Conferência manual dos resultados gerados

### ✅ **Configuração Produtiva**
- [ ] Notificações configuradas na interface
- [ ] Backup automático habilitado nas configurações
- [ ] Interface acessível para toda a equipe
- [ ] Procedimentos de contingência documentados
- [ ] Treinamento da equipe realizado

## 🆘 Problemas Comuns

### **"Aplicação não inicia"**
1. **Verificar dependências**: Execute `go mod tidy` no terminal
2. **Recompilar**: Execute `wails build` novamente
3. **Verificar permissões**: Confirmar que pode executar o arquivo gerado

### **"Erro ao ler planilhas"**
1. **Verificar pasta selecionada**: Botão "📁 Selecionar Pasta" na interface
2. **Ver detalhes do erro**: Botão "📋 Ver Logs" no dashboard
3. **Validar formato**: Deve ser arquivos .xlsx (não .xls ou .csv)

### **"API Key inválida"**
1. **Ir em Configurações**: Botão "⚙️" na interface
2. **Aba "🔑 API"**: Inserir chave válida
3. **Testar conexão**: Botão "🧪 Testar Conexão"

## 🎯 Próximos Passos

Depois de dominar o básico:

1. 📖 **[Guia Completo do Usuário](user-guide.md)** - Funcionalidades avançadas
2. ⚙️ **[Workflows Detalhados](workflows.md)** - Customização de processos
3. 🔧 **[Configurações Avançadas](../developer/architecture.md)** - Otimizações
4. 💡 **[Exemplos Práticos](examples/)** - Cases reais de uso

---

**💡 Dica**: Use o chat integrado na aplicação para obter ajuda sobre qualquer funcionalidade!

*Pronto para processar seu primeiro VR? Abra a aplicação e vá em Workflows → Processamento Completo! 🚀*