# 📖 Guia Completo do Usuário - Agente de IA BrxAgente

Este guia detalha todas as funcionalidades da aplicação desktop BrxAgente e como utilizá-las efetivamente no seu trabalho diário.

## 📑 Índice

1. Interface da Aplicação
1. Funcionalidades Principais
1. Chat Avançado com IA
1. Sistema de Agentes de IA (LangflowGo)
1. Configurações
1. Melhores Práticas

## 🖥️ Interface da Aplicação

### **Layout Principal**

A aplicação desktop BrxAgente possui uma interface simples e funcional:

```
┌─────────────────────────────────────────────────────┐
│  🏷️ Automação de VR/VA                              │
│  Sistema de processamento automatizado              │
├─────────────────────────────────────────────────────┤
│                                                     │
│  📁 Seleção de Planilhas                           │
│  ┌─────────────────────────────────────────────────┐ │
│  │ [📁 Selecionar Diretório]                      │ │
│  │ Diretório: /caminho/para/planilhas             │ │
│  │ Status: ✅ Válido                               │ │
│  └─────────────────────────────────────────────────┘ │
│                                                     │
│  ▶️ Processamento                                   │
│  ┌─────────────────────────────────────────────────┐ │
│  │ [▶️ Iniciar Processamento]                      │ │
│  └─────────────────────────────────────────────────┘ │
│                                                     │
│  📊 Resultados                                      │
│  ┌─────────────────────────────────────────────────┐ │
│  │ Status: Análise concluída com sucesso!         │ │
│  │ X colaboradores processados                     │ │
│  │ Arquivo salvo em Downloads                      │ │
│  └─────────────────────────────────────────────────┘ │
│                                                     │
├─────────────────────────────────────────────────────┤
│ [⚙️ Configurações]               💬 Chat (expandir) │
└─────────────────────────────────────────────────────┘
```

### **Seções da Interface**

#### **1. Cabeçalho**
- 🏷️ **Logo e Título**: Identificação da aplicação
- 📝 **Subtítulo**: Descrição do propósito

#### **2. Área Principal**
- **Seleção de Planilhas**: Botão para escolher diretório
- **Processamento**: Botão para iniciar análise (aparece após seleção válida)
- **Resultados**: Exibe progresso e resultados do processamento

#### **3. Rodapé**
- ⚙️ **Configurações**: Acesso ao modal de configuração
- 💬 **Chat**: Painel expansível de chat inteligente

## 🔍 Funcionalidades Principais

### **1. Consultor de Políticas**

**Como Usar:**
1. Clique em **"Chat"** no menu principal
2. Digite sua pergunta na caixa de texto
3. Aguarde a resposta do agente
4. Use **histórico de conversa** para contexto

**Tipos de Perguntas Suportadas:**

**Elegibilidade e Regras:**
```
👤 "Diretores têm direito a VR?"
🤖 "Não. Diretores não são elegíveis para VR conforme 
    Política VR-2025, Seção 2.3. Isso inclui:
    • Diretores Estatutários
    • Diretores Executivos  
    • Diretores de Área
    
    Exceção: Diretores com salário < R$ 15.000 podem 
    ser incluídos mediante aprovação especial."
```

**Cálculos Específicos:**
```
👤 "Como calcular VR para licença médica de 20 dias?"
🤖 "Para licença médica > 15 dias, usar cálculo proporcional:
    
    📊 Fórmula: (Dias trabalhados ÷ Dias úteis totais) × VR base
    
    📝 Exemplo:
    • Dias úteis outubro: 22 dias
    • Dias trabalhados: 2 dias (22-20)
    • Cálculo: (2 ÷ 22) × R$ 460,00 = R$ 41,82
    
    📚 Fonte: Manual VR-2025, Seção 5.4"
```

## 🤖 Chat Avançado com IA

### **Tipos de Consulta Suportadas**

O sistema de chat inteligente pode responder diferentes tipos de perguntas com base em conhecimento estruturado e dados processados.

#### 📋 Consultas sobre Políticas
O sistema pode responder perguntas sobre políticas de VR com base em sua base de conhecimento:

**Exemplos:**
- "Diretores têm direito a VR?"
- "Estagiários podem receber Vale Refeição?"
- "Qual a política para colaboradores de 4h/dia?"

**Resposta típica:**
> ## 📋 Consulta de Política
> **Pergunta:** Diretores têm direito a VR?
> **Resposta:** Não, diretores não têm direito ao Vale Refeição.
> **Fonte:** Política VR_003 - Manual de RH v2.1
> **Confiança:** 95% (Muito Alta)

#### 🧮 Cálculos Específicos
Explica como calcular VR para cenários específicos:

**Exemplos:**
- "Como calcular VR para licença médica de 20 dias?"
- "Colaborador admitido dia 25, qual valor?"
- "Desligamento comunicado dia 18, tem direito?"

**Resposta típica:**
> ## 🧮 Cálculo de VR
> **Cenário:** Licença médica de 20 dias
> 
> Para licença médica > 15 dias:
> - **Fórmula:** (Dias trabalhados ÷ Dias úteis) × VR base
> - **Exemplo:** (2 ÷ 22) × R$ 460,00 = R$ 41,82
> 
> **Fonte:** Manual VR-2025, Seção 5.4
> **Confiança:** 90% (Alta)

#### 📊 Dados Processados
Consulta dados do último processamento realizado:

**Exemplos:**
- "Quantos colaboradores foram processados?"
- "Qual o total de VR este mês?"
- "Distribuição por empresa"

### **Interface de Chat**

**Layout do Chat:**
```
┌─────────────────────────────────────────────────────┐
│ 💬 Consultor Inteligente - Agente VR              │
├─────────────────────────────────────────────────────┤
│                                                     │
│ 👤 Quantos colaboradores do SINDPD temos?          │
│                                                     │
│ 🤖 ## 📊 Consulta de Dados                         │
│    **Colaboradores SINDPD Processados:**           │
│    • Total: 247 colaboradores                      │
│    • Elegíveis para VR: 234 (94.7%)               │
│    • Não elegíveis: 13 (5.3%)                     │
│                                                     │
│    **Detalhamento:**                                │
│    • Ativos e elegíveis: 234                       │
│    • Em licença: 8                                 │
│    • Estagiários: 3 (não elegíveis)               │
│    • Diretores: 2 (não elegíveis)                 │
│                                                     │
│    **Fonte:** Último processamento (Out/2025)      │
│    **Confiança:** 100% (Dados Reais)              │
│                                                     │
├─────────────────────────────────────────────────────┤
│ Digite sua pergunta...                    [Enviar] │
└─────────────────────────────────────────────────────┘
```

**Funcionalidades do Chat:**
- 🤖 **Classificação Inteligente**: Identifica automaticamente o tipo de pergunta
- 📋 **Base de Conhecimento**: Respostas baseadas em políticas oficiais
- 📊 **Contexto de Dados**: Acesso aos dados processados (1794+ colaboradores)
- 🐌⚡ **Duas Modalidades**: Análise Detalhada (padrão) ou Análise Rápida
- 💾 **Histórico Persistente**: Conversas salvas
- 🔍 **Citação de Fontes**: Referencias oficiais nas respostas
- 📈 **Formatação Inteligente**: Profissional (rápida) ou Raw (detalhada)
- 🎯 **Confiança**: Indica nível de certeza da resposta

### **🎯 Modalidades de Análise**

#### **🐌 Análise Detalhada (Modo Padrão)**

**Características:**
- ✅ **Cache desabilitado** por padrão para máxima precisão
- ✅ **Dados completos**: Envia todos os 1794+ colaboradores individuais (~216KB)
- ✅ **Respostas brutas**: Direto da LLM, sem formatação pré-definida
- ⏱️ **Tempo**: ~13 segundos para análises profundas
- 🧠 **Capacidade**: Processa OpenAI com contexto completo

**Ideal Para:**
- 📊 Análises estatísticas complexas
- 🔍 Detecção de anomalias e padrões
- 📈 Comparações multi-dimensionais
- 🎯 Relatórios customizados e detalhados

**Exemplos de Perguntas Avançadas:**
```
💬 "Identifique colaboradores com valores de VR que desviam mais de 20% da média do seu sindicato e analise se há padrões relacionados aos dias úteis efetivos"

💬 "Compare a eficiência de VR por sindicato considerando a relação entre valor médio,variabilidade dos dias úteis e distribuição geográfica"

💬 "Analise se existem colaboradores que, baseado no padrão de dias úteis e valores de VR do seu sindicato, podem estar recebendo valores incorretos"
```

#### **⚡ Análise Rápida (Modo Opcional)**

**Características:**
- 🚀 **Cache habilitado** para máxima velocidade
- 📊 **Dados agregados**: Informações otimizadas (~702 chars)
- 🎨 **Respostas formatadas**: Interface padronizada e consistente
- ⏱️ **Tempo**: ~2 segundos para consultas diretas
- 💫 **PolicyConsultant**: Acesso à base de conhecimento estruturada

**Como Ativar:**
1. **Clique no ícone do cache** no cabeçalho do chat (🐌 → ⚡)
2. **Observe a mudança** no tooltip e comportamento
3. **Use para consultas** simples e validações rápidas

**Ideal Para:**
- ❓ Consultas simples e diretas
- 📋 Validação de políticas
- 🔄 Verificações rápidas de status
- 📊 Dados agregados básicos

**Exemplos de Perguntas Rápidas:**
```
💬 "Quantos colaboradores foram processados neste mês?"

💬 "Qual o valor total de VR calculado?"

💬 "Estagiários têm direito a VR segundo as políticas?"

💬 "Há alguma anomalia crítica detectada?"
```

### **🔄 Alternando Entre Modalidades**

**Na Interface:**
- **🐌 Padrão**: Aparece automaticamente ao abrir o chat
- **⚡ Rápida**: Clique no ícone do cache para alternar
- **Tooltip**: Mostra claramente qual modo está ativo
- **Feedback**: Sistema confirma a mudança com mensagem explicativa

## 🤖 Sistema de Agentes de IA (LangflowGo)

### **Arquitetura dos Agentes**

O BrxAgente utiliza uma arquitetura avançada de agentes de IA baseada na framework **LangflowGo** (baseado no LangChain para Go), proporcionando funcionalidades inteligentes e automatizadas para o processamento de VR.

#### **🧠 Agente Principal (VRAgent)**

O sistema é centrado no `VRAgent`, que integra múltiplos componentes especializados:

**Componentes Ativos:**
- 💬 **Chat Service**: Interface para comunicação com LLMs (OpenAI/Ollama)
- 🧮 **Question Classifier**: Classificação inteligente de perguntas
- ⚡ **Performance Optimizer**: Otimização de contexto e cache
- 📊 **Response Formatter**: Formatação profissional de respostas
- 🛠️ **Tool Registry**: Registro de ferramentas especializadas

**Componentes Backend (Implementados mas não integrados na UI):**
- 🔍 **Auditor Inteligente**: Detecção de anomalias e inconsistências
- 🔄 **Workflow Orchestrator**: Orquestração de processos automatizados
- 📈 **Intelligence Analyzer**: Análise estatística avançada

### **🛠️ Ferramentas do Agente**

O sistema possui um **Tool Registry** com ferramentas especializadas:

#### **📋 PolicyConsultantTool**
- **Função**: Consulta à base de conhecimento de políticas de VR
- **Recursos**:
  - Cache inteligente para respostas rápidas
  - Classificação automática de perguntas
  - Citação de fontes oficiais
  - Níveis de confiança nas respostas

#### **📊 Ferramentas de Cálculo**
- **CalculoVRTool**: Cálculos específicos de VR
- **ValidacaoDataTool**: Validação de datas e períodos
- **ConsolidacaoTool**: Consolidação de dados de múltiplas bases

#### **📁 Ferramentas de Excel**
- **ExcelReaderTool**: Leitura avançada de planilhas
- **ExcelValidatorTool**: Validação de estrutura e dados
- **ExcelWriterTool**: Geração de relatórios

#### **⚖️ ComplianceAssistant**
- **Função**: Assistência em questões de conformidade
- **Recursos**: Análise de regras e regulamentações

### **🔄 Sistema de Workflows (Backend)**

O sistema implementa workflows automatizados através do **Workflow Orchestrator**:

#### **📋 Simple Validation Workflow**
- Validação básica de dados e estruturas
- Verificação de consistência
- Relatório de status

#### **📊 VR Processing Workflow**
- Processamento completo mensal de VR
- Cálculos automatizados
- Geração de relatórios finais

#### **📈 Reporting Workflow**
- Geração de relatórios customizados
- Análises estatísticas
- Exportação em múltiplos formatos

**Status:** ⚠️ _Workflows implementados no backend, pendente integração na interface_

### **🔍 Auditor Inteligente (Backend)**

O **Intelligence Analyzer** oferece capacidades avançadas de auditoria:

#### **Detecção de Anomalias**
- Análise estatística de valores de VR
- Identificação de outliers por sindicato
- Detecção de padrões suspeitos
- Cálculo de scores de confiança

#### **Relatórios de Auditoria**
- Relatórios detalhados de inconsistências
- Sugestões de correção automatizadas
- Métricas de qualidade dos dados
- Rastreabilidade completa

**Status:** ⚠️ _Auditor implementado no backend, pendente integração na interface_

### **⚡ Modos de Operação**

#### **🐌 Modo Análise Detalhada (Padrão)**
- **LLM**: Processamento completo com todos os dados
- **Context**: Dados completos de 1794+ colaboradores (~216KB)
- **Cache**: Desabilitado para máxima precisão
- **Tempo**: ~13 segundos
- **Ideal**: Análises estatísticas complexas

#### **⚡ Modo Análise Rápida**
- **LLM**: PolicyConsultant com cache otimizado
- **Context**: Dados agregados (~702 chars)
- **Cache**: Habilitado para velocidade
- **Tempo**: ~2 segundos
- **Ideal**: Consultas simples e validações

### **🎯 Roteamento Inteligente**

O agente utiliza **classificação automática** para rotear perguntas:

```
🤖 Question Classifier
├── 📚 PolicyQuestion → PolicyConsultantTool
├── 🧮 CalculationQuestion → PolicyConsultant + Fallback
├── 📊 ProcessedDataQuestion → Dados Reais
├── 🤔 WhatIfQuestion → PolicyConsultantTool
└── ❓ UnknownQuestion → Tentativa + Fallback
```

### **💾 Sistema de Cache Inteligente**

**Cache Habilitado (⚡):**
- Respostas formatadas padronizadas
- Dados agregados otimizados
- PolicyConsultant com knowledge base

**Cache Desabilitado (🐌):**
- Dados completos individuais de colaboradores
- Respostas raw diretas da LLM
- Análise estatística profunda

### **🔧 Configuração dos Agentes**

**Parâmetros Configuráveis:**
- Thresholds de confiança do classificador
- Timeout de workflows
- Limites de execuções concorrentes
- Níveis de logging detalhado
- Habilitação de rollback automático

**APIs Integradas:**
- **OpenAI GPT**: Para análises complexas
- **Ollama**: Para processamento local
- **LangChain Memory**: Histórico de conversação
- **Tool Registry**: Descoberta dinâmica de ferramentas

### **🔗 Integração VRAgent + Chat Avançado**

#### **Como o VRAgent Atua Dentro do Chat**

Quando você digita uma pergunta no Chat Avançado, não está comunicando diretamente com OpenAI ou Ollama. Na verdade, sua pergunta passa pelo **VRAgent** que atua como um intermediário inteligente, coordenando e otimizando toda a interação.

#### **📊 Evidência no Código**

**Arquivo: `/internal/chat/chat.go:371-392`**
```go
// Ask sends a question to the configured AI service
func (c *Chat) Ask(question string, systemPrompt string, context []Message) (string, error) {
    // Try agent first if configured and enabled
    if c.agent != nil && c.agent.IsEnabled() {
        fmt.Printf("Ask: Trying VR Agent...\n")
        response, err := c.agent.Ask(question)  // ← AQUI: VRAgent processa primeiro
        if err == nil {
            fmt.Printf("Ask: Resposta obtida via agente\n")
            return response, nil
        }
        // Fallback only if agent fails
        fmt.Printf("Warning: Agent request failed, fallback to other services: %v\n", err)
    }

    // Only tries OpenAI/Ollama directly if agent is unavailable/disabled
    // ... fallback logic
}
```

**Configuração do Agente: `/internal/chat/chat.go:360-363`**
```go
// SetAgent configura o agente de IA para o chat
func (c *Chat) SetAgent(agent AgentInterface) {
    c.agent = agent  // ← Chat Service recebe uma instância do VRAgent
}
```

#### **🔄 Fluxo Completo de Consulta no Chat**

```
┌─────────────────────────────────────────────────────┐
│                    USUÁRIO                          │
│           "Quantos colaboradores do SINDPD?"        │
└─────────────────────┬───────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────┐
│                INTERFACE CHAT                       │
│  • Captura pergunta do usuário                     │
│  • Chama chat.Ask(question, systemPrompt, context) │
└─────────────────────┬───────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────┐
│              CHAT SERVICE (Coordenador)             │
│  🎯 Prioridade 1: agent.Ask(question)              │
│  🎯 Prioridade 2: OpenAI (se VRAgent falhar)       │
│  🎯 Prioridade 3: Ollama (se OpenAI falhar)        │
└─────────────────────┬───────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────┐
│                   VR AGENT                          │
│  🧠 Question Classifier: "ProcessedDataQuestion"    │
│  🎯 Roteamento: "usar dados processados"           │
│  ⚡ Performance Optimizer: formatar contexto       │
│  📊 Contextualizar com 1794+ colaboradores         │
└─────────────────────┬───────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────┐
│                  LLM EXECUTION                      │
│  💻 Ollama: processamento com dados reais          │
│  🌐 OpenAI: fallback se Ollama falhar              │
│  📝 System Prompt: instruções específicas de VR    │
└─────────────────────┬───────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────┐
│              RESPONSE PROCESSING                    │
│  🎨 Response Formatter: formatar resposta          │
│  📋 Cache Management: decidir se cachear           │
│  ✅ Quality Check: validar resposta                │
└─────────────────────┬───────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────┐
│                RESPOSTA FINAL                       │
│  "## 📊 Consulta de Dados                         │
│   **Colaboradores SINDPD Processados:**            │
│   • Total: 247 colaboradores                       │
│   • Elegíveis para VR: 234 (94.7%)"              │
└─────────────────────────────────────────────────────┘
```

#### **🧠 Como o VRAgent Atua como Intermediário**

**1. 🎯 Análise Inteligente:**
- Classifica automaticamente o tipo de pergunta
- Decide qual ferramenta ou dados usar
- Otimiza o contexto conforme necessário

**2. 📊 Contextualização Especializada:**
- Adiciona dados específicos de VR (colaboradores, sindicatos, valores)
- Aplica system prompts especializados em VR/VA
- Gerencia cache para otimizar performance

**3. 🔄 Roteamento Inteligente:**
- PolicyConsultant para perguntas de política
- Dados processados para análises estatísticas
- Hybrid approach com fallbacks automáticos

**4. 🎨 Formatação Profissional:**
- Respostas padronizadas com ícones e estrutura
- Citação de fontes quando aplicável
- Indicadores de confiança

#### **✅ Benefícios da Arquitetura VRAgent + Chat**

- **🎯 Respostas Especializadas**: Sistema otimizado para domínio de VR
- **⚡ Performance Inteligente**: Cache e otimizações automáticas
- **🔄 Redundância**: Múltiplos LLMs com fallback automático
- **📊 Contexto Rico**: Acesso a dados reais dos colaboradores
- **🎨 Formatação Consistente**: Interface padronizada e profissional
- **🧠 Aprendizado**: Question Classifier melhora com uso

**Conclusão:** O Chat não é apenas uma interface de LLM - é um sistema inteligente orquestrado pelo VRAgent que fornece respostas especializadas e contextualizadas sobre processamento de VR.

---

## ⚙️ Configurações

### **Tela de Configurações**

**Abas de Configuração:**
```
┌─────────────────────────────────────────────────────┐
│ [🔑 API]                                            │
├─────────────────────────────────────────────────────┤
│                                                     │
│ 🔑 Configuração de API                              │
│                                                     │
│ Provedor LLM: OpenAI                                │
│ API Key: [sk-...] [👁️]                              │
│                                                     │
│ ──────────────────────────                          │
│                                                     │
│ 🔄 Provedor Ollama:                                 │
│ API Key: [sk-ant-...] [👁️]                          │
│                                                     │
│ [Cancelar][💾 Salvar][💾 Salvar e Fechar]           │
└─────────────────────────────────────────────────────┘
```

## 🎯 Melhores Práticas

### **Fluxo de Trabalho Recomendado**

1. **📋 Preparação Semanal**
   - ✅ Verifique atualizações das planilhas
   - ✅ Confirme configurações de API
   - ✅ Execute teste de conectividade
   - ✅ Revise alertas pendentes

2. **🔄 Processamento Mensal**
   - ✅ Backup das planilhas originais
   - ✅ Execute "Validação de Planilhas" primeiro
   - ✅ Processe apenas após resolver críticas
   
3. **📊 Monitoramento Contínuo**
   - ✅ Responda alertas em até 24h
   - ✅ Mantenha histórico organizado
   - ✅ Documente exceções tratadas


---

## 🆘 Suporte e Recursos Adicionais

- 🔧 **[API Reference](api-reference.md)** - Métodos técnicos disponíveis
- 🆘 **[Troubleshooting](troubleshooting.md)** - Soluções para problemas
- 💡 **[Exemplos](examples/)** - Cases de uso práticos

*A interface gráfica foi projetada para ser intuitiva. Explore e experimente! 🚀*