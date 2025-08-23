# 🆘 Troubleshooting - Agente de IA BrxAgente

Guia completo para diagnóstico e resolução de problemas na aplicação desktop BrxAgente.

## 📑 Índice

1. [Problemas de Instalação](#-problemas-de-instalação)
2. [Problemas de Configuração](#-problemas-de-configuração)
3. [Problemas de Processamento](#-problemas-de-processamento)
4. [Problemas do Chat](#-problemas-do-chat)
5. [Problemas de Performance](#-problemas-de-performance)
6. [Logs e Diagnóstico](#-logs-e-diagnóstico)
7. [FAQ](#-faq)

## 🔧 Problemas de Instalação

### **❌ Problema: "Aplicação não inicia"**

**Sintomas:**
- Clique duplo não funciona
- Nada acontece ao executar
- Erro "arquivo não encontrado"

**Soluções:**

1. **Verificar dependências:**
```bash
# Verificar Go instalado
go version
# Deve retornar >= 1.19

# Verificar Node.js (para desenvolvimento)
node --version
# Deve retornar >= 16

# Verificar Wails CLI
wails version
```

2. **Recompilar a aplicação:**
```bash
cd BrxAgente-desafio4

# Limpar builds anteriores
rm -rf build/

# Recompilar
wails build

# Ou para desenvolvimento
wails dev
```

3. **Verificar permissões:**
```bash
# Linux/Mac: Adicionar permissão de execução
chmod +x ./build/bin/BrxAgente-desafio4

# Windows: Executar como administrador se necessário
```

---

### **❌ Problema: "Erro ao compilar"**

**Sintomas:**
```
Error: failed to build frontend
npm ERR! Missing script: "build"
```

**Soluções:**

1. **Instalar dependências do frontend:**
```bash
cd frontend/
npm install

# Ou usar yarn
yarn install
```

2. **Verificar scripts no package.json:**
```json
{
  "scripts": {
    "build": "react-scripts build",
    "dev": "react-scripts start"
  }
}
```

3. **Limpar cache e reinstalar:**
```bash
rm -rf frontend/node_modules/
rm frontend/package-lock.json
cd frontend/
npm install
```

---

### **❌ Problema: "Erro de dependências Go"**

**Sintomas:**
```
go: module not found
cannot find package
```

**Soluções:**

1. **Atualizar módulos Go:**
```bash
go mod tidy
go mod download
```

2. **Verificar versão do Go:**
```bash
go version
# Se < 1.19, atualize o Go
```

3. **Limpar cache do Go:**
```bash
go clean -modcache
go mod download
```

## ⚙️ Problemas de Configuração

### **❌ Problema: "API Key inválida"**

**Sintomas:**
- "Chave da API do OpenAI inválida"
- Chat não responde
- Erro de autenticação

**Soluções:**

1. **Verificar formato da chave:**
```
✅ Formato correto: sk-proj-abc123...
❌ Formato incorreto: sk-abc123... (sem 'proj')
❌ Formato incorreto: chave incompleta
```

2. **Testar chave manualmente:**
```bash
# Teste via curl (substitua YOUR_KEY)
curl -H "Authorization: Bearer sk-proj-YOUR_KEY" \
     -H "Content-Type: application/json" \
     -d '{"model":"gpt-3.5-turbo","messages":[{"role":"user","content":"test"}],"max_tokens":5}' \
     https://api.openai.com/v1/chat/completions
```

3. **Verificar configurações na interface:**
   - Vá em "Configurações" → "API"
   - Cole a chave novamente
   - Clique em "Testar Conexão"
   - Aguarde confirmação ✅

---

### **❌ Problema: "Ollama não conecta"**

**Sintomas:**
- "Configuração do Ollama inválida"
- Timeout na conexão
- Erro de URL

**Soluções:**

1. **Verificar se Ollama está rodando:**
```bash
# Verificar se está ativo
curl http://localhost:11434/api/version

# Iniciar Ollama se necessário
ollama serve
```

2. **Testar modelo específico:**
```bash
# Baixar modelo se necessário
ollama pull llama2

# Testar modelo
ollama run llama2 "Olá, como você está?"
```

3. **Configurar na interface:**
   - URL: `http://localhost:11434`
   - Modelo: `llama2` (ou outro instalado)
   - Timeout: `30` segundos

---

### **❌ Problema: "Configurações não salvam"**

**Sintomas:**
- Configurações resetam ao reiniciar
- "Falha ao salvar a configuração"
- Arquivo config.json não encontrado

**Soluções:**

1. **Verificar permissões do diretório:**
```bash
# Verificar se existe diretório config
ls -la config/

# Criar se não existir
mkdir -p config/
chmod 755 config/
```

2. **Verificar espaço em disco:**
```bash
# Linux/Mac
df -h .

# Windows
dir
```

3. **Executar como administrador** (Windows) se necessário

## 📊 Problemas de Processamento

### **❌ Problema: "Diretório não encontrado"**

**Sintomas:**
- "Diretório não encontrado: /path/to/planilhas"
- Seletor de arquivos não funciona
- Planilhas não são reconhecidas

**Soluções:**

1. **Verificar caminho:**
```bash
# Verificar se diretório existe
ls -la /caminho/para/planilhas/

# Verificar arquivos .xlsx
ls -la /caminho/para/planilhas/*.xlsx
```

2. **Usar seletor de diretório da interface:**
   - Clique em "📁 Selecionar Pasta"
   - Navegue até a pasta correta
   - Confirme a seleção

3. **Verificar nomes dos arquivos esperados:**
   - ✅ `colaboradores.xlsx`
   - ✅ `afastamentos.xlsx`
   - ✅ `feriados.xlsx`

---

### **❌ Problema: "Erro ao ler planilhas"**

**Sintomas:**
```
Error: failed to read Excel file
Planilha corrompida
Formato não suportado
```

**Soluções:**

1. **Verificar formato dos arquivos:**
   - Deve ser `.xlsx` (não `.xls` ou `.csv`)
   - Salvar como Excel 2007+ se necessário

2. **Verificar integridade:**
   - Abrir arquivos no Excel/LibreOffice
   - Verificar se não estão corrompidos
   - Resalvar se necessário

3. **Verificar estrutura:**
```
📋 Colaboradores.xlsx deve ter colunas:
• MATRICULA
• NOME  
• SINDICATO
• DATA_ADMISSAO
• STATUS (ativo/inativo)

📋 Afastamentos.xlsx deve ter colunas:
• MATRICULA
• TIPO_AFASTAMENTO
• DATA_INICIO
• DATA_FIM

📋 Feriados.xlsx deve ter colunas:
• DATA
• DESCRICAO
• TIPO (nacional/estadual/municipal)
```

---

### **❌ Problema: "Workflow trava/não progride"**

**Sintomas:**
- Barra de progresso parou
- "Etapa 3 de 6" há muito tempo
- Interface não responde

**Soluções:**

1. **Verificar logs em tempo real:**
   - Clique em "📊 Ver Logs"
   - Procure por erros ou mensagens de stuck
   - Identifique última operação bem-sucedida

2. **Recursos do sistema:**
```bash
# Verificar uso de CPU/Memória
# Linux/Mac
top -p $(pgrep BrxAgente)

# Windows (Task Manager)
# Procurar por BrxAgente-desafio4.exe
```

3. **Ações corretivas:**
   - Se travou: Clique em "❌ Cancelar"
   - Aguarde limpeza (pode demorar 30s)
   - Reinicie workflow com menos dados (teste)

---

### **❌ Problema: "Muitas anomalias detectadas"**

**Sintomas:**
```
🚨 47 anomalias críticas detectadas
Processing stopped due to excessive anomalies
Threshold exceeded: >5 critical anomalies
```

**Soluções:**

1. **Revisar dados de entrada:**
   - Verificar se planilhas estão atualizadas
   - Confirmar formatos de data (DD/MM/YYYY)
   - Validar matrículas e sindicatos

2. **Ajustar sensibilidade:**
   - Ir em "Configurações" → "Detecção"
   - Reduzir sensibilidade para "Média" ou "Baixa"
   - Aumentar threshold multiplier (2.0x → 2.5x)

3. **Análise individual:**
   - Clicar em "🔍 Ver Detalhes" de cada anomalia
   - Determinar se são erros reais ou falsos positivos
   - Corrigir dados na fonte se necessário

## 💬 Problemas do Chat

### **❌ Problema: "Chat não responde"**

**Sintomas:**
- Mensagem fica "enviando..." indefinidamente
- Erro "Falha ao obter resposta da IA"
- Timeout na resposta

**Soluções:**

1. **Verificar conexão com API:**
   - Ir em "Configurações" → "API"
   - Clicar em "🧪 Testar Conexão"
   - Aguardar resultado

2. **Verificar contexto de dados:**
   - Garantir que dados foram processados
   - Verificar se `SetChatContext()` foi executado
   - Dados vazios = respostas limitadas

3. **Testar pergunta simples:**
   - Começar com: "Olá, você está funcionando?"
   - Perguntas complexas podem falhar se contexto vazio

---

### **❌ Problema: "Respostas inconsistentes ou incorretas"**

**Sintomas:**
- Números não batem com planilhas
- Informações contraditórias
- "Não encontrei dados sobre..."

**Soluções:**

1. **Atualizar contexto:**
   - Reprocessar dados mais recentes
   - Verificar se dados foram carregados corretamente
   - Status deve mostrar número correto de colaboradores

2. **Reformular pergunta:**
```
❌ Vago: "Quantos colaboradores temos?"
✅ Específico: "Quantos colaboradores ativos do SINDPD temos?"

❌ Ambíguo: "Há problemas nos dados?"  
✅ Claro: "Há anomalias nos valores de VR detectadas?"
```

3. **Verificar logs do sistema:**
   - Procurar por erros de contexto
   - Confirmar que dados foram indexados

## ⚡ Problemas de Performance

### **❌ Problema: "Processamento muito lento"**

**Sintomas:**
- Processamento > 30 minutos
- Throughput < 10 colaboradores/min
- Interface travada

**Soluções:**

1. **Verificar recursos do sistema:**
```
Requisitos mínimos:
• RAM: 4GB livre
• CPU: 2+ cores
• Disco: 1GB espaço livre
```

2. **Otimizar configurações:**
   - Ir em "Configurações" → "Performance"
   - Aumentar workers: 4 → 8
   - Ativar cache se desabilitado
   - Reduzir validações não essenciais

3. **Processar em lotes menores:**
   - Filtrar por sindicato específico
   - Processar meses separadamente
   - Usar apenas colaboradores ativos

---

### **❌ Problema: "Alto uso de memória"**

**Sintomas:**
```
Memory usage: 1.2GB / 4GB (90%)
System running low on memory
Garbage collection frequent
```

**Soluções:**

1. **Reduzir tamanho do cache:**
   - Configurações → Performance
   - Cache size: 1GB → 500MB
   - TTL: 24h → 12h

2. **Processar dados menores:**
   - Filtrar registros desnecessários
   - Processar por sindicato
   - Limitar período de análise

3. **Reiniciar aplicação periodicamente:**
   - Fechar aplicação após processamentos grandes
   - Reiniciar para limpar memória
   - Não deixar executando 24/7

## 📋 Logs e Diagnóstico

### **Como Acessar Logs**

1. **Interface da Aplicação:**
   - Dashboard → "📊 Ver Logs"
   - Workflow → "📋 Logs do Workflow"
   - Configurações → "🔍 Logs do Sistema"

2. **Arquivos de Log:**
```bash
# Localização dos logs
logs/
├── application.log    # Log geral da aplicação
├── workflow.log      # Logs específicos de workflows
├── chat.log         # Logs do sistema de chat
└── errors.log       # Apenas erros críticos
```

### **Interpretando Logs**

**Níveis de Log:**
- `[INFO ]` - Informações normais
- `[WARN ]` - Avisos (não críticos)
- `[ERROR]` - Erros que requerem atenção
- `[DEBUG]` - Informações detalhadas (apenas dev)

**Exemplos de Logs Importantes:**

```log
✅ Log Normal:
[INFO ] 14:32:15 - Workflow started: analise-vr-mensal
[INFO ] 14:32:45 - Excel reading completed: 2847 records
[INFO ] 14:35:22 - Workflow completed successfully

⚠️  Log com Aviso:
[WARN ] 14:33:12 - Anomaly detected: MAT001234 value 340% above average
[WARN ] 14:33:15 - Missing data for 3 collaborators in afastamentos.xlsx

❌ Log com Erro:
[ERROR] 14:34:22 - Failed to read Excel file: files/colaboradores.xlsx
[ERROR] 14:34:23 - OpenAI API error: rate limit exceeded
[ERROR] 14:35:01 - Workflow terminated due to critical error
```

### **Coleta de Informações para Suporte**

Antes de abrir um ticket, colete:

1. **Informações do Sistema:**
```bash
# Versão da aplicação
cat version.txt

# Sistema operacional
uname -a    # Linux/Mac
ver         # Windows

# Versões das dependências
go version
node --version
wails version
```

2. **Logs Relevantes:**
   - Últimas 50 linhas do log de erro
   - Log completo do workflow que falhou
   - Screenshot da tela de erro

3. **Configurações (sem dados sensíveis):**
   - Configurações de performance
   - Parâmetros do workflow
   - Tamanho dos arquivos processados

## ❓ FAQ

### **P: Posso processar planilhas com mais de 10.000 colaboradores?**
**R:** Sim, mas recomendamos:
- RAM mínima de 8GB
- Processar por sindicato
- Aumentar timeout para 60+ minutos
- Usar máquina dedicada

### **P: É possível executar múltiplos workflows simultaneamente?**
**R:** Não. O sistema executa um workflow por vez para garantir integridade. Use workflow combinado se precisar de operações múltiplas.

### **P: Como fazer backup dos dados processados?**
**R:** A aplicação cria backups automáticos em:
- `backups/[data]/planilhas-originais/`
- `backups/[data]/resultados-processados/`
- Configure retenção em Configurações → Backup

### **P: Posso usar o sistema offline?**
**R:** Parcialmente:
- ✅ Processamento de planilhas (sem IA)
- ✅ Validação e cálculos básicos
- ❌ Chat inteligente (requer internet)
- ❌ Análise preditiva (requer API)

### **P: Como atualizar a aplicação?**
**R:** 
1. Baixar nova versão do repositório
2. Fazer backup das configurações
3. Executar `wails build` 
4. Restaurar configurações se necessário

### **P: O sistema guarda histórico de conversas do chat?**
**R:** Sim, por 30 dias por padrão. Configure retenção em Configurações → Chat → Histórico.

---

## 🚨 Situações de Emergência

### **Sistema Travado Completamente**
1. Force-close a aplicação (Ctrl+Alt+T / Cmd+Q)
2. Verifique se processo ainda está rodando
3. Kill processo se necessário
4. Reinicie aplicação
5. Verifique se dados foram corrompidos

### **Dados Corrompidos**
1. **NÃO** execute novos processamentos
2. Pare todos os workflows ativos  
3. Restaure backup mais recente
4. Verifique integridade dos backups
5. Execute validação completa antes de continuar

### **Performance Crítica**
1. Reduza workers para 2
2. Limpe todos os caches
3. Reinicie aplicação
4. Processe apenas dados essenciais
5. Considere migrar para máquina mais potente

---

**🆘 Precisa de mais ajuda?**
1. Consulte a [documentação completa](README.md)
2. Verifique [issues conhecidas](https://github.com/vstram/BrxAgente-desafio4/issues)
3. Abra novo issue com logs e informações coletadas

*Lembre-se: a maioria dos problemas pode ser resolvida com logs adequados e análise sistemática! 🔍*