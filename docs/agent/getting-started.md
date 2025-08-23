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
- 📊 **Dashboard Principal**: Status do agente e métricas
- 🗂️ **Seletor de Diretório**: Para escolher pasta com planilhas Excel
- ⚙️ **Configurações**: API keys e parâmetros
- 💬 **Chat Inteligente**: Para fazer perguntas ao agente
- 🔄 **Workflows**: Botões para executar processamentos

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

1. **Coloque suas planilhas** na pasta selecionada:
   - colaboradores.xlsx
   - afastamentos.xlsx  
   - feriados.xlsx

2. **Na tela principal, clique em "Processar VR"**

3. **Acompanhe o progresso** na interface:
   - ✅ Lendo planilhas... (30s)
   - ✅ Validando dados... (45s)
   - 🔍 Detectando anomalias... (2m)
   - ⚙️ Calculando VR... (3m)
   - 📊 Gerando relatórios... (1m)

4. **Baixe os resultados** através dos botões na interface

### **Cenário 3: Chat Inteligente**
*"Tenho dúvidas sobre políticas de VR"*

1. **Clique na aba "💬 Chat"** na interface
2. **Digite suas perguntas** na caixa de texto
3. **Receba respostas inteligentes** do agente

**Exemplos de perguntas:**
```
Você: "Quantos colaboradores do SINDPD temos ativos?"
Agente: "Encontrei 247 colaboradores do SINDPD ativos. 
        Destes, 234 são elegíveis para VR. 
        Deseja ver a distribuição por setor?"

Você: "Um colaborador admitido no dia 20 tem direito a VR integral?"
Agente: "Colaboradores admitidos após o dia 15 têm direito a VR 
        proporcional. Para admissão no dia 20, o cálculo seria:
        - Dias úteis restantes: 8 dias
        - VR proporcional: R$ 184,00 (considerando VR base R$ 460,00)
        
        Fonte: Política interna VR-2025, seção 3.2"
```

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