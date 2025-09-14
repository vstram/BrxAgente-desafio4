# 📖 Guia Completo do Usuário - Agente de IA BrxAgente

Este guia detalha todas as funcionalidades da aplicação desktop BrxAgente e como utilizá-las efetivamente no seu trabalho diário.

## 📑 Índice

1. Interface da Aplicação
1. Funcionalidades Principais
1. Chat Avançado com IA
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